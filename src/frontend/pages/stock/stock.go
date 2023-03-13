package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/comp/movementupdate"
	"github.com/lpuig/batec/stockmanagement/src/frontend/comp/stockarticletable"
	"github.com/lpuig/batec/stockmanagement/src/frontend/comp/stockcatalogtable"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/femovement"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/festock"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/feuser"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/feworksite"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/elements/message"
	"github.com/lpuigo/hvue"
)

//go:generate bash ./makejs.sh

func main() {
	mpm := NewMainPageModel()

	hvue.NewVM(
		hvue.El("#stock_app"),
		stockarticletable.RegisterComponent(),
		stockcatalogtable.RegisterComponent(),
		movementupdate.RegisterComponent(),
		hvue.DataS(mpm),
		hvue.MethodsOf(mpm),
		hvue.Mounted(func(vm *hvue.VM) {
			mpm := &MainPageModel{Object: vm.Object}
			mpm.GetUserSession(func() {
				mpm.InitPage(vm)
			})
		}),
		hvue.Computed("Title", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			if mpm.Stock.Object == js.Undefined {
				return ""
			}
			return "BATEC " + mpm.Stock.Ref
		}),
		hvue.Computed("LoggedUser", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			if mpm.User.Name == "" {
				return "Non connecté"
			}
			return mpm.User.Name
		}),
		hvue.Computed("IsDirty", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			return !(!mpm.StockStore.Ref.IsDirty() && !mpm.MovementStore.Ref.IsDirty())
		}),
	)

	// TODO to remove after debug
	js.Global.Set("mpm", mpm)
}

type MainPageModel struct {
	*js.Object

	VM   *hvue.VM     `js:"VM"`
	User *feuser.User `js:"User"`

	AvailableArticles *fearticle.ArticleStore   `js:"AvailableArticles"`
	StockStore        *festock.StockStore       `js:"StockStore"`
	Stock             *festock.Stock            `js:"Stock"`
	MovementStore     *femovement.MovementStore `js:"MovementStore"`
	Worksites         *feworksite.WorksiteStore `js:"Worksites"`
	SaveInProgress    bool                      `js:"SaveInProgress"`

	ActiveMode string `js:"ActiveMode"`
	Filter     string `js:"Filter"`
	FilterType string `js:"FilterType"`
}

func NewMainPageModel() *MainPageModel {
	m := &MainPageModel{Object: tools.O()}
	m.VM = nil
	m.User = feuser.NewUser()

	m.AvailableArticles = fearticle.NewArticleStore()
	m.StockStore = festock.NewStockStore()
	m.Stock = festock.NewStock()
	m.MovementStore = femovement.NewMovementStore()
	m.Worksites = feworksite.NewWorksiteStore()
	m.SaveInProgress = false

	m.ActiveMode = "article"
	m.Filter = ""
	m.FilterType = ""

	return m
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Action Methods

func (m *MainPageModel) InitPage(vm *hvue.VM) {
	m = &MainPageModel{Object: vm.Object}
	sid := tools.GetURLSearchParam("sid")
	if sid == nil {
		message.ErrorMsgStr(m.VM, "Identifiant de stock non trouvé", sid, true)
		return
	}
	stockId := sid.Int()
	m.LoadStockWithId(stockId)
}

func (m *MainPageModel) LoadStockWithId(id int) {
	onStockLoaded := func() {
		// update page title
		if len(m.StockStore.Stocks) > 0 {
			m.Stock = m.StockStore.Stocks[0]
			js.Global.Get("document").Set("title", m.Stock.Ref)
		}
		m.LoadMovementWithStockId(id)
	}
	m.StockStore.CallGetStockById(m.VM, id, onStockLoaded)

	onWorkistesLoaded := func() {}
	m.Worksites.CallGetWorksites(m.VM, onWorkistesLoaded)
}

func (m *MainPageModel) LoadMovementWithStockId(id int) {
	// load pertaining movement
	onLoadedMovements := func() {
		m.LoadArticles()
	}
	m.MovementStore.CallGetMovementsForStockId(m.VM, id, onLoadedMovements)

}

func (m *MainPageModel) LoadArticles() {
	onLoadedArticles := func() {
		// Set Article status : available in stock or not
		m.AvailableArticles.SetArticleStatusFromStock(m.Stock)
	}
	m.AvailableArticles.CallGetArticles(m.VM, onLoadedArticles)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// HTML Methods

func (m *MainPageModel) LoadStock(vm *hvue.VM) {
	m = &MainPageModel{Object: vm.Object}
	m.LoadStockWithId(m.Stock.Id)
}

func (m *MainPageModel) SaveStock(vm *hvue.VM) {
	m = &MainPageModel{Object: vm.Object}
	if m.StockStore.Ref.IsDirty() {
		m.SaveInProgress = true
		onUpdated := func() {
			if len(m.StockStore.Stocks) > 0 {
				m.Stock = m.StockStore.Stocks[0]
				js.Global.Get("document").Set("title", m.Stock.Ref)
			}
			m.SaveInProgress = false
		}
		m.StockStore.CallUpdateStocks(m.VM, onUpdated)
	}
	if m.MovementStore.Ref.IsDirty() {
		m.SaveInProgress = true
		onUpdated := func() {
			m.SaveInProgress = false
		}
		m.MovementStore.CallUpdateMovements(m.VM, onUpdated)
	}
}

// SwitchActiveMode handles ActiveMode change
func (m *MainPageModel) SwitchActiveMode(vm *hvue.VM) {
	//m = &MainPageModel{Object: vm.Object}
	//switch m.ActiveMode {
	//case "article":
	//	// do something specific
	//case "stock":
	//	// do something specific
	//}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// User Management Methods

func (m *MainPageModel) GetUserSession(callback func()) {
	onUnauthorized := func() {}
	onUserLogged := func() {
		callback()
	}
	m.User.CallGetUser(m.VM, onUnauthorized, onUserLogged)
}
