package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/comp/movementtable"
	"github.com/lpuig/batec/stockmanagement/src/frontend/comp/stockarticletable"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/femovement"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/festock"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/feuser"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/elements/message"
	"github.com/lpuigo/hvue"
	"honnef.co/go/js/xhr"
)

//go:generate bash ./makejs.sh

func main() {
	mpm := NewMainPageModel()

	hvue.NewVM(
		hvue.El("#stock_app"),
		stockarticletable.RegisterComponent(),
		movementtable.RegisterComponent(),
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
			return mpm.StockStore.Ref.IsDirty()
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
	m.SaveInProgress = false

	m.ActiveMode = "article"
	m.Filter = ""
	m.FilterType = ""

	return m
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Action Methods

func (m *MainPageModel) InitPage(vm *hvue.VM) {
	m.LoadStock(vm)
	onLoadedArticles := func() {}
	m.AvailableArticles.CallGetArticles(vm, onLoadedArticles)
}

func (m *MainPageModel) LoadStock(vm *hvue.VM) {
	m = &MainPageModel{Object: vm.Object}
	sid := tools.GetURLSearchParam("sid")
	if sid == nil {
		message.ErrorMsgStr(m.VM, "Identifiant de stock non trouvé", sid, true)
		return
	}
	stockId := sid.Int()
	onLoaded := func() {
		if len(m.StockStore.Stocks) > 0 {
			m.Stock = m.StockStore.Stocks[0]
			js.Global.Get("document").Set("title", m.Stock.Ref)
		}
		onLoadedMovements := func() {}
		m.MovementStore.CallGetMovementsForStockId(vm, stockId, onLoadedMovements)
	}
	m.StockStore.CallGetStockById(m.VM, stockId, onLoaded)
}

func (m *MainPageModel) SaveStock(vm *hvue.VM) {
	m = &MainPageModel{Object: vm.Object}
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
	go m.callGetUser(onUnauthorized, onUserLogged)
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// WS call Methods

func (m *MainPageModel) callGetUser(notloggedCallback, loggedCallback func()) {
	req := xhr.NewRequest("GET", "/api/login")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status == tools.HttpUnauthorized {
		notloggedCallback()
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(m.VM, req)
		return
	}
	m.User.Copy(feuser.UserFromJS(req.Response))
	if m.User.Name == "" {
		m.User = feuser.NewUser()
		return
	}
	m.User.Connected = true
	loggedCallback()
}

func (m *MainPageModel) callLogout(callBack func()) {
	req := xhr.NewRequest("DELETE", "/api/login")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON
	err := req.Send(nil)
	if err != nil {
		message.ErrorStr(m.VM, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(m.VM, req)
		return
	}
	m.User = feuser.NewUser()
	m.User.Connected = false
	callBack()
}
