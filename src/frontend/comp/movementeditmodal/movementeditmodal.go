package movementeditmodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/femovement"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/feuser"
	"github.com/lpuigo/hvue"
)

type MovementEditModalModel struct {
	*MovementModalModel

	EditMode string `js:"EditMode"`
}

const (
	modeMovement   string = "acc"
	modeRentalStay string = "stay"
)

func NewMovementEditModalModel(vm *hvue.VM) *MovementEditModalModel {
	aemm := &MovementEditModalModel{MovementModalModel: NewMovementModalModel(vm)}
	aemm.EditMode = modeMovement
	return aemm
}

func MovementEditModalModelFromJS(o *js.Object) *MovementEditModalModel {
	return &MovementEditModalModel{MovementModalModel: MovementModalModelFromJS(o)}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Component Methods

func RegisterComponent() hvue.ComponentOption {
	return hvue.Component("movement-edit-modal", componentOptions()...)
}

func componentOptions() []hvue.ComponentOption {
	return []hvue.ComponentOption{
		hvue.Template(template),
		hvue.DataFunc(func(vm *hvue.VM) interface{} {
			return NewMovementEditModalModel(vm)
		}),
		hvue.MethodsOf(&MovementEditModalModel{}),
		hvue.Computed("hasChanged", func(vm *hvue.VM) interface{} {
			aemm := MovementEditModalModelFromJS(vm.Object)
			return aemm.HasChanged()
		}),
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Modal Methods

func (memm *MovementEditModalModel) Show(editedMvt *femovement.Movement, user *feuser.User) {
	memm.MovementModalModel.Show(editedMvt, user)
}

func (memm *MovementEditModalModel) ConfirmChange(vm *hvue.VM) {
	memm = MovementEditModalModelFromJS(vm.Object)
	memm.MovementModalModel.ConfirmChange()
	vm.Emit("edited-movement", memm.EditedMovement)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// HTML Methods
