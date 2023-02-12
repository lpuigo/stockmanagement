package movementeditmodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/femovement"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/festock"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/feuser"
	"github.com/lpuigo/hvue"
)

type MovementEditModalModel struct {
	*MovementModalModel

	EditMode string `js:"EditMode"`

	StockArticles *fearticle.ArticleStore `js:"articles"`
	Stock         *festock.Stock          `js:"stock"`
}

const (
	modeMovement   string = "acc"
	modeRentalStay string = "stay"
)

func NewMovementEditModalModel(vm *hvue.VM) *MovementEditModalModel {
	aemm := &MovementEditModalModel{MovementModalModel: NewMovementModalModel(vm)}
	aemm.EditMode = modeMovement
	aemm.StockArticles = fearticle.NewArticleStore()
	aemm.Stock = festock.NewStock()
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
		hvue.Props("stock", "articles"),
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

func (memm *MovementEditModalModel) FormatType(t string) string {
	return femovement.GetTypeLabel(t)
}

func (memm *MovementEditModalModel) UpdateDate(vm *hvue.VM) {
	//memm = MovementEditModalModelFromJS(vm.Object)
}
