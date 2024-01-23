package restHandler

import (
	auth "attendance/services"
	"encoding/json"
	"fmt"
	"net/http"
)

type NewUser struct {
	UserName string `json:"username"`
	Password string `json:"password"`
	FullName string `json:"fullname"`
	Class    int    `json:"class"`
	Email    string `json:"email"`
	Type     string `json:"type"`
}

func (newUser NewUser) VerifyNewUserData() bool {
	if newUser.UserName != "" && newUser.Password != "" && newUser.FullName != "" && newUser.Email != "" && (newUser.Type == "teacher" || newUser.Type == "student") && newUser.Class != 0 {
		return true
	}
	return false
}

func AddNewUserHandler(w http.ResponseWriter, r *http.Request) {

	status, username := auth.VerifyToken(r)

	if status != http.StatusAccepted {
		w.WriteHeader(status)
		return
	}

	// verify is user is principal using the username
	fmt.Println(username)

	newUser := NewUser{}
	err := json.NewDecoder(r.Body).Decode(&newUser)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// checking if all data is present

	w.WriteHeader(http.StatusAccepted)

	return
}
