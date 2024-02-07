package services

import (
	"attendance/repository"
	"encoding/json"
	"net/http"
)

func FetchStatusService(w http.ResponseWriter, r *http.Request, username string) {
	status, _ := repository.GetCurrentStatus(username)
	json.NewEncoder(w).Encode(map[string]bool{"status": status})
}
