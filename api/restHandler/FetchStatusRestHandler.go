package restHandler

import (
	auth "attendance/api/auth"
	"attendance/services"
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
)

type FetchStatusHandler interface {
	FetchStatus(w http.ResponseWriter, r *http.Request)
}

type FetchStatusImpl struct {
	fetchStatusService services.FetchStatusService
	authService        auth.AuthService
}

func NewFetchStatusImpl(fetchStatusService services.FetchStatusService, authService auth.AuthService) *FetchStatusImpl {
	return &FetchStatusImpl{fetchStatusService: fetchStatusService, authService: authService}
}

func (impl *FetchStatusImpl) FetchStatus(w http.ResponseWriter, r *http.Request) {
	username := context.Get(r, "username").(string)
	student_status := impl.fetchStatusService.FetchStatus(username)
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(student_status)
}
