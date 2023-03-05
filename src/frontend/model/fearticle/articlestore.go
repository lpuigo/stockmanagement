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

	Articles     []*Article       `js:"Articles"`
	Loaded       bool             `js:"Loaded"`
	ArticleIndex map[int]*Article `js:"ArticleIndex"`

	Ref       *ref.Ref `js:"Ref"`
	NextNewId int      `js:"NextNewId"`
}

func NewArticleStore() *ArticleStore {
	as := &ArticleStore{Object: tools.O()}
	as.Articles = []*Article{}
	as.Loaded = false
	as.Ref = ref.NewRef(func() string {
		return json.Stringify(as.Articles)
	})
	as.NextNewId = 0
	as.ArticleIndex = make(map[int]*Article)
	return as
}

// IsDirty returns true if Loading in progress or Ref is dirty
func (as *ArticleStore) IsDirty() bool {
	return !(as.Loaded && !as.Ref.IsDirty())
}

func (as *ArticleStore) GetNextNewId() int {
	as.NextNewId--
	return as.NextNewId
}

// ArticleSliceFromJS returns a slice of Article extracted from the given JS Array of Article
func ArticleSliceFromJS(jsArticleArray *js.Object) []*Article {
	articles := []*Article{}
	jsArticleArray.Call("forEach", func(item *js.Object) {
		a := ArticleFromJS(item)
		articles = append(articles, a)
	})
	return articles
}

// UpdateArticleIndex sets receiver's ArticleIndex
func (as *ArticleStore) UpdateArticleIndex() {
	dict := make(map[int]*Article)
	for _, article := range as.Articles {
		dict[article.Id] = article
	}
	as.ArticleIndex = dict
}

func (as *ArticleStore) SetArticles(arts []*Article) {
	as.Articles = arts
	as.UpdateArticleIndex()
	as.Ref.SetReference()
	as.Loaded = true
}

// UpdateWith updates receiver Articles by checking given articles
//
// - updatedArticle with new or unknown id is added
//
// - updatedArticle with known Id will update already existing article
func (as *ArticleStore) UpdateWith(updatedArticles []*Article) {
	for _, updatedArticle := range updatedArticles {
		matchingArticle, found := as.ArticleIndex[updatedArticle.Id]
		if !found {
			// unknown article ... add it
			updatedArticle.Id = as.GetNextNewId()
			as.Articles = append(as.Articles, updatedArticle)
			continue
		}
		matchingArticle.Clone(updatedArticle)
	}
	as.UpdateArticleIndex()
}

func (as *ArticleStore) CallGetArticles(vm *hvue.VM, onSuccess func()) {
	as.Loaded = false
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
	as.SetArticles(ArticleSliceFromJS(req.Response))
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
	as.Loaded = false
	go as.callUpdateArticles(vm, onSuccess)
}

func (as *ArticleStore) callUpdateArticles(vm *hvue.VM, onSuccess func()) {
	req := xhr.NewRequest("PUT", "/api/articles")
	req.Timeout = tools.LongTimeOut
	req.ResponseType = xhr.JSON

	defer func() {
		as.Loaded = true
	}()

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

func (as *ArticleStore) GetExportArticlestoXlsxURL() string {
	return "/api/articles/export"
}

func (as *ArticleStore) GetImportArticlesFromXlsxURL() string {
	return "/api/articles/import"
}
