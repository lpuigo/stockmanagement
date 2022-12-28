package beuser

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
)

// type BeUSer reflects backend/model/user.User struct
type BeUser struct {
	*js.Object
	Id          int             `js:"Id"`
	Name        string          `js:"Name"`
	Password    string          `js:"Password"`
	Permissions map[string]bool `js:"Permissions"`
	CTime       string          `js:"CTime"`
	UTime       string          `js:"UTime"`
	DTime       string          `js:"DTime"`
}

func BeUserFromJS(obj *js.Object) *BeUser {
	bu := &BeUser{Object: obj}
	return bu
}

func NewBeUser() *BeUser {
	usr := &BeUser{Object: tools.O()}
	usr.Id = -1
	usr.Name = ""
	usr.Password = ""
	usr.Permissions = make(map[string]bool)
	usr.CTime = ""
	usr.UTime = ""
	usr.DTime = ""
	return usr
}
