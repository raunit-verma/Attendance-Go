package restHandler

import (
	auth "attendance/api/auth"
	"attendance/repository"
	"attendance/services"
	"attendance/util"
	"encoding/json"
	"net/http"
)

func PunchInHandler(w http.ResponseWriter, r *http.Request) {
	status, username := auth.VerifyToken(r)
	if status != http.StatusAccepted {
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}
	services.PunchInService(username, w, r)

	return
}

func PunchOutHandler(w http.ResponseWriter, r *http.Request) {
	status, username := auth.VerifyToken(r)
	if status != http.StatusAccepted {
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}
	services.PunchOutService(username, w, r)

	return
}
