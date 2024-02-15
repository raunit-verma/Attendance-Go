package restHandler

import (
	auth "attendance/api/auth"
	"attendance/repository"
	"attendance/services"
	"attendance/util"
	"encoding/json"
	"net/http"
)

type FetchStatusHandler interface {
	FetchStatus(w http.ResponseWriter, r *http.Request)
}

type FetchStatusImpl struct {
	fetchStatusService services.FetchStatusService
}

func NewFetchStatusImpl(fetchStatusService services.FetchStatusService) *FetchStatusImpl {
	return &FetchStatusImpl{fetchStatusService: fetchStatusService}
}

func (impl *FetchStatusImpl) FetchStatus(w http.ResponseWriter, r *http.Request) {
	status, username, _ := auth.VerifyToken(r)
	if status != http.StatusAccepted {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}
	impl.fetchStatusService.FetchStatus(w, r, username)
}
