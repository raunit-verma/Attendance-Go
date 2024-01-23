package repository

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID   uuid.UUID
	Username string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"fullname"`
	Class    int    `json:"class"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type Attendance struct {
	UserID       uuid.UUID
	AttendanceID uuid.UUID
	PunchInDate  time.Time
	PunchOutDate time.Time
}
