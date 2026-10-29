package httpx

import (
	"crcs/core/customerr"
	"crcs/platform/o11y/logging"
	"encoding/json"
	"errors"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type Page[T any] struct {
	Items  []T `json:"items"`
	Total  int `json:"total"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

func WriteOKResponse(w http.ResponseWriter, v any) {
	writeResponse(w, http.StatusOK, v)
}

func WriteCreatedResponse(w http.ResponseWriter, v any) {
	writeResponse(w, http.StatusCreated, v)
}

func WriteNoContentResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func WriteAcceptedResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusAccepted)
}

func WriteBadRequest(w http.ResponseWriter, r *http.Request, err error) {
	HandleErr(w, r, customerr.Validation(err.Error()))
}

func HandleErr(w http.ResponseWriter, r *http.Request, err error) {
	log := logging.LoggerFromContext(r.Context())

	if cerr, ok := errors.AsType[*customerr.Error](err); ok {
		httpErr := customerr.MapToHTTP(*cerr)

		err := cerr.Unwrap()
		originalErrMsg := ""
		if err != nil {
			originalErrMsg = err.Error()
		}

		log.Error("application error",
			"type", httpErr.Type,
			"message", httpErr.Message,
			"original_error", originalErrMsg,
		)

		writeResponse(w, httpErr.Status, ErrorResponse{
			Message: httpErr.Message,
		})

		return
	}

	log.Error("internal server error", "error", err)

	writeResponse(w, customerr.UnhandledHTTPError.Status, ErrorResponse{
		Message: customerr.UnhandledHTTPError.Message,
	})
}

func writeResponse(w http.ResponseWriter, status int, v any) {
	if v == nil {
		w.WriteHeader(status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	body, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"message":"Internal Server Error"}`))
		return
	}

	w.WriteHeader(status)
	_, _ = w.Write(body)
}
