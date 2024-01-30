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

	return r
}
