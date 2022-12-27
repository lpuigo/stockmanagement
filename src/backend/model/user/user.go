package user

import "github.com/lpuig/batec/stockmanagement/src/backend/model/timestamp"

type User struct {
	Id          int
	Name        string
	Password    string
	Permissions map[string]bool
	timestamp.TimeStamp
}

func NewUser(name string) *User {
	u := &User{
		Id:          0,
		Name:        name,
		Password:    "",
		Permissions: make(map[string]bool),
	}
	u.SetCreateDate()
	return u
}

func (u *User) HasPermissionHR() bool {
	return u.Permissions["HR"]
}

func (u *User) HasPermissionUpdate() bool {
	return u.Permissions["Update"]
}
