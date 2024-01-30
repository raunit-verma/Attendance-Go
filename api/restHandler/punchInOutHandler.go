package restHandler

import (
	auth "attendance/api/auth"
	"attendance/services"
	"net/http"
)

func PunchInHandler(w http.ResponseWriter, r *http.Request) {
	status, username := auth.VerifyToken(r)
	if status != http.StatusAccepted {
		w.WriteHeader(status)
		return
	}
	services.PunchInService(username, w, r)

	return
}

func PunchOutHandler(w http.ResponseWriter, r *http.Request) {
	status, username := auth.VerifyToken(r)
	if status != http.StatusAccepted {
		w.WriteHeader(status)
		return
	}
	services.PunchOutService(username, w, r)

	return
}
