package fearticle

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle/articleconst"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/festock"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/ref"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/elements/message"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/json"
	"github.com/lpuigo/hvue"
	"honnef.co/go/js/xhr"
	"strconv"
)

type ArticleStore struct {
	*js.Object

	Articles []*Article         `js:"Articles"`
	Loaded   bool               `js:"Loaded"`
	GetById  func(int) *Article `js:"GetById"`

	Ref *ref.Ref `js:"Ref"`
}

func NewArticleStore() *ArticleStore {
	as := &ArticleStore{Object: tools.O()}
	as.Articles = []*Article{}
	as.Loaded = false
	as.Ref = ref.NewRef(func() string {
		return json.Stringify(as.Articles)
	})
	as.GetById = func(int) *Article { return nil }
	return as
}

// GenGetById returns a GetByArticleId func, which, given an article Id, returns the pertaining article if exists, or null
func (as *ArticleStore) GenGetById() func(id int) *Article {
	dict := make(map[int]*Article)
	for _, article := range as.Articles {
		dict[article.Id] = article
	}
	return func(id int) *Article {
		return dict[id]
	}
}

func (as *ArticleStore) CallGetArticles(vm *hvue.VM, onSuccess func()) {
	go as.callGetArticles(vm, onSuccess)
}

func (as *ArticleStore) callGetArticles(vm *hvue.VM, onSuccess func()) {
	req := xhr.NewRequest("GET", "/api/articles")
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
	loadedArticles := []*Article{}
	req.Response.Call("forEach", func(item *js.Object) {
		a := ArticleFromJS(item)
		loadedArticles = append(loadedArticles, a)
	})
	as.Articles = loadedArticles
	as.Ref.SetReference()
	as.GetById = as.GenGetById()
	as.Loaded = true
	onSuccess()
}

// SetArticleStatusFromStock sets Article status depending on it is declared in stock or not
func (as *ArticleStore) SetArticleStatusFromStock(stock *festock.Stock) {
	isArticleInStockById := stock.GetArticleAvailability()
	for _, art := range as.Articles {
		if isArticleInStockById[art.Id] {
			art.Status = articleconst.StatusValueOutOfStock
			continue
		}
		art.Status = articleconst.StatusValueUnavailable
	}
}

func (as *ArticleStore) CallUpdateArticles(vm *hvue.VM, onSuccess func()) {
	go as.callUpdateArticles(vm, onSuccess)
}

func (as *ArticleStore) callUpdateArticles(vm *hvue.VM, onSuccess func()) {
	req := xhr.NewRequest("PUT", "/api/articles")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON

	toUpdates := as.getUpdatedArticles()
	nbToUpd := len(toUpdates)
	if nbToUpd == 0 {
		onSuccess()
		return
	}

	err := req.Send(json.Stringify(toUpdates))
	if err != nil {
		message.ErrorStr(vm, "Oups! "+err.Error(), true)
		return
	}
	if req.Status != tools.HttpOK {
		message.ErrorRequestMessage(vm, req)
		return
	}

	msg := " article mis à jour"
	if nbToUpd > 1 {
		msg = " articles mis à jour"
	}
	message.NotifySuccess(vm, "Sauvegarde des articles", strconv.Itoa(nbToUpd)+msg)
	onSuccess()
}

func (as *ArticleStore) getUpdatedArticles() []*Article {
	refArticles := []*Article{}
	json.Parse(as.Ref.Reference).Call("forEach", func(acc *Article) {
		refArticles = append(refArticles, acc)
	})
	refDict := makeDictArticles(refArticles)

	updtArticles := []*Article{}
	for _, article := range as.Articles {
		refAcc := refDict[article.Id]
		if !(refAcc != nil && json.Stringify(article) == json.Stringify(refAcc)) {
			// this article has been updated ...
			updtArticles = append(updtArticles, article)
		}
	}
	return updtArticles
}

func makeDictArticles(accs []*Article) map[int]*Article {
	res := make(map[int]*Article)
	for _, acc := range accs {
		res[acc.Id] = acc
	}
	return res
}
