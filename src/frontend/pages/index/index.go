package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/comp/adminmodal"
	"github.com/lpuig/batec/stockmanagement/src/frontend/comp/articlescatalog"
	"github.com/lpuig/batec/stockmanagement/src/frontend/comp/userloginmodal"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle/articleconst"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/festock"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/feuser"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
	"github.com/lpuigo/hvue"
	"strconv"
)

//go:generate bash ./makejs.sh

func main() {
	mpm := NewMainPageModel()

	hvue.NewVM(
		hvue.El("#app"),
		userloginmodal.RegisterComponent(),
		adminmodal.RegisterComponent(),
		articlescatalog.RegisterComponent(),
		hvue.DataS(mpm),
		hvue.MethodsOf(mpm),
		hvue.Mounted(func(vm *hvue.VM) {
			mpm := &MainPageModel{Object: vm.Object}
			mpm.GetUserSession()
		}),
		hvue.Computed("LoggedUser", func(vm *hvue.VM) interface{} {
			mpm := &MainPageModel{Object: vm.Object}
			if mpm.User.Name == "" {
				return "Non connecté"
			}
			return mpm.User.Name
		}),
	)

	// TODO to remove after debug
	js.Global.Set("mpm", mpm)
}

type MainPageModel struct {
	*js.Object

	VM *hvue.VM `js:"VM"`

	User *feuser.User `js:"User"`

	AvailableArticles *fearticle.ArticleStore `js:"AvailableArticles"`
	AvailableStocks   *festock.StockStore     `js:"AvailableStocks"`

	ActiveMode string `js:"ActiveMode"`

	Filter     string `js:"Filter"`
	FilterType string `js:"FilterType"`
}

func NewMainPageModel() *MainPageModel {
	mpm := &MainPageModel{Object: tools.O()}
	mpm.VM = nil
	mpm.User = feuser.NewUser()
	mpm.AvailableArticles = fearticle.NewArticleStore()
	mpm.AvailableStocks = festock.NewStockStore()
	mpm.ClearModes()
	mpm.ClearSiteInfos()
	//mpm.SetMode()

	return mpm
}

func (m *MainPageModel) ClearModes() {
	m.ActiveMode = "stock"
	m.Filter = ""
	m.FilterType = articleconst.FilterValueAll
}

func (m *MainPageModel) ClearSiteInfos() {
	//m.WorksiteInfos = []*fm.WorksiteInfo{}
}

///////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Action Methods

func (m *MainPageModel) GetUserSession() {
	onUnauthorized := func() {
		m.ShowUserLogin()
	}
	onUserLogged := func() {
		m.GetInfos()
	}
	m.User.CallGetUser(m.VM, onUnauthorized, onUserLogged)
}

func (m *MainPageModel) ShowUserLogin() {
	m.VM.Refs("UserLoginModal").Call("Show", m.User)
}

func (m *MainPageModel) ShowAdmin() {
	m.VM.Refs("AdminModal").Call("Show", m.User)
}

func (m *MainPageModel) UserLogout() {
	m.User.CallLogout(m.VM, func() {})
}

func (m *MainPageModel) GetInfos() {
	onSuccessArticles := func() {}
	m.AvailableArticles.CallGetArticles(m.VM, onSuccessArticles)
	onSuccessStocks := func() {}
	m.AvailableStocks.CallGetStocks(m.VM, onSuccessStocks)
}

// OpenStockPage opens another web page for stock management
func (m *MainPageModel) OpenStockPage(stock *festock.Stock) {
	pageUrl := "stock.html?sid=" + strconv.Itoa(stock.Id)
	js.Global.Get("window").Call("open", pageUrl)
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
