package restHandler

import (
	auth "attendance/api/auth"
	"attendance/bean"
	"attendance/services"
	"attendance/util"
	"encoding/json"
	"net/http"
)

type PunchInOutHandler interface {
	PunchIn(w http.ResponseWriter, r *http.Request)
	PunchOut(w http.ResponseWriter, r *http.Request)
}

type PunchInOutImpl struct {
	punchInOutService services.PunchInOutService
}

func NewPunchInOutImpl(punchInOutService services.PunchInOutService) *PunchInOutImpl {
	return &PunchInOutImpl{punchInOutService: punchInOutService}
}

func (impl *PunchInOutImpl) PunchIn(w http.ResponseWriter, r *http.Request) {
	status, username, _ := auth.VerifyToken(r)
	if status != http.StatusAccepted {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(bean.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}
	status, errorJSON := impl.punchInOutService.PunchIn(username)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorJSON)
}

func (impl *PunchInOutImpl) PunchOut(w http.ResponseWriter, r *http.Request) {
	status, username, _ := auth.VerifyToken(r)
	if status != http.StatusAccepted {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(bean.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}
	status, errorJSON := impl.punchInOutService.PunchOut(username)
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorJSON)
}
