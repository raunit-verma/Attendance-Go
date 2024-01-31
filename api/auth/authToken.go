package auth

import (
	"attendance/repository"
	"attendance/util"
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
		return http.StatusBadRequest, "", ""
	}

	user := repository.GetUser(credentials.Username)

	if user == nil || user.Password != credentials.Password {
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
		return http.StatusInternalServerError, "", ""
	}

	return http.StatusAccepted, tokenString, user.Username
}

func VerifyToken(r *http.Request) (int, string) {
	cookie, err := r.Cookie("Authorization")

	if err != nil {
		if err == http.ErrNoCookie {
			return http.StatusUnauthorized, ""
		}
		return http.StatusBadRequest, ""
	}

	tokenStr := cookie.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return http.StatusUnauthorized, ""
		}
		return http.StatusBadRequest, ""
	}

	if !token.Valid {
		return http.StatusUnauthorized, ""
	}
	return http.StatusAccepted, claims.Username
}
