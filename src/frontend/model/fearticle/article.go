package fearticle

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/fearticle/articleconst"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/elements"
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

func (a *Article) SearchString(filter string) string {
	searchItem := func(prefix, typ, value string) string {
		if value == "" {
			return ""
		}
		if filter != articleconst.FilterLabelAll && filter != typ {
			return ""
		}
		return prefix + typ + value
	}

	res := searchItem("", articleconst.FilterValueDes, a.Designation)
	res += searchItem("", articleconst.FilterValueRef, a.Ref)
	res += searchItem("", articleconst.FilterValueCat, a.Category)
	return res
}

func GetFilterTypeValueLabel() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(articleconst.FilterValueAll, articleconst.FilterLabelAll),
		elements.NewValueLabel(articleconst.FilterValueDes, articleconst.FilterLabelDes),
		elements.NewValueLabel(articleconst.FilterValueRef, articleconst.FilterLabelRef),
		elements.NewValueLabel(articleconst.FilterValueCat, articleconst.FilterLabelCat),
	}
}

func GetStatusLabel(status string) string {
	switch status {
	case articleconst.StatusValueAvailable:
		return articleconst.StatusLabelAvailable
	case articleconst.StatusValueUnavailable:
		return articleconst.StatusLabelUnavailable
	default:
		return articleconst.StatusLabelError
	}
}
