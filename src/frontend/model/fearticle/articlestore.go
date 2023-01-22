package fearticle

import (
	"github.com/gopherjs/gopherjs/js"
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

	Articles []*Article `js:"Articles"`
	Loaded   bool       `js:"Loaded"`

	Ref *ref.Ref `js:"Ref"`
}

func NewArticleStore() *ArticleStore {
	as := &ArticleStore{Object: tools.O()}
	as.Articles = []*Article{}
	as.Loaded = false
	as.Ref = ref.NewRef(func() string {
		return json.Stringify(as.Articles)
	})
	return as
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
	as.Loaded = true
	onSuccess()
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

	as.Ref.SetReference()
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
