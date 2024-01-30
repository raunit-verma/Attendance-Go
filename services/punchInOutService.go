package services

import (
	"attendance/repository"
	"net/http"
)

func PunchInService(username string, w http.ResponseWriter, r *http.Request) {
	user := repository.GetUser(username)
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	currentStatus := repository.GetCurrentStatus(user.Username)
	if currentStatus {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	return
}
