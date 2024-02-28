package router

import (
	"attendance/api/restHandler"

	"github.com/gorilla/mux"
)

type MuxRouter interface {
	NewMUXRouter() *mux.Router
}

type MUXRouterImpl struct {
	homeHandler                 restHandler.HomeHandler
	loginHandler                restHandler.LoginHandler
	addNewUserHandler           restHandler.NewUserHandler
	punchInOutHandler           restHandler.PunchInOutHandler
	getTeacherAttendanceHandler restHandler.TeacherAttendanceHandler
	getClassAttendanceHandler   restHandler.ClassAttendanceHandler
	getStudentAttandance        restHandler.StudentAttendanceHandler
	fetchStatusHandler          restHandler.FetchStatusHandler
}

func NewMUXRouterImpl(
	homeHandler restHandler.HomeHandler,
	loginHandler restHandler.LoginHandler,
	addNewUserHandler restHandler.NewUserHandler,
	punchInOutHandler restHandler.PunchInOutHandler,
	getTeacherAttendanceHandler restHandler.TeacherAttendanceHandler,
	getClassAttendanceHandler restHandler.ClassAttendanceHandler,
	getStudentAttandance restHandler.StudentAttendanceHandler,
	fetchStatusHandler restHandler.FetchStatusHandler) *MUXRouterImpl {
	return &MUXRouterImpl{
		homeHandler:                 homeHandler,
		loginHandler:                loginHandler,
		addNewUserHandler:           addNewUserHandler,
		punchInOutHandler:           punchInOutHandler,
		getTeacherAttendanceHandler: getTeacherAttendanceHandler,
		getClassAttendanceHandler:   getClassAttendanceHandler,
		getStudentAttandance:        getStudentAttandance,
		fetchStatusHandler:          fetchStatusHandler}
}

func (impl *MUXRouterImpl) NewMUXRouter() *mux.Router {
	// creating a new mux router

	r := mux.NewRouter()

	// all the routes are defined here
	// home route to display all stats
	r.HandleFunc("/home", impl.homeHandler.Home).Methods("POST")

	// Route for accepting username and password
	r.HandleFunc("/login", impl.loginHandler.Login).Methods("POST")

	// Route for adding new users
	r.HandleFunc("/user", impl.addNewUserHandler.AddNewUser).Methods("POST")

	// Route for Punch-in and Punch-out
	r.HandleFunc("/punch/in", impl.punchInOutHandler.PunchIn).Methods("GET")
	r.HandleFunc("/punch/out", impl.punchInOutHandler.PunchOut).Methods("GET")

	// Route for Teacher attendence for particular month accessible by Principal and Teacher
	r.HandleFunc("/teacher/attendance", impl.getTeacherAttendanceHandler.GetTeacherAttendance).Methods("POST")

	// Route to get class attendance for day, month and year
	r.HandleFunc("/getclassattendance", impl.getClassAttendanceHandler.GetClassAttendance).Methods("POST")

	// Route to get particular student attendance for month and year
	r.HandleFunc("/getstudentattendance", impl.getStudentAttandance.GetStudentAttendance).Methods("POST")

	// Route to get currentStatus of any user
	r.HandleFunc("/fetchstatus", impl.fetchStatusHandler.FetchStatus).Methods("Get")
	return r
}
