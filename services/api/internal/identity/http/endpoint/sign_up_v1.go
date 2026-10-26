package endpoint

import (
	"crcs/platform/httpx"
	"net/http"

	"github.com/charmingruby/crcs/services/api/internal/identity/usecase"
)

type SignUpV1Request struct {
	FirstName string `json:"first_name" validate:"required,min=1"`
	LastName  string `json:"last_name"  validate:"required,min=1"`
	Email     string `json:"email"      validate:"required,email"`
	Password  string `json:"password"   validate:"required,min=8"`
}

type SignUpV1Response struct {
	AccessToken string `json:"access_token"`
}

func (e *Endpoint) SignUpV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	request, err := httpx.ParseRequest[SignUpV1Request](w, r)
	if err != nil {
		return
	}

	output, err := e.uc.SignUp(ctx, usecase.SignUpInput{
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Email:     request.Email,
		Password:  request.Password,
	})
	if err != nil {
		httpx.HandleErr(w, r, err)
		return
	}

	httpx.WriteCreatedResponse(w, SignUpV1Response{
		AccessToken: output.AccessToken,
	})
}
