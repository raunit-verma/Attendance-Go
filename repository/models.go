package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
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
			log.Fatal(err)
			return err
		} else {
			fmt.Printf("Schema create for %T\n", model)
		}
	}
	return nil
}
