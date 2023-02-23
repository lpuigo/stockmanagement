package movementeditmodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/femovement"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/feuser"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/elements/message"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/fedate"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/json"
	"github.com/lpuigo/hvue"
)

type MovementModalModel struct {
	*js.Object

	Visible         bool                 `js:"Visible"`
	VM              *hvue.VM             `js:"VM"`
	User            *feuser.User         `js:"User"`
	EditedMovement  *femovement.Movement `js:"edited_movement"`
	CurrentMovement *femovement.Movement `js:"current_movement"`

	ShowConfirmDelete bool `js:"ShowConfirmDelete"`
}

func NewMovementModalModel(vm *hvue.VM) *MovementModalModel {
	mmm := &MovementModalModel{Object: tools.O()}
	mmm.Visible = false
	mmm.VM = vm
	mmm.User = feuser.NewUser()
	mmm.EditedMovement = femovement.NewMovement()
	mmm.CurrentMovement = femovement.NewMovement()
	mmm.ShowConfirmDelete = false

	return mmm
}

func MovementModalModelFromJS(o *js.Object) *MovementModalModel {
	return &MovementModalModel{Object: o}
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (mmm *MovementModalModel) HasChanged() bool {
	if mmm.EditedMovement.Object == js.Undefined {
		return true
	}
	return json.Stringify(mmm.CurrentMovement) != json.Stringify(mmm.EditedMovement)
}

func (mmm *MovementModalModel) Show(editedMvt *femovement.Movement, user *feuser.User) {
	mmm.EditedMovement = editedMvt
	mmm.CurrentMovement = editedMvt.Copy()
	mmm.User = user
	mmm.ShowConfirmDelete = false
	mmm.Visible = true
}

func (mmm *MovementModalModel) HideWithControl(onConfirm func()) {
	exitFn := func() {
		onConfirm()
		mmm.Hide()
	}
	if mmm.HasChanged() {
		message.ConfirmWarning(mmm.VM, "OK pour perdre les changements effectués ?", exitFn)
		return
	}
	exitFn()
}

func (mmm *MovementModalModel) Hide() {
	mmm.Visible = false
	mmm.ShowConfirmDelete = false
}

//////////////////////////////////////////////////////////////////////////////////////////////
// Action Button Methods

func (mmm *MovementModalModel) ConfirmChange() {
	mmm.EditedMovement.Clone(mmm.CurrentMovement)
	//mmm.EditedMovement.UpdateState()
	mmm.Hide()
}

func (mmm *MovementModalModel) UndoChange() {
	mmm.CurrentMovement.Clone(mmm.EditedMovement)
}

func (mmm *MovementModalModel) FormatDate(d string) string {
	return fedate.DateString(d)
}
