package restHandler

import (
	auth "attendance/api/auth"
	"attendance/repository"
	"encoding/json"
	"net/http"
	"time"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	status, tokenString, username := auth.CreateToken(r)

	if status != http.StatusAccepted {
		w.WriteHeader(status)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "Authorization",
		Value:   tokenString,
		Expires: time.Now().Add(time.Hour * 24),
	})
	user := repository.GetUser(username)
	user.Password = ""
	json.NewEncoder(w).Encode(user)
	w.WriteHeader(status)
}
