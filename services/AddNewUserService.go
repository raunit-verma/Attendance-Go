package services

import (
	"attendance/repository"
	"attendance/util"
	"net/http"

	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

type AddNewUserService interface {
	AddNewUser(newUser repository.User, username string) (int, repository.ErrorJSON)
}

type AddNewUserServiceImpl struct {
	repository repository.Repository
}

func NewAddNewUserServiceImpl(repository repository.Repository) *AddNewUserServiceImpl {
	return &AddNewUserServiceImpl{repository: repository}
}

func (impl *AddNewUserServiceImpl) AddNewUser(newUser repository.User, username string) (int, repository.ErrorJSON) {
	util.TrimSpacesFromStruct(&newUser)
	status, flag, errorJSON := newUser.IsNewUserDataMissing()
	if flag {
		zap.L().Error("New user data is missing")
		return status, errorJSON
	}

	flag, message := util.IsStrongPassword(newUser.Password)

	if !flag {
		return http.StatusBadRequest, repository.ErrorJSON{Message: util.PasswordNotStrong_Ten + message, ErrorCode: 10}
	}

	user := impl.repository.GetUser(username)

	if user != nil && user.Role != "principal" {
		zap.L().Warn("Unauthorized to add new user")
		return http.StatusUnauthorized, repository.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1}
	}

	hashedPassword, err := util.GenerateHashFromPassword(newUser.Password)
	if err != nil {
		zap.L().Info("Error in hashing password.", zap.Error(err))
		return http.StatusInternalServerError, repository.ErrorJSON{Message: util.InternalServererror_Eleven, ErrorCode: 11}
	}
	newUser.Password = hashedPassword

	err = impl.repository.AddNewUser(&newUser)
	pgErr, ok := err.(pg.Error)
	if err != nil {
		if ok && pgErr.Field('C') == "23505" {
			zap.L().Error("Username already exist", zap.String("username", newUser.Username))
			return http.StatusBadRequest, repository.ErrorJSON{ErrorCode: 4, Message: util.UsernameTaken_Four}
		}
		return http.StatusInternalServerError, repository.ErrorJSON{Message: util.DBError_Seven, ErrorCode: 7}
	}

	zap.L().Info("New user added succesfully.")
	return http.StatusAccepted, repository.ErrorJSON{Message: util.Success_Eight + " Added new user " + username, ErrorCode: 8}
}
