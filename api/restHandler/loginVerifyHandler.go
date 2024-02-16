package restHandler

import (
	auth "attendance/api/auth"
	"attendance/repository"
	"attendance/util"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

type LoginHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
	VerifyToken(w http.ResponseWriter, r *http.Request)
}

type LoginImpl struct {
	repository repository.Repository
	authToken  auth.AuthToken
}

func NewLoginImpl(repository repository.Repository) *LoginImpl {
	return &LoginImpl{repository: repository}
}

func (impl *LoginImpl) Login(w http.ResponseWriter, r *http.Request) {
	status, tokenString, username := impl.authToken.CreateToken(r)

	if status != http.StatusAccepted {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}

	if os.Getenv("TYPE") == "Production" {
		http.SetCookie(w, &http.Cookie{
			Name:     "Authorization",
			Value:    tokenString,
			Expires:  time.Now().Add(time.Hour * 24),
			HttpOnly: false,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Domain:   os.Getenv("DOMAIN"),
			Path:     "/",
		})
	} else {
		http.SetCookie(w, &http.Cookie{
			Name:    "Authorization",
			Value:   tokenString,
			Expires: time.Now().Add(time.Hour * 24),
			Path:    "/",
		})
	}

	user := impl.repository.GetUser(username)
	if user.Username == "" {
		zap.L().Info("No user found", zap.String("Username", username))
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.UserNotFound_Six, ErrorCode: 6})
		return
	}
	user.Password = ""
	json.NewEncoder(w).Encode(user)
	zap.L().Info("User logged in", zap.String("username", user.Username))
}

func (impl *LoginImpl) VerifyToken(w http.ResponseWriter, r *http.Request) {
	status, username, _ := auth.VerifyToken(r)
	if status != http.StatusAccepted {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(repository.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}
	user := impl.repository.GetUser(username)
	user.Password = ""
	json.NewEncoder(w).Encode(user)
}
