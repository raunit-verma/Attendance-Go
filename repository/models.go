package repository

import (
	"attendance/util"
	"time"

	"github.com/go-pg/pg"
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
	// models := []interface{}{
	// 	(*User)(nil),
	// 	(*Attendance)(nil),
	// }
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS "users" (
		"user_id" VARCHAR(255),
		"password" VARCHAR(255),
		"full_name" VARCHAR(255),
		"class" INTEGER,
		"email" VARCHAR(255),
		"role" VARCHAR(255),
		PRIMARY KEY ("user_id")
	  );`)

	if err != nil {
		zap.L().Fatal("Error creating schema for users", zap.Error(err))
		return err
	} else {
		zap.L().Info("Schema created for users")
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "attendances" (
		"user_id" VARCHAR(255),
		"attendance_id" VARCHAR(255),
		"punch_in_date" DATE,
		"punch_out_date" DATE,
		PRIMARY KEY ("attendance_id")
	  );
	  `)

	if err != nil {
		zap.L().Fatal("Error creating schema for attendance", zap.Error(err))
		return err
	} else {
		zap.L().Info("Schema created for attendance")
	}

	_, err = db.Exec(`ALTER TABLE IF EXISTS "attendances" ADD CONSTRAINT "fk_user_attendance" FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");`)
	pgErr, ok := err.(pg.Error)

	if ok && pgErr.Field('C') == "42710" {
		zap.L().Warn("Foreign key already exist")
	} else if err != nil {
		zap.L().Fatal("Error creating foregin key for attendance", zap.Error(err))
		return err
	} else {
		zap.L().Info("Foreign key created")
	}

	// for _, model := range models {
	// 	err := db.Model(model).CreateTable(&orm.CreateTableOptions{
	// 		Temp:        false,
	// 		IfNotExists: true,
	// 	})
	// 	if err != nil {
	// 		zap.L().Fatal("Error creating schema", zap.Error(err))
	// 		return err
	// 	} else {
	// 		zap.L().Info("Schema created for ", zap.String("type", fmt.Sprintf("%T", model)))
	// 	}
	// }
	return nil
}
