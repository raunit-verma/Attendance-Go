package router

import (
	"attendance/api/restHandler"

	"github.com/gorilla/mux"
)

type ServerConfig struct {
	Port string
}

func NewMUXRouter() *mux.Router {
	// creating a new mux router

	r := mux.NewRouter()

	// all the routes are defined here

	// Route for accepting username and password
	r.HandleFunc("/login", restHandler.LoginHandler).Methods("POST")

	// Route for adding new users
	r.HandleFunc("/addNewUser", restHandler.AddNewUserHandler).Methods("POST")

	// Route for Punch-in and Punch-out
	r.HandleFunc("/punchIn", restHandler.PunchInHandler).Methods("GET")
	r.HandleFunc("/punchOut", restHandler.PunchOutHandler).Methods("GET")

	// Route for Teacher attendence for particular month accessible by Principal and Teacher
	r.HandleFunc("/getTeacherAttendance", restHandler.GetTeacherAttendanceHandler).Methods("POST")

	// Route to get class attendance for day, month and year
	r.HandleFunc("/getClassAttendance", restHandler.GetClassAttendanceHandler).Methods("POST")

	// Route to get particular student attendance for month and year
	r.HandleFunc("/getStudentAttendance", restHandler.GetStudentsAttendanceHandler).Methods("POST")

	return r
}
