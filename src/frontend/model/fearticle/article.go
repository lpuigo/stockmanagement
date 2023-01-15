package fearticle

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
)

// type Article reflects backend/model/article.Article
type Article struct {
	*js.Object

	Id                  int     `js:"Id"`
	Category            string  `js:"Category"`
	SubCategory         string  `js:"SubCategory"`
	Designation         string  `js:"Designation"`
	Ref                 string  `js:"Ref"`
	Manufacturer        string  `js:"Manufacturer"`
	PhotoId             string  `js:"PhotoId"`
	Location            string  `js:"Location"`
	UnitStock           string  `js:"UnitStock"`
	UnitAccounting      string  `js:"UnitAccounting"`
	ConvStockAccounting float64 `js:"ConvStockAccounting"`
	Status              string  `js:"Status"`
	CTime               string  `js:"CTime"`
	UTime               string  `js:"UTime"`
	DTime               string  `js:"DTime"`
}

func NewArticle() *Article {
	a := &Article{Object: tools.O()}
	a.Id = -1
	a.Category = ""
	a.SubCategory = ""
	a.Designation = ""
	a.Ref = ""
	a.Manufacturer = ""
	a.PhotoId = ""
	a.Location = ""
	a.UnitStock = ""
	a.UnitAccounting = ""
	a.ConvStockAccounting = 0
	a.Status = ""
	a.CTime = ""
	a.UTime = ""
	a.DTime = ""
	return a
}

func ArticleFromJS(o *js.Object) *Article {
	return &Article{Object: o}
}
