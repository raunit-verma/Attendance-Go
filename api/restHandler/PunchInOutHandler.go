package restHandler

import (
	auth "attendance/api/auth"
	"attendance/services"
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
)

type PunchInOutHandler interface {
	PunchIn(w http.ResponseWriter, r *http.Request)
	PunchOut(w http.ResponseWriter, r *http.Request)
}

type PunchInOutImpl struct {
	punchInOutService services.PunchInOutService
	authService       auth.AuthService
}

func NewPunchInOutImpl(punchInOutService services.PunchInOutService, auth auth.AuthService) *PunchInOutImpl {
	return &PunchInOutImpl{punchInOutService: punchInOutService, authService: auth}
}

func (impl *PunchInOutImpl) PunchIn(w http.ResponseWriter, r *http.Request) {
	username := context.Get(r, "username").(string)

	status, errorJSON := impl.punchInOutService.PunchIn(username)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorJSON)
}

func (impl *PunchInOutImpl) PunchOut(w http.ResponseWriter, r *http.Request) {
	username := context.Get(r, "username").(string)

	status, errorJSON := impl.punchInOutService.PunchOut(username)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorJSON)
}
