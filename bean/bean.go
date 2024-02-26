package bean

import (
	"attendance/util"
	"net/http"
	"time"

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

type GetHomeJSON struct {
	Date  int `json:"date"`
	Month int `json:"month"`
	Year  int `json:"year"`
}

type DashboardJSON struct {
	MonthlyAttendance []bool `json:"monthly_attendance"`
	Hour              int    `json:"hour"`
	Minute            int    `json:"minute"`
	Second            int    `json:"second"`
}

func (newUser User) IsNewUserDataMissing() (int, bool, ErrorJSON) {
	IsDataMissing := false
	Message := ""

	if newUser.Username == "" {
		IsDataMissing = true
		Message = " Username is missing."
		zap.L().Info("Username is empty. ")
	} else if newUser.Password == "" {
		IsDataMissing = true
		zap.L().Info("Password is empty")
		Message = " Password is missing. "
	} else if newUser.FullName == "" {
		IsDataMissing = true
		zap.L().Info("Fullname is empty")
		Message = " Fullname is missing. "
	} else if newUser.Class <= 0 || newUser.Class > 12 {
		IsDataMissing = true
		zap.L().Info("Class constraint failed")
		Message = " Class should be between 1 to 12. "
	} else if newUser.Email != "" && !util.IsValidEmail(newUser.Email) {
		IsDataMissing = true
		zap.L().Info("Not a valid email")
		Message = " Email is missing or not a valid email. "
	} else if newUser.Role != "teacher" && newUser.Role != "student" {
		IsDataMissing = true
		zap.L().Info("Not a valid role")
		Message = " Role is missing. "
	}

	if IsDataMissing {
		return http.StatusBadRequest, IsDataMissing, ErrorJSON{ErrorCode: 3, Message: Message + util.UserDataMissing_Three}
	}

	return http.StatusAccepted, IsDataMissing, ErrorJSON{}
}
