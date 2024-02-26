package bean

import "time"

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
