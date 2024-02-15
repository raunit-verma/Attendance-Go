package restHandler

import (
	auth "attendance/api/auth"
	"attendance/repository"
	"attendance/services"
	"attendance/util"
	"encoding/json"
	"net/http"
)

type PunchInOutHandler interface {
	PunchInHandler(w http.ResponseWriter, r *http.Request)
	PunchOutHandler(w http.ResponseWriter, r *http.Request)
}

type PunchInOutImpl struct {
	punchInOutService services.PunchInOutService
}

func NewPunchInOutImpl(punchInOutService services.PunchInOutService) *PunchInOutImpl {
	return &PunchInOutImpl{punchInOutService: punchInOutService}
}

func (impl *PunchInOutImpl) PunchInHandler(w http.ResponseWriter, r *http.Request) {
	status, username, _ := auth.VerifyToken(r)
	if status != http.StatusAccepted {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}
	impl.punchInOutService.PunchInService(username, w, r)

	return
}

func (impl *PunchInOutImpl) PunchOutHandler(w http.ResponseWriter, r *http.Request) {
	status, username, _ := auth.VerifyToken(r)
	if status != http.StatusAccepted {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}
	impl.punchInOutService.PunchOutService(username, w, r)

	return
}
