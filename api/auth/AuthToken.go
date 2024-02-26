package auth

import (
	"attendance/bean"
	"attendance/repository"
	"attendance/util"
	"encoding/json"
	"net/http"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type AuthToken interface {
	CreateToken(r *http.Request) (int, string, string)
	VerifyToken(r *http.Request) (int, string, string)
}

type AuthTokenImpl struct {
	repository repository.Repository
	cfg        bean.AuthConfig
}

func NewAuthTokenImpl(repository repository.Repository) *AuthTokenImpl {
	cfg := bean.AuthConfig{}
	if err := env.Parse(&cfg); err != nil {
		zap.L().Error("Error loading env.", zap.Error(err))
	}
	return &AuthTokenImpl{repository: repository, cfg: cfg}
}

type Claims struct {
	Username string `json:"username"`
	FullName string `json:"fullname"`
	Class    int    `json:"class"`
	Email    string `json:"email,omitempty"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type Credentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func (impl *AuthTokenImpl) CreateToken(r *http.Request) (int, string, string) {

	var credentials Credentials

	err := json.NewDecoder(r.Body).Decode(&credentials)
	util.TrimSpacesFromStruct(&credentials)

	if err != nil {
		zap.L().Error("Cannot decode username and password", zap.Error(err))
		return http.StatusBadRequest, "", ""
	}

	user := impl.repository.GetUser(credentials.Username)

	if user.Username != credentials.Username || user.Username == "" || user.Password == "" || !util.MatchPassword([]byte(user.Password), []byte(credentials.Password)) {
		zap.L().Info("Wrong username or password", zap.String("Passed Credentials", user.Username))
		return http.StatusUnauthorized, "", ""
	}

	expirationtime := time.Now().Add(time.Hour * 24)

	claims := &Claims{
		Username: credentials.Username,
		FullName: user.FullName,
		Email:    user.Email,
		Class:    user.Class,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationtime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(impl.cfg.JwtKey))

	if err != nil {
		zap.L().Error("Couldn't create token string", zap.Error(err))
		return http.StatusInternalServerError, "", ""
	}

	return http.StatusAccepted, tokenString, user.Username
}

func (impl *AuthTokenImpl) VerifyToken(r *http.Request) (int, string, string) {
	cookie, err := r.Cookie("Authorization")

	if err != nil {
		if err == http.ErrNoCookie {
			zap.L().Error("No Cookie found", zap.Error(err))
			return http.StatusUnauthorized, "", ""
		}
		zap.L().Error("Cannot retrieve cookie", zap.Error(err))
		return http.StatusBadRequest, "", ""
	}

	tokenStr := cookie.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(impl.cfg.JwtKey), nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			zap.L().Error("Invalid token", zap.Error(err))
			return http.StatusUnauthorized, "", ""
		}
		zap.L().Error("Error verifying token", zap.Error(err))
		return http.StatusBadRequest, "", ""
	}

	if !token.Valid {
		zap.L().Error("Token not valid")
		return http.StatusUnauthorized, "", ""
	}
	return http.StatusAccepted, claims.Username, claims.Role
}
