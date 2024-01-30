package services

import (
	"attendance/repository"
	"attendance/util"
	"net/http"

	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

func AddNewUserService(newUser repository.User, username string, w http.ResponseWriter, r *http.Request) {
	util.TrimSpacesFromStruct(&newUser)
	if newUser.IsNewUserDataMissing() {
		zap.L().Error("New user data is missing")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user := repository.GetUser(username)
	if user.Role != "principal" {
		zap.L().Warn("Unauthorized to add new user")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err := repository.AddNewUser(&newUser)
	pgErr, ok := err.(pg.Error)
	if err != nil {
		if ok && pgErr.Field('C') == "23505" {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	zap.L().Info("New user added succesfully.")
	w.WriteHeader(http.StatusAccepted)
	return
}
