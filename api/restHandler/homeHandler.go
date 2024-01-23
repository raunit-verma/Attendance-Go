package restHandler

import (
	auth "attendance/services"
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	status, username := auth.VerifyToken(r)

	if status != http.StatusAccepted {
		w.WriteHeader(status)
		return
	}
	fmt.Println(username)
	w.Write([]byte(fmt.Sprintf("Welcome %v", username)))
}
