package feuser

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/elements/message"
	"github.com/lpuigo/hvue"
	"honnef.co/go/js/xhr"
)

// reflects and extend github.com\lpuigo\batec\stockmanagement\src\backend\route\session.go route.authentUser struct
type User struct {
	*js.Object

	Name        string          `js:"Name"`
	Pwd         string          `js:"Pwd"`
	Connected   bool            `js:"Connected"`
	Permissions map[string]bool `js:"Permissions"`
}

func NewUser() *User {
	user := &User{Object: tools.O()}
	user.Name = ""
	user.Pwd = ""
	user.Connected = false
	user.Permissions = make(map[string]bool)
	return user
}

func UserFromJS(o *js.Object) *User {
	return &User{Object: o}
}

// Copy copies all given User's attributes on receiver's
func (u *User) Copy(ou *User) {
	u.Name = ou.Name
	u.Pwd = ou.Pwd
	u.Connected = ou.Connected
	p := make(map[string]bool)
	for perm, value := range ou.Permissions {
		p[perm] = value
	}
	u.Permissions = p
}

func (u *User) HasPermissionInvoice() bool {
	return u.Permissions["Invoice"]
}

func (u *User) HasPermissionHR() bool {
	return u.Permissions["HR"]
}

func (u *User) HasPermissionUpdate() bool {
	return u.Permissions["Update"]
}

// CallGetUser calls the server to request the connected User.
func (u *User) CallGetUser(vm *hvue.VM, notloggedCallback, loggedCallback func()) {
	go u.callGetUser(vm, notloggedCallback, loggedCallback)
}

func (u *User) callGetUser(vm *hvue.VM, notloggedCallback, loggedCallback func()) {
	req := xhr.NewRequest("GET", "/api/login")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(vm, "Oups! "+err.Error(), true)
		return
	}
	if req.Status == tools.HttpUnauthorized {
		notloggedCallback()
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(vm, req)
		return
	}
	newUser := UserFromJS(req.Response)

	if newUser.Name == "" {
		notloggedCallback()
		return
	}
	u.Copy(newUser)
	u.Connected = true
	loggedCallback()
}

func (u *User) callLogout(vm *hvue.VM, callBack func()) {
	req := xhr.NewRequest("DELETE", "/api/login")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(vm, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(vm, req)
		return
	}
	u.Copy(NewUser())
	callBack()
}
