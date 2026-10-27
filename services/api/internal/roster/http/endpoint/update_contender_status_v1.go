package endpoint

import (
	"crcs/platform/httpx"
	"net/http"

	"crcs/core/model/enum"

	"github.com/charmingruby/crcs/services/api/internal/roster/usecase"
)

type UpdateContenderStatusV1Request struct {
	Status string `json:"status" validate:"required,oneof=idle active contracted rejected"`
}

func (e *Endpoint) UpdateContenderStatusV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	seasonID, err := httpx.GetPathParam(r, "id")
	if err != nil {
		httpx.WriteBadRequest(w, r, err)
		return
	}

	contenderID, err := httpx.GetPathParam(r, "cid")
	if err != nil {
		httpx.WriteBadRequest(w, r, err)
		return
	}

	request, err := httpx.ParseRequest[UpdateContenderStatusV1Request](w, r)
	if err != nil {
		return
	}

	if err := e.uc.UpdateContenderStatus(ctx, usecase.UpdateContenderStatusInput{
		SeasonID:    seasonID,
		ContenderID: contenderID,
		Status:      enum.ContenderStatus(request.Status),
	}); err != nil {
		httpx.HandleErr(w, r, err)
		return
	}

	httpx.WriteNoContentResponse(w)
}
