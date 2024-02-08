package restHandler

import (
	auth "attendance/api/auth"
	"attendance/repository"
	"attendance/util"
	"encoding/json"
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	status, username := auth.VerifyToken(r)

	if status != http.StatusAccepted {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}
	w.Write([]byte(fmt.Sprintf("Welcome %v", username)))
}
