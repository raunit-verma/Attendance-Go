package repository

import (
	"attendance/util"
	"fmt"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// write sql scripts
type User struct {
	UserID   string
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"fullname"`
	Class    int    `json:"class"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type Attendance struct {
	UserID       string
	AttendanceID uuid.UUID
	PunchInDate  time.Time
	PunchOutDate time.Time
}

func (newUser User) IsNewUserDataMissing() bool {
	if newUser.Username == "" {
		zap.L().Info("Username is empty")
		return true
	} else if newUser.Password == "" {
		zap.L().Info("Password is empty")
		return true
	} else if newUser.FullName == "" {
		zap.L().Info("Fullname is empty")
		return true
	} else if newUser.Class <= 0 || newUser.Class > 12 {
		zap.L().Info("Class constraint failed")
		return true
	} else if newUser.Email != "" && util.IsValidEmail(newUser.Email) {
		zap.L().Info("Not a valid email")
		return true
	} else if newUser.Role != "teacher" && newUser.Role != "student" {
		zap.L().Info("Not a valid role")
		return true
	}
	return false
}

func CreateSchema(db *pg.DB) error {
	models := []interface{}{
		(*User)(nil),
		(*Attendance)(nil),
	}

	for _, model := range models {
		err := db.Model(model).CreateTable(&orm.CreateTableOptions{
			Temp:        false,
			IfNotExists: true,
		})
		if err != nil {
			zap.L().Fatal("Error creating schema", zap.Error(err))
			return err
		} else {
			zap.L().Info("Schema created for ", zap.String("type", fmt.Sprintf("%T", model)))
		}
	}
	return nil
}
