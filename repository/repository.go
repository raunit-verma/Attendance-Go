package repository

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

func GetUser(username string) *User {
	db := GetDB()
	user := &User{}

	err := db.Model(user).Where("username=?", username).Select()

	if err != nil {
		zap.L().Info(fmt.Sprintf("Username %v doesn't exist", username))
		return nil
	}
	return user
}

func AddNewUser(user *User) error {
	db := GetDB()
	_, err := db.Model(user).Insert()
	if err != nil {
		zap.L().Info("Error adding new user to DB.", zap.Error(err))
		return err
	}
	return nil
}

func GetCurrentStatus(username string) bool {
	currentTimeUTC := time.Now().UTC()
	indiaTimeZone, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		zap.L().Info("Error loading Indian timezone")
	}
	currentTimeIndia := currentTimeUTC.In(indiaTimeZone)
	fmt.Print(currentTimeIndia)
	return false
}

func AddNewPunchIn(attendance *Attendance) {

}
