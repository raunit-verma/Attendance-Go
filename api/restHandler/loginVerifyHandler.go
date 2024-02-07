package restHandler

import (
	auth "attendance/api/auth"
	"attendance/repository"
	"attendance/util"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	status, tokenString, username := auth.CreateToken(r)

	if status != http.StatusAccepted {
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: false,
		Secure:   false,
		// Domain:   os.Getenv("DOMAIN"), // production
		Domain: "",
		Path:   "/",
	})
	user := repository.GetUser(username)
	if user.Username == "" {
		zap.L().Info("No user found", zap.String("Username", username))
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.UserNotFound_Six, ErrorCode: 6})
		return
	}
	user.Password = ""
	json.NewEncoder(w).Encode(user)
	zap.L().Info("User logged in", zap.String("username", user.Username))
}

func VerifyToken(w http.ResponseWriter, r *http.Request) {
	status, username := auth.VerifyToken(r)
	if status != http.StatusAccepted {
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}
	user := repository.GetUser(username)
	user.Password = ""
	json.NewEncoder(w).Encode(user)
}
