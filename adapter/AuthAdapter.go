package adapter

import (
	"attendance/bean"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func SetUserToClaims(user *bean.User) *bean.Claims {
	return &bean.Claims{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
		Class:    user.Class,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
}
