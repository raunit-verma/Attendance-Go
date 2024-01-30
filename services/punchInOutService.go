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

	currentStatus, _ := repository.GetCurrentStatus(user.Username)

	if currentStatus {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	err := repository.AddNewPunchIn(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	return
}

func PunchOutService(username string, w http.ResponseWriter, r *http.Request) {
	user := repository.GetUser(username)
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	currentStatus, punchIn := repository.GetCurrentStatus(user.Username)

	if !currentStatus {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	err := repository.AddNewPunchOut(user.Username, punchIn[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	return
}
