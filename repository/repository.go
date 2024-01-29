package repository

import (
	"fmt"

	"go.uber.org/zap"
)

func CheckIfUserExists(username string) *User {
	db := GetDB()
	user := User{
		Username: username,
	}

	err := db.Model(user).WherePK().Select()
	if err != nil {
		return nil
	}
	return &user
}

func AddNewUser(user *User) error {
	db := GetDB()
	_, err := db.Model(user).Insert()
	if err != nil {
		zap.L().Error("Error adding new user to db", zap.Error(err))
		fmt.Println(err)
		return err
	}
	return nil
}
