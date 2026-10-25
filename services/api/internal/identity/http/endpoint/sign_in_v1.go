package endpoint

import (
	"net/http"

	"github.com/charmingruby/crcs/services/api/internal/identity/usecase"
	"github.com/charmingruby/crcs/services/api/internal/platform/httpx"
)

type SignInV1Request struct {
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=1"`
}

type SignInV1Response struct {
	AccessToken string `json:"access_token"`
}

func (e *Endpoint) SignInV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request, err := httpx.ParseRequest[SignInV1Request](w, r)
	if err != nil {
		return
	}

	output, err := e.uc.SignIn(ctx, usecase.SignInInput{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		httpx.HandleErr(w, r, err)
		return
	}

	httpx.WriteOKResponse(w, SignInV1Response{
		AccessToken: output.AccessToken,
	})
}
