package services

import (
	"attendance/repository"
	"encoding/json"
	"net/http"
)

type FetchStatusService interface {
	FetchStatus(w http.ResponseWriter, r *http.Request, username string)
}

type FetchStatusImpl struct {
	repository repository.Repository
}

func NewFetchStatusImpl(repository repository.Repository) *FetchStatusImpl {
	return &FetchStatusImpl{repository: repository}
}

func (impl *FetchStatusImpl) FetchStatus(w http.ResponseWriter, r *http.Request, username string) {
	status, _ := impl.repository.GetCurrentStatus(username)
	json.NewEncoder(w).Encode(map[string]bool{"status": status})
}
