package restHandler

import (
	auth "attendance/api/auth"
	"attendance/bean"
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
	auth               auth.AuthToken
}

func NewFetchStatusImpl(fetchStatusService services.FetchStatusService, auth auth.AuthToken) *FetchStatusImpl {
	return &FetchStatusImpl{fetchStatusService: fetchStatusService, auth: auth}
}

func (impl *FetchStatusImpl) FetchStatus(w http.ResponseWriter, r *http.Request) {
	status, username, _ := impl.auth.VerifyToken(r)
	if status != http.StatusAccepted {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(bean.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}
	student_status := impl.fetchStatusService.FetchStatus(username)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(student_status)
}
