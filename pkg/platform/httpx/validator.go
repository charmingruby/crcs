package httpx

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	val "github.com/go-playground/validator/v10"
)

type validator struct {
	validator *val.Validate
}

type ctxKey struct{}

func newValidator() *validator {
	return &validator{
		validator: val.New(),
	}
}

func (v *validator) Validate(s any) error {
	// Top-level array payloads have no struct to attach tags to;
	// validate each element instead.
	if reflect.ValueOf(s).Kind() == reflect.Slice {
		return v.validateSlice(s)
	}

	err := v.validator.Struct(s)

	if err == nil {
		return nil
	}

	errs := v.unwrapValidationErr(err)

	return errors.New(strings.Join(errs, ", "))
}

func (v *validator) validateSlice(s any) error {
	items := reflect.ValueOf(s)

	for i := range items.Len() {
		item := items.Index(i).Interface()

		if err := v.Validate(item); err != nil {
			return fmt.Errorf("item %d: %w", i, err)
		}
	}

	return nil
}

func validatorToContext(ctx context.Context, val *validator) context.Context {
	return context.WithValue(ctx, ctxKey{}, val)
}

func validatorFromContext(ctx context.Context) *validator {
	val, ok := ctx.Value(ctxKey{}).(*validator)
	if !ok {
		return newValidator()
	}

	return val
}

func (v *validator) unwrapValidationErr(err error) []string {
	var valErr *val.ValidationErrors
	if !errors.As(err, &valErr) {
		return []string{err.Error()}
	}

	errs := *valErr

	reasonsWrapper := make([]string, 0, len(errs))

	for _, vErr := range errs {
		reason := fmt.Sprintf("field `%s` does not satisfy %s rule", vErr.Field(), vErr.Tag())

		reasonsWrapper = append(reasonsWrapper, reason)
	}

	return reasonsWrapper
}
