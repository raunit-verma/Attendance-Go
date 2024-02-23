package repository

import (
	"attendance/util"
	"fmt"
	"time"

	"github.com/go-pg/pg"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Repository interface {
	GetUser(username string) *User
	AddNewUser(user *User) error
	GetCurrentStatus(username string, startDate string, endDate string) (bool, []Attendance)
	AddNewPunchIn(username string) error
	AddNewPunchOut(username string, attendance Attendance) error
	GetTeacherAttendance(username string, data GetTeacherAttendanceJSON) []Attendance
	GetClassAttendance(data GetClassAttendanceJSON) []StudentAttendanceJSON
	GetStudentAttendance(username string, data GetStudentAttendanceJSON) []Attendance
	GetDailyStats(data GetHomeJSON) (int, int, int, int)
}

type RepositoryImpl struct {
	db *pg.DB
}

func NewRepositoryImpl(db *pg.DB) *RepositoryImpl {
	return &RepositoryImpl{db: db}
}

func (impl *RepositoryImpl) GetUser(username string) *User {
	user := &User{}

	err := impl.db.Model(user).Where("username=?", username).Select()

	if err != nil {
		zap.L().Info("No record found in DB", zap.String("username", username))
		return nil
	}
	return user
}

func (impl *RepositoryImpl) AddNewUser(user *User) error {

	_, err := impl.db.Model(user).Insert()

	if err != nil {
		zap.L().Error("Error in inserting new user.", zap.Error(err))
		return err
	}

	return nil
}

func (impl *RepositoryImpl) GetCurrentStatus(username string, startDate string, endDate string) (bool, []Attendance) {

	var attendances []Attendance
	err := impl.db.Model(&attendances).
		Where("username = ?", username).
		Where("punch_in_date BETWEEN ? AND ?", startDate, endDate).
		Where("punch_out_date IS NULL").
		Select()
	if err != nil {
		zap.L().Error("Error in DB operation", zap.Error(err))
		return false, nil
	}
	if len(attendances) == 0 {
		return false, nil
	}
	return true, attendances
}

func (impl *RepositoryImpl) AddNewPunchIn(username string) error {
	t := time.Now()
	_, currentTime := util.FormateDateTime(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	attendance := Attendance{
		AttendanceID: uuid.New().String(),
		PunchInDate:  currentTime,
		Username:     username,
	}
	_, err := impl.db.Model(&attendance).Insert()
	if err != nil {
		zap.L().Error("Cannot add new punch in of user "+username, zap.Error(err))
		return err
	}
	return nil
}

func (impl *RepositoryImpl) AddNewPunchOut(username string, attendance Attendance) error {
	t := time.Now()
	_, currentTime := util.FormateDateTime(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	attendance.PunchOutDate = currentTime
	_, err := impl.db.Model(&attendance).Where("attendance_id = ?", attendance.AttendanceID).Column("punch_out_date").Update()
	if err != nil {
		zap.L().Error("Cannot add new punch out of user "+username, zap.Error(err))
		return err
	}
	return nil
}

func (impl *RepositoryImpl) GetTeacherAttendance(username string, data GetTeacherAttendanceJSON) []Attendance {
	var attendances []Attendance

	startDate, _ := util.FormateDateTime(data.Year, time.Month(data.Month), 1, 0, 0, 0)
	endDate, _ := util.FormateDateTime(data.Year, time.Month(data.Month), 31, 23, 59, 59)

	err := impl.db.Model(&attendances).Where("username=?", username).Where("punch_in_date BETWEEN ? AND ?", startDate, endDate).Select()

	if err != nil {
		zap.L().Error("Cannot retrieve data for teacher "+username, zap.Error(err))
		return nil
	}

	return attendances
}

func (impl *RepositoryImpl) GetClassAttendance(data GetClassAttendanceJSON) []StudentAttendanceJSON {
	var results []StudentAttendanceJSON

	startDate, _ := util.FormateDateTime(data.Year, time.Month(data.Month), data.Day, 0, 0, 0)
	endDate, _ := util.FormateDateTime(data.Year, time.Month(data.Month), data.Day, 23, 59, 59)

	err := impl.db.Model(&results).
		ColumnExpr("DISTINCT users.username").
		Column("users.full_name").
		Join("JOIN attendances a ON users.username = a.username").
		Table("users").
		Where("users.class = ?", data.Class).
		Where("a.punch_in_date BETWEEN ? AND ?", startDate, endDate).
		Where("users.role=?", "student").
		Select()

	if err != nil {
		fmt.Println(err)
		zap.L().Error("Error in retrieving particular class attendance ", zap.Error(err))
		return nil
	}
	return results
}

func (impl *RepositoryImpl) GetStudentAttendance(username string, data GetStudentAttendanceJSON) []Attendance {
	var results []Attendance
	startDate, _ := util.FormateDateTime(data.Year, time.Month(data.Month), 1, 0, 0, 0)
	endDate, _ := util.FormateDateTime(data.Year, time.Month(data.Month), 31, 23, 59, 59)

	err := impl.db.Model(&results).Where("username=?", username).Where("punch_in_date BETWEEN ? AND ?", startDate, endDate).Select()

	if err != nil {
		zap.L().Error("Error in retrieving particular student attendance "+username, zap.Error(err))
		return nil
	}

	return results
}

func (impl *RepositoryImpl) GetDailyStats(data GetHomeJSON) (int, int, int, int) {
	startDate, _ := util.FormateDateTime(data.Year, time.Month(data.Month), data.Date, 0, 0, 0)
	endDate, _ := util.FormateDateTime(data.Year, time.Month(data.Month), data.Date, 23, 59, 59)
	var totalTeacherPresent int
	var totalStudentPresent int
	var totalStudent int
	var totalTeacher int
	err := impl.db.Model(&Attendance{}).
		ColumnExpr("COUNT (DISTINCT attendances.username)").
		Join("JOIN users ON users.username = attendances.username").
		Table("attendances").
		Where("attendances.punch_in_date BETWEEN ? AND ?", startDate, endDate).
		Where("users.role=?", "teacher").
		Select(&totalTeacherPresent)

	if err != nil {
		zap.L().Error("Error in counting distinct present username teacher", zap.Error(err))
		return -1, -1, -1, -1
	}

	err = impl.db.Model(&Attendance{}).
		ColumnExpr("COUNT (DISTINCT attendances.username)").
		Join("JOIN users ON users.username = attendances.username").
		Table("attendances").
		Where("attendances.punch_in_date BETWEEN ? AND ?", startDate, endDate).
		Where("users.role=?", "student").
		Select(&totalStudentPresent)

	if err != nil {
		zap.L().Error("Error in counting distinct present username student", zap.Error(err))
		return -1, -1, -1, -1
	}

	err = impl.db.Model(&User{}).
		ColumnExpr("COUNT (DISTINCT username)").
		Where("role=?", "student").Select(&totalStudent)

	if err != nil {
		zap.L().Error("Error in counting distinct username student", zap.Error(err))
		return -1, -1, -1, -1
	}

	err = impl.db.Model(&User{}).
		ColumnExpr("COUNT (DISTINCT username)").
		Where("role=?", "teacher").Select(&totalTeacher)

	if err != nil {
		zap.L().Error("Error in counting distinct username teacher", zap.Error(err))
		return -1, -1, -1, -1
	}

	return totalStudentPresent, totalTeacherPresent, totalStudent, totalTeacher
}
