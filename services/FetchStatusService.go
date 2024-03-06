package services

import (
	"attendance/repository"
	"attendance/util"
	"time"
)

type FetchStatusService interface {
	FetchStatus(username string) map[string]bool
}

type FetchStatusImpl struct {
	repository repository.Repository
}

func NewFetchStatusImpl(repository repository.Repository) *FetchStatusImpl {
	return &FetchStatusImpl{repository: repository}
}

func (impl *FetchStatusImpl) FetchStatus(username string) map[string]bool {

	t := time.Now()
	startDate, _ := util.FormateDateTime(t.Year(), t.Month(), t.Day(), 0, 0, 0)
	endDate, _ := util.FormateDateTime(t.Year(), t.Month(), t.Day(), 23, 59, 59)
	status, _ := impl.repository.GetCurrentStatus(username, startDate, endDate)

	return map[string]bool{"status": status}
}
