package router

import (
	"attendance/api/restHandler"
	"net/http"

	"github.com/gorilla/mux"
)

type GetTeacherAttendanceHandler interface {
	GetTeacherAttendance(w http.ResponseWriter, r *http.Request)
}

type MUXRouterImpl struct {
	homeHandler  restHandler.HomeHandler
	loginHandler restHandler.LoginHandler
	// verifyToken auth.
}

func NewMUXRouterImpl(homeHandler restHandler.HomeHandler) *MUXRouterImpl {
	return &MUXRouterImpl{homeHandler: homeHandler}
}

type ServerConfig struct {
	Port string
}

func (impl *MUXRouterImpl) NewMUXRouter() *mux.Router {
	// creating a new mux router

	r := mux.NewRouter()

	// all the routes are defined here
	// home route to display all stats
	r.HandleFunc("/home", impl.homeHandler.Home).Methods("POST")

	// Route for accepting username and password
	r.HandleFunc("/login", restHandler.LoginHandler).Methods("POST")

	// Route to verify token
	r.HandleFunc("/verify", restHandler.VerifyToken).Methods("GET")

	// Route for adding new users
	r.HandleFunc("/addnewuser", restHandler.AddNewUserHandler).Methods("POST")

	// Route for Punch-in and Punch-out
	r.HandleFunc("/punchin", restHandler.PunchInHandler).Methods("GET")
	r.HandleFunc("/punchout", restHandler.PunchOutHandler).Methods("GET")

	// Route for Teacher attendence for particular month accessible by Principal and Teacher
	r.HandleFunc("/getteacherattendance", restHandler.GetTeacherAttendanceHandler).Methods("POST")

	// Route to get class attendance for day, month and year
	r.HandleFunc("/getclassattendance", restHandler.GetClassAttendanceHandler).Methods("POST")

	// Route to get particular student attendance for month and year
	r.HandleFunc("/getstudentattendance", restHandler.GetStudentAttendanceHandler).Methods("POST")

	// Route to get currentStatus of any user
	r.HandleFunc("/fetchstatus", restHandler.FetchStatus).Methods("Get")
	return r
}
