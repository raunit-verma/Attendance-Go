package repository

import (
	"attendance/util"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-pg/pg"
	"go.uber.org/zap"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	FullName string `json:"fullname"`
	Class    int    `json:"class"`
	Email    string `json:"email,omitempty"`
	Role     string `json:"role"`
}

type Attendance struct {
	Username     string    `pg:"username" json:"-"`
	AttendanceID string    `pg:"attendance_id,pk" json:"-"`
	PunchInDate  time.Time `pg:"punch_in_date"`
	PunchOutDate time.Time `pg:"punch_out_date"`
}

type AttendanceJSON struct {
	PunchInDate  time.Time `pg:"punch_in_date"`
	PunchOutDate time.Time `pg:"punch_out_date"`
}

type GetTeacherAttendanceJSON struct {
	ID    string `json:"id"`
	Month int    `json:"month"`
	Year  int    `json:"year"`
}

type StudentAttendanceJSON struct {
	TableName struct{} `sql:"users" json:"-"`
	Username  string   `pg:"username"`
	FullName  string   `pg:"full_name"`
}

type ErrorJSON struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"code"`
}

type GetClassAttendanceJSON struct {
	Class int `json:"class"`
	Day   int `json:"day"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

type GetStudentAttendanceJSON struct {
	Month int `json:"month"`
	Year  int `json:"year"`
}

func (newUser User) IsNewUserDataMissing(w http.ResponseWriter, r *http.Request) bool {
	IsDataMissing := false
	Message := ""

	if newUser.Username == "" {
		IsDataMissing = true
		Message = " Username is missing."
		zap.L().Info("Username is empty")
	} else if newUser.Password == "" {
		IsDataMissing = true
		zap.L().Info("Password is empty")
		Message = " Password is missing."
	} else if newUser.FullName == "" {
		IsDataMissing = true
		zap.L().Info("Fullname is empty")
		Message = " Fullname is missing."
	} else if newUser.Class <= 0 || newUser.Class > 12 {
		IsDataMissing = true
		zap.L().Info("Class constraint failed")
		Message = " Class should be between 1 to 12."
	} else if newUser.Email != "" && !util.IsValidEmail(newUser.Email) {
		IsDataMissing = true
		zap.L().Info("Not a valid email")
		Message = " Email is missing or not a valid email."
	} else if newUser.Role != "teacher" && newUser.Role != "student" {
		zap.L().Info("Not a valid role")
		Message = " Role is missing."
	}

	if IsDataMissing {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorJSON{ErrorCode: 3, Message: util.UserDataMissing_Three + Message})
	}

	return IsDataMissing
}

func CreateSchema(db *pg.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS "users" (
		"username" VARCHAR(255) PRIMARY KEY,
		"password" VARCHAR(255),
		"full_name" VARCHAR(255),
		"class" INTEGER,
		"email" VARCHAR(255),
		"role" VARCHAR(255)
	  );`)

	if err != nil {
		zap.L().Fatal("Error creating schema for users", zap.Error(err))
		return err
	} else {
		zap.L().Info("Schema created for users")
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "attendances" (
		"username" VARCHAR(255),
		"attendance_id" VARCHAR(255) PRIMARY KEY,
		"punch_in_date" TIMESTAMP WITH TIME ZONE,
		"punch_out_date" TIMESTAMP WITH TIME ZONE
	  );
	  `)

	if err != nil {
		zap.L().Fatal("Error creating schema for attendance", zap.Error(err))
		return err
	} else {
		zap.L().Info("Schema created for attendance")
	}

	_, err = db.Exec(`ALTER TABLE IF EXISTS "attendances" ADD CONSTRAINT "fk_user_attendance" FOREIGN KEY ("username") REFERENCES "users" ("username");`)
	pgErr, ok := err.(pg.Error)

	if ok && pgErr.Field('C') == "42710" {
		zap.L().Warn("Foreign key already exist")
	} else if err != nil {
		zap.L().Fatal("Error creating foregin key for attendance", zap.Error(err))
		return err
	} else {
		zap.L().Info("Foreign key created")
	}

	return nil

	// models := []interface{}{
	// 	(*User)(nil),
	// 	(*Attendance)(nil),
	// }

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
	// return nil
}
