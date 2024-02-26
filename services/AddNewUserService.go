package services

import (
	"attendance/bean"
	"attendance/repository"
	"attendance/util"
	"net/http"

	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

type AddNewUserService interface {
	AddNewUser(newUser bean.User, username string) (int, bean.ErrorJSON)
}

type AddNewUserServiceImpl struct {
	repository repository.Repository
}

func NewAddNewUserServiceImpl(repository repository.Repository) *AddNewUserServiceImpl {
	return &AddNewUserServiceImpl{repository: repository}
}

func (impl *AddNewUserServiceImpl) AddNewUser(newUser bean.User, username string) (int, bean.ErrorJSON) {
	util.TrimSpacesFromStruct(&newUser)
	status, flag, errorJSON := newUser.IsNewUserDataMissing()
	if flag {
		zap.L().Error("New user data is missing")
		return status, errorJSON
	}

	flag, message := util.IsStrongPassword(newUser.Password)

	if !flag {
		return http.StatusBadRequest, bean.ErrorJSON{Message: util.PasswordNotStrong_Ten + message, ErrorCode: 10}
	}

	user := impl.repository.GetUser(username)

	if user != nil && user.Role != "principal" {
		zap.L().Warn("Unauthorized to add new user")
		return http.StatusUnauthorized, bean.ErrorJSON{Message: util.NotAuthorized_One, ErrorCode: 1}
	}

	hashedPassword, err := util.GenerateHashFromPassword(newUser.Password)
	if err != nil {
		zap.L().Info("Error in hashing password.", zap.Error(err))
		return http.StatusInternalServerError, bean.ErrorJSON{Message: util.InternalServererror_Eleven, ErrorCode: 11}
	}
	newUser.Password = hashedPassword
	repoUser := repository.User(newUser)
	err = impl.repository.AddNewUser(&repoUser)
	pgErr, ok := err.(pg.Error)
	if err != nil {
		if ok && pgErr.Field('C') == "23505" {
			zap.L().Error("Username already exist", zap.String("username", newUser.Username))
			return http.StatusBadRequest, bean.ErrorJSON{ErrorCode: 4, Message: util.UsernameTaken_Four}
		}
		return http.StatusInternalServerError, bean.ErrorJSON{Message: util.DBError_Seven, ErrorCode: 7}
	}

	zap.L().Info("New user added succesfully.")
	return http.StatusAccepted, bean.ErrorJSON{Message: util.Success_Eight + " Added new user " + username, ErrorCode: 8}
}
