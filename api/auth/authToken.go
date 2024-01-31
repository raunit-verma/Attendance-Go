package auth

import (
	"attendance/repository"
	"attendance/util"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

var jwtKey = []byte("Raunit-Verma")

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateToken(r *http.Request) (int, string, string) {

	var credentials Credentials

	err := json.NewDecoder(r.Body).Decode(&credentials)
	util.TrimSpacesFromStruct(&credentials)

	if err != nil {
		zap.L().Error("Cannot decode username and password", zap.Error(err))
		return http.StatusBadRequest, "", ""
	}

	user := repository.GetUser(credentials.Username)

	if user == nil || user.Password != credentials.Password {
		zap.L().Info("Wrong username or password", zap.String("Data", user.Username))
		return http.StatusUnauthorized, "", ""
	}

	expirationtime := time.Now().Add(time.Hour * 24)

	claims := &Claims{
		Username: credentials.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationtime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		zap.L().Error("Couldn't create token string", zap.Error(err))
		return http.StatusInternalServerError, "", ""
	}

	return http.StatusAccepted, tokenString, user.Username
}

func VerifyToken(r *http.Request) (int, string) {
	cookie, err := r.Cookie("Authorization")

	if err != nil {
		if err == http.ErrNoCookie {
			zap.L().Error("No Cookie found", zap.Error(err))
			return http.StatusUnauthorized, ""
		}
		zap.L().Error("Cannot retrieve cookie", zap.Error(err))
		return http.StatusBadRequest, ""
	}

	tokenStr := cookie.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			zap.L().Error("Invalid token", zap.Error(err))
			return http.StatusUnauthorized, ""
		}
		zap.L().Error("Error verifying token", zap.Error(err))
		return http.StatusBadRequest, ""
	}

	if !token.Valid {
		zap.L().Error("Token not valid")
		return http.StatusUnauthorized, ""
	}
	zap.L().Info("Token verified", zap.String("user", claims.Username))
	return http.StatusAccepted, claims.Username
}
