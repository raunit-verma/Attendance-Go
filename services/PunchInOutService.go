package services

import (
	"attendance/bean"
	"attendance/repository"
	"attendance/util"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type PunchInOutService interface {
	PunchIn(username string) (int, bean.ErrorJSON)
	PunchOut(username string) (int, bean.ErrorJSON)
}

type PunchInOutServiceImpl struct {
	repository repository.Repository
}

func NewPunchInOutServiceImpl(repository repository.Repository) *PunchInOutServiceImpl {
	return &PunchInOutServiceImpl{repository: repository}
}

func (impl *PunchInOutServiceImpl) PunchIn(username string) (int, bean.ErrorJSON) {
	user := impl.repository.GetUser(username)
	if user == nil {
		zap.L().Error("User not authorized.", zap.String("Username", username))
		return http.StatusUnauthorized, bean.ErrorJSON{ErrorCode: 1, Message: util.NotAuthorized_One}
	}

	t := time.Now()
	startDate, _ := util.FormateDateTime(t.Year(), t.Month(), t.Day(), 0, 0, 0)
	endDate, _ := util.FormateDateTime(t.Year(), t.Month(), t.Day(), 23, 59, 59)

	currentStatus, _ := impl.repository.GetCurrentStatus(user.Username, startDate, endDate)

	if currentStatus {
		return http.StatusBadRequest, bean.ErrorJSON{ErrorCode: 9, Message: util.OperationNotAllowed_Nine + " Punch out first before punching in again."}
	}

	_, currentTime := util.FormateDateTime(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

	err := impl.repository.AddNewPunchIn(user.Username, currentTime)
	if err != nil {
		zap.L().Error("Error doing operation on DB.", zap.Error(err))
		return http.StatusInternalServerError, bean.ErrorJSON{Message: util.DBError_Seven, ErrorCode: 7}
	}
	return http.StatusAccepted, bean.ErrorJSON{Message: util.Success_Eight + " Punched in successfully.", ErrorCode: 8}
}

func (impl *PunchInOutServiceImpl) PunchOut(username string) (int, bean.ErrorJSON) {
	user := impl.repository.GetUser(username)
	if user == nil {
		zap.L().Error("User not authorized.", zap.String("Username", username))

		return http.StatusUnauthorized, bean.ErrorJSON{ErrorCode: 1, Message: util.NotAuthorized_One}
	}

	t := time.Now()
	startDate, _ := util.FormateDateTime(t.Year(), t.Month(), t.Day(), 0, 0, 0)
	endDate, _ := util.FormateDateTime(t.Year(), t.Month(), t.Day(), 23, 59, 59)

	currentStatus, punchIn := impl.repository.GetCurrentStatus(user.Username, startDate, endDate)

	if !currentStatus {
		return http.StatusBadRequest, bean.ErrorJSON{ErrorCode: 9, Message: util.OperationNotAllowed_Nine + " Punch in first before punching out."}
	}

	_, currentTime := util.FormateDateTime(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	punchIn[0].PunchOutDate = currentTime

	err := impl.repository.AddNewPunchOut(user.Username, punchIn[0], currentTime)

	if err != nil {
		zap.L().Error("Error doing operation on DB.", zap.Error(err))
		return http.StatusInternalServerError, bean.ErrorJSON{Message: util.DBError_Seven, ErrorCode: 7}
	}
	return http.StatusAccepted, bean.ErrorJSON{Message: util.Success_Eight + " Punched out successfully.", ErrorCode: 8}
}
