package services

import (
	"attendance/repository"
	"attendance/util"
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func PunchInService(username string, w http.ResponseWriter, r *http.Request) {
	user := repository.GetUser(username)
	if user == nil {
		zap.L().Error("User not authorized.", zap.String("Username", username))
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(repository.ErrorJSON{ErrorCode: 1, Message: util.NotAuthorized_One})
		return
	}

	currentStatus, _ := repository.GetCurrentStatus(user.Username)

	if currentStatus {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(repository.ErrorJSON{ErrorCode: 9, Message: util.OperationNotAllowed_Nine + " Punch out first before punching in again."})
		return
	}
	err := repository.AddNewPunchIn(user.Username)
	if err != nil {
		zap.L().Error("Error doing operation on DB.", zap.Error(err))
		w.WriteHeader(util.InternalServererror)
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.DBError_Seven, ErrorCode: 7})
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.Success_Eight + " Punched in successfully.", ErrorCode: 8})
	return
}

func PunchOutService(username string, w http.ResponseWriter, r *http.Request) {
	user := repository.GetUser(username)
	if user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		zap.L().Error("User not authorized.", zap.String("Username", username))
		json.NewEncoder(w).Encode(repository.ErrorJSON{ErrorCode: 1, Message: util.NotAuthorized_One})
		return
	}

	currentStatus, punchIn := repository.GetCurrentStatus(user.Username)

	if !currentStatus {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(repository.ErrorJSON{ErrorCode: 9, Message: util.OperationNotAllowed_Nine + " Punch in first before punching out."})
		return
	}
	err := repository.AddNewPunchOut(user.Username, punchIn[0])
	if err != nil {
		w.WriteHeader(util.InternalServererror)
		zap.L().Error("Error doing operation on DB.", zap.Error(err))
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.DBError_Seven, ErrorCode: 7})
		return
	}
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.Success_Eight + " Punched out successfully.", ErrorCode: 8})
	return
}
