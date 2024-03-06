package restHandler

import (
	"attendance/adapter"
	auth "attendance/api/auth"
	"attendance/bean"
	"attendance/repository"
	"attendance/util"
	"encoding/json"
	"net/http"

	"github.com/caarlos0/env/v10"
	"go.uber.org/zap"
)

type LoginHandler interface {
	Login(w http.ResponseWriter, r *http.Request)
}

type LoginImpl struct {
	repository  repository.Repository
	authService auth.AuthService
	cfg         bean.CookieConfig
}

func NewLoginImpl(repository repository.Repository, auth auth.AuthService) *LoginImpl {
	cfg := bean.CookieConfig{}
	if err := env.Parse(&cfg); err != nil {
		zap.L().Error("Error loading env.", zap.Error(err))
	}
	return &LoginImpl{repository: repository, authService: auth, cfg: cfg}
}

func (impl *LoginImpl) Login(w http.ResponseWriter, r *http.Request) {
	status, tokenString, username := impl.authService.CreateToken(r)

	if status != http.StatusAccepted {
		w.WriteHeader(status)
		json.NewEncoder(w).Encode(bean.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
		return
	}

	if impl.cfg.Type == util.PRODUCTION {
		http.SetCookie(w, adapter.SetCookie(util.PRODUCTION, tokenString, impl.cfg.Domain))
	} else {
		http.SetCookie(w, adapter.SetCookie(util.DEVELOPMENT, tokenString, impl.cfg.Domain))
	}

	user, err := impl.repository.GetUser(username)

	if err != nil {
		zap.L().Error(err.Error())
	}

	if user == nil || user.Username == "" {
		zap.L().Info("No user found", zap.String("Username", username))
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(bean.ErrorJSON{Message: util.UserNotFound_Six, ErrorCode: 6})
		return
	}
	user.Password = ""
	json.NewEncoder(w).Encode(user)
	zap.L().Info("User logged in", zap.String("username", user.Username))
}
