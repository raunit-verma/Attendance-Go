package restHandler

import (
	"attendance/bean"

	"go.uber.org/zap"
)

func ValidateClassRequestData(data bean.GetClassAttendanceJSON) (bool, string) {
	if data.Class <= 0 || data.Class > 12 {
		zap.L().Info("Requested class is not valid")
		return true, "Class is not valid. "
	} else if data.Month <= 0 || data.Month > 12 {
		zap.L().Info("Requested month is not valid")
		return true, "Month is not valid. "
	} else if data.Year <= 2020 || data.Year >= 2100 {
		zap.L().Info("Request year is not valid")
		return true, "Year is not valid. "
	} else if data.Day <= 0 || data.Day > 31 {
		zap.L().Info("Requested day is not valid")
		return true, "Day is not valid. "
	}
	return false, ""
}

func ValidateHomeRequestData(data bean.GetHomeJSON) (bool, string) {
	if data.Month <= 0 || data.Month > 12 {
		zap.L().Info("Requested month is not valid")
		return true, "Month is not valid. "
	} else if data.Year <= 2020 || data.Year >= 2100 {
		zap.L().Info("Request year is not valid")
		return true, "Year is not valid. "
	}
	return false, ""
}

func ValidateStudentRequestData(data bean.GetStudentAttendanceJSON) (bool, string) {
	if data.Month <= 0 || data.Month > 12 {
		zap.L().Info("Requested month is not valid")
		return true, "Month is not valid. "
	} else if data.Year <= 2020 || data.Year >= 2100 {
		zap.L().Info("Request year is not valid")
		return true, "Year is not valid. "
	}
	return false, ""
}

func ValidateTeacherRequestData(data bean.GetTeacherAttendanceJSON) (bool, string) {
	if data.ID == "" {
		zap.L().Info("Teacher id is null")
		return true, "Teacher id is null. "
	} else if data.Month <= 0 || data.Month > 12 {
		zap.L().Info("Month is not valid")
		return true, "Month is not valid. "
	} else if data.Year < 2020 || data.Year >= 2100 {
		zap.L().Info("Year is not valid")
		return true, "Year is not valid. "
	}
	return false, ""
}
