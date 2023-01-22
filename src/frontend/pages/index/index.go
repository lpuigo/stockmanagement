package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/comp/adminmodal"
	"github.com/lpuig/batec/stockmanagement/src/frontend/comp/articletable"
	"github.com/lpuig/batec/stockmanagement/src/frontend/comp/userloginmodal"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle/articleconst"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/festock"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/feuser"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/elements/message"
	"github.com/lpuigo/hvue"
	"honnef.co/go/js/xhr"
	"strconv"
)

//go:generate bash ./makejs.sh

func main() {
	mpm := NewMainPageModel()

	hvue.NewVM(
		hvue.El("#app"),
		userloginmodal.RegisterComponent(),
		adminmodal.RegisterComponent(),
		articletable.RegisterComponent(),
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
	go m.callGetUser(onUnauthorized, onUserLogged)
}

func (m *MainPageModel) ShowUserLogin() {
	m.VM.Refs("UserLoginModal").Call("Show", m.User)
}

func (m *MainPageModel) ShowAdmin() {
	m.VM.Refs("AdminModal").Call("Show", m.User)
}

func (m *MainPageModel) UserLogout() {
	go m.callLogout()
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

func (m *MainPageModel) callLogout() {
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
	m.ClearSiteInfos()
	m.ClearModes()
}
