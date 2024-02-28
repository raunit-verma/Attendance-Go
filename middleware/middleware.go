package middleware

import (
	"attendance/bean"
	"attendance/util"
	"encoding/json"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func verifyToken(r *http.Request) (int, string, string) {
	cookie, err := r.Cookie("Authorization")

	if err != nil {
		zap.L().Error("No Cookie found", zap.Error(err))
		return http.StatusUnauthorized, "", ""
	}

	tokenStr := cookie.Value
	claims := bean.Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_KEY")), nil
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

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route, _ := mux.CurrentRoute(r).GetPathTemplate()
		status, username, role := verifyToken(r)
		if status != http.StatusAccepted && route != "/login" {
			zap.L().Error("User not verified", zap.String("Code", "1"))
			w.WriteHeader(status)
			json.NewEncoder(w).Encode(bean.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1})
			return
		}
		context.Set(r, "username", username)
		context.Set(r, "role", role)

		next.ServeHTTP(w, r)
	})
}
