package httpx

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

const (
	maxRequestBodyBytes = 1 << 20

	DefaultLimit = 20
	MaxLimit     = 100
)

var (
	ErrInvalidPayload = errors.New("invalid payload")
	ErrMissingParam   = errors.New("missing param")
	ErrInvalidQuery   = errors.New("invalid query param")
)

func ParseRequest[T any](w http.ResponseWriter, r *http.Request) (*T, error) {
	var obj T

	r.Body = http.MaxBytesReader(w, r.Body, maxRequestBodyBytes)

	if err := json.UnmarshalRead(r.Body, &obj); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]string{
			"message": fmt.Sprintf("%s: %s", ErrInvalidPayload.Error(), err.Error()),
		})

		return nil, err
	}

	val := validatorFromContext(r.Context())
	if err := val.Validate(obj); err != nil {
		writeResponse(w, http.StatusBadRequest, map[string]string{
			"message": fmt.Sprintf("%s: %s", ErrInvalidPayload.Error(), err.Error()),
		})

		return nil, err
	}

	return &obj, nil
}

func GetPathParam(r *http.Request, key string) (string, error) {
	param := chi.URLParam(r, key)
	if param == "" {
		return "", fmt.Errorf("%w: %s", ErrMissingParam, key)
	}

	return param, nil
}

func GetQueryParam(r *http.Request, key string) (string, error) {
	param := r.URL.Query().Get(key)
	if param == "" {
		return "", fmt.Errorf("%w: %s", ErrMissingParam, key)
	}

	return param, nil
}

func QueryInt(r *http.Request, key string, def int) (int, error) {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return def, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrInvalidQuery, key)
	}

	return value, nil
}

func Pagination(r *http.Request) (limit, offset int, err error) {
	limit, err = QueryInt(r, "limit", DefaultLimit)
	if err != nil {
		return 0, 0, err
	}

	if limit < 1 {
		limit = DefaultLimit
	}

	if limit > MaxLimit {
		limit = MaxLimit
	}

	offset, err = QueryInt(r, "offset", 0)
	if err != nil {
		return 0, 0, err
	}

	if offset < 0 {
		offset = 0
	}

	return limit, offset, nil
}
