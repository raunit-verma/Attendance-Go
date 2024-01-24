package repository

func CheckIfUserExists(username string) *User {
	db := GetDB()
	user := User{
		Username: username,
	}
	err := db.Model(user).WherePK().Select()
	if err != nil {
		return nil
	}
	return &user
}

func AddNewUser(user User) error {
	db := GetDB()
	_, err := db.Model(user).Insert()
	if err != nil {
		return err
	}
	return nil
}
