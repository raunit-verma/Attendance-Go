package restHandler

import (
	auth "attendance/api/auth"
	"net/http"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	status, tokenString := auth.CreateToken(r)

	if status != http.StatusAccepted {
		w.WriteHeader(status)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "Authorization",
		Value:   tokenString,
		Expires: time.Now().Add(time.Hour * 24),
	})
	w.WriteHeader(status)
}
