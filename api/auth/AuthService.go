package auth

import (
	"attendance/adapter"
	"attendance/bean"
	"attendance/repository"
	"attendance/util"
	"encoding/json"
	"net/http"

	"github.com/caarlos0/env/v10"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type AuthService interface {
	CreateToken(r *http.Request) (int, string, string)
	VerifyToken(r *http.Request) (int, string, string)
}

type AuthServiceImpl struct {
	repository repository.Repository
	cfg        bean.AuthConfig
}

func NewAuthServiceImpl(repository repository.Repository) *AuthServiceImpl {
	cfg := bean.AuthConfig{}
	if err := env.Parse(&cfg); err != nil {
		zap.L().Error("Error loading env.", zap.Error(err))
	}
	return &AuthServiceImpl{repository: repository, cfg: cfg}
}

func (impl *AuthServiceImpl) CreateToken(r *http.Request) (int, string, string) {

	var credentials bean.Credentials

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		zap.L().Error("Cannot decode username and password", zap.Error(err))
		return http.StatusBadRequest, "", ""
	}

	util.TrimSpacesFromStruct(&credentials)

	user, _ := impl.repository.GetUser(credentials.Username)

	if user == nil || user.Username == "" || user.Password == "" || !util.MatchPassword([]byte(user.Password), []byte(credentials.Password)) {
		zap.L().Info("Wrong username or password", zap.String("Passed Credentials", user.Username))
		return http.StatusUnauthorized, "", ""
	}

	convertedUser := (bean.User)(*user)

	claims := adapter.SetUserToClaims(&convertedUser)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(impl.cfg.JwtKey))

	if err != nil {
		zap.L().Error("Couldn't create token string", zap.Error(err))
		return http.StatusInternalServerError, "", ""
	}

	return http.StatusAccepted, tokenString, user.Username
}

func (impl *AuthServiceImpl) VerifyToken(r *http.Request) (int, string, string) {
	cookie, err := r.Cookie("Authorization")

	if err != nil {
		zap.L().Error("No Cookie found", zap.Error(err))
		return http.StatusUnauthorized, "", ""
	}

	tokenStr := cookie.Value
	claims := bean.Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(impl.cfg.JwtKey), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			zap.L().Error("Invalid token", zap.Error(err))
			return http.StatusUnauthorized, "", ""
		}
		zap.L().Error("Error verifying token", zap.Error(err))
		return http.StatusForbidden, "", ""
	}

	if !token.Valid {
		zap.L().Error("Token not valid")
		return http.StatusUnauthorized, "", ""
	}
	return http.StatusAccepted, claims.Username, claims.Role
}
