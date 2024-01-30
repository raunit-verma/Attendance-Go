package restHandler

import (
	auth "attendance/api/auth"
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	status, username := auth.VerifyToken(r)

	if status != http.StatusAccepted {
		w.WriteHeader(status)
		return
	}
	w.Write([]byte(fmt.Sprintf("Welcome %v", username)))
}
