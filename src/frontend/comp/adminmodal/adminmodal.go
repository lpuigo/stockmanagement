package adminmodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/comp/modal"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/beuser"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/feuser"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/elements"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/elements/message"
	"github.com/lpuigo/hvue"
	"honnef.co/go/js/xhr"
	"sort"
	"strconv"
)

type AdminModalModel struct {
	*modal.ModalModel

	User       *feuser.User  `js:"user"`
	UsersStore *beuser.Store `js:"UsersStore"`
}

func NewAdminModalModel(vm *hvue.VM) *AdminModalModel {
	tpmm := &AdminModalModel{
		ModalModel: modal.NewModalModel(vm),
	}
	tpmm.User = feuser.NewUser()
	tpmm.UsersStore = beuser.NewStore()
	return tpmm
}

func AdminModalModelFromJS(o *js.Object) *AdminModalModel {
	tpmm := &AdminModalModel{
		ModalModel: &modal.ModalModel{Object: o},
	}
	return tpmm
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("admin-modal", componentOption()...)
}

func componentOption() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewAdminModalModel(vm)
		}),
		hvue.MethodsOf(&AdminModalModel{}),
		hvue.Computed("filteredUsers", func(vm *hvue.VM) interface{} {
			amm := AdminModalModelFromJS(vm.Object)
			return amm.UsersStore.Users
		}),
		hvue.Computed("hasChanged", func(vm *hvue.VM) interface{} {
			amm := AdminModalModelFromJS(vm.Object)
			userDirty := amm.UsersStore.Ref.IsDirty()
			return userDirty
		}),
	}
}

func (amm *AdminModalModel) ReloadData() {
	go amm.callReloadData()
}

func (amm *AdminModalModel) SaveArchive() {
	go amm.callSaveArchive()
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (amm *AdminModalModel) Show(user *feuser.User) {
	amm.User = user
	amm.Loading = false
	amm.UsersStore.CallGetUsers(amm.VM, func() {})
	amm.ModalModel.Show()
}

func (amm *AdminModalModel) HideWithControl(user *feuser.User) {
	if amm.UsersStore.Ref.Dirty {
		message.ConfirmCancelWarning(amm.VM, "Sauvegarder les modifications utilisateurs apportées ?",
			func() { // confirm
				amm.UsersStore.CallUpdateUsers(amm.VM, func() {
					amm.ModalModel.Hide()
				})
			},
			func() {
				amm.ModalModel.Hide()
			},
		)
	}
}

func (amm *AdminModalModel) UndoChange() {
	if amm.UsersStore.Ref.Dirty {
		amm.UsersStore.Users = amm.UsersStore.GetReferenceUsers()
	}
}

func (amm *AdminModalModel) ConfirmChange() {
	if amm.UsersStore.Ref.Dirty {
		amm.UsersStore.CallUpdateUsers(amm.VM, func() {})
	}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Users Tabs Methods

func (amm *AdminModalModel) TableRowClassName(vm *hvue.VM) string {
	return ""
}

func (amm *AdminModalModel) AddNewUser(vm *hvue.VM) {
	amm = AdminModalModelFromJS(vm.Object)
	nuser := beuser.NewBeUser()
	nuser.Name = "Nouvel Utilisateur"
	nuser.Password = "default"
	amm.UsersStore.AddNewUser(nuser)
}

// Column Filtering Related Methods

func (amm *AdminModalModel) FilterHandler(vm *hvue.VM, value string, p *js.Object, col *js.Object) bool {
	prop := col.Get("property").String()
	return p.Get(prop).String() == value
}

func (amm *AdminModalModel) FilterList(vm *hvue.VM, prop string) []*elements.ValText {
	amm = AdminModalModelFromJS(vm.Object)
	count := map[string]int{}
	attribs := []string{}

	var translate func(string) string
	translate = func(val string) string { return val }

	for _, usr := range amm.UsersStore.Users {
		var attrs []string
		switch prop {
		default:
			attrs = []string{usr.Object.Get(prop).String()}
		}
		for _, a := range attrs {
			if _, exist := count[a]; !exist {
				attribs = append(attribs, a)
			}
			count[a]++
		}
	}
	sort.Strings(attribs)
	res := []*elements.ValText{}
	for _, a := range attribs {
		fa := a
		if fa == "" {
			fa = "Vide"
		}
		res = append(res, elements.NewValText(a, translate(fa)+" ("+strconv.Itoa(count[a])+")"))
	}
	return res
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// WS call Methods

func (amm *AdminModalModel) callReloadData() {
	defer func() { amm.Loading = false }()
	req := xhr.NewRequest("GET", "/api/admin/reload")
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(amm.VM, "Oups! "+err.Error(), true)
		amm.Hide()
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(amm.VM, req)
		amm.Hide()
		return
	}
	message.NotifySuccess(amm.VM, "Données", "Rechargement des données effectué")
	amm.VM.Emit("reload")
	return
}

func (amm *AdminModalModel) callSaveArchive() {
	defer func() { amm.Loading = false }()
	req := xhr.NewRequest("GET", "/api/archive")
	req.Timeout = tools.TimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(amm.VM, "Oups! "+err.Error(), true)
		amm.Hide()
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(amm.VM, req)
		amm.Hide()
		return
	}
	message.NotifySuccess(amm.VM, "Archivage", "Sauvegarde des archives demandée")
	return
}
