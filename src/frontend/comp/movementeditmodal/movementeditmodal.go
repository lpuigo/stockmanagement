package movementeditmodal

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/comp/articleflowtable"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle/articleconst"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/femovement"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/festock"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/feuser"
	"github.com/lpuigo/hvue"
)

type MovementEditModalModel struct {
	*MovementModalModel

	EditMode string `js:"EditMode"`

	Articles *fearticle.ArticleStore `js:"articles"`
	Stock    *festock.Stock          `js:"stock"`

	StockArticles *fearticle.ArticleStore `js:"StockArticles"`
}

const (
	modeMovement   string = "acc"
	modeRentalStay string = "stay"
)

func NewMovementEditModalModel(vm *hvue.VM) *MovementEditModalModel {
	aemm := &MovementEditModalModel{MovementModalModel: NewMovementModalModel(vm)}
	aemm.EditMode = modeMovement
	aemm.Articles = fearticle.NewArticleStore()
	aemm.Stock = festock.NewStock()
	aemm.StockArticles = fearticle.NewArticleStore()
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
		articleflowtable.RegisterComponent(),
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
	memm.SetStockArticles()
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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// inner Methods

// SetStockArticles sets StockArticles store with articles that are defined in attached stock
func (memm *MovementEditModalModel) SetStockArticles() {
	memm.StockArticles = fearticle.NewArticleStore()
	stockArticles := []*fearticle.Article{}
	for _, article := range memm.Articles.Articles {
		if article.Status == articleconst.StatusValueUnavailable {
			continue
		}
		stockArticles = append(stockArticles, article)
	}
	memm.StockArticles.SetArticles(stockArticles)
}
