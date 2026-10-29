package endpoint

import (
	"crcs/platform/httpx"
	"net/http"
)

func (e *Endpoint) GetTaskV1(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	taskID, err := httpx.GetPathParam(r, "id")
	if err != nil {
		httpx.WriteBadRequest(w, r, err)
		return
	}

	task, err := e.uc.GetTask(ctx, taskID)
	if err != nil {
		httpx.HandleErr(w, r, err)
		return
	}

	httpx.WriteOKResponse(w, taskFromModel(task))
}
