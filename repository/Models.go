package repository

import (
	"attendance/bean"
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

func CreateSchema(db *pg.DB, cfg bean.DBConfig) error {

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

	query := `insert into users values ('user','` + cfg.PrincipalPassword + `','Principal',1,'ramverma@gmail.com','principal')`

	_, err = db.Exec(query)

	if err != nil {
		if pgErr, ok := err.(pg.Error); ok && pgErr.IntegrityViolation() {
			zap.L().Info("Principal already exists.")
		} else {
			zap.L().Error("Error creating principal in users table", zap.Error(err))
		}
	} else {
		zap.L().Info("Pricipal created in users")
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS "attendances" (
		"username" VARCHAR(255),
		"attendance_id" VARCHAR(255) PRIMARY KEY,
		"punch_in_date" TIMESTAMP,
		"punch_out_date" TIMESTAMP
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
		zap.L().Info("Foreign key already exist")
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
