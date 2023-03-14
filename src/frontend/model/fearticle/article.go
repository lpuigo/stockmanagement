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

	Id           int    `js:"Id"`
	Category     string `js:"Category"`
	SubCategory  string `js:"SubCategory"`
	Designation  string `js:"Designation"`
	Ref          string `js:"Ref"`
	Manufacturer string `js:"Manufacturer"`
	PhotoId      string `js:"PhotoId"`
	Location     string `js:"Location"`

	InvoiceUnit          string  `js:"InvoiceUnit"`
	InvoiceUnitPrice     float64 `js:"InvoiceUnitPrice"`
	InvoiceUnitRetailQty float64 `js:"InvoiceUnitRetailQty"`
	RetailUnit           string  `js:"RetailUnit"`
	RetailUnitStockQty   float64 `js:"RetailUnitStockQty"`
	StockUnit            string  `js:"StockUnit"`

	Status string `js:"Status"`
	CTime  string `js:"CTime"`
	UTime  string `js:"UTime"`
	DTime  string `js:"DTime"`
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
	a.InvoiceUnit = articleconst.UnitPiece
	a.InvoiceUnitPrice = 0
	a.InvoiceUnitRetailQty = 1
	a.RetailUnit = articleconst.UnitPiece
	a.RetailUnitStockQty = 1
	a.StockUnit = articleconst.UnitPiece
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
		if filter != articleconst.FilterValueAll && filter != typ {
			return ""
		}
		return prefix + typ + value
	}

	res := searchItem("", articleconst.FilterValueDes, a.Designation)
	res += searchItem("", articleconst.FilterValueRef, a.Ref)
	res += searchItem("", articleconst.FilterValueCat, a.Category)
	return res
}

func (a *Article) Clone(oa *Article) {
	//a.Id = -1
	a.Category = oa.Category
	a.SubCategory = oa.SubCategory
	a.Designation = oa.Designation
	a.Ref = oa.Ref
	a.Manufacturer = oa.Manufacturer
	a.PhotoId = oa.PhotoId
	a.Location = oa.Location
	a.InvoiceUnit = oa.InvoiceUnit
	a.InvoiceUnitPrice = oa.InvoiceUnitPrice
	a.InvoiceUnitRetailQty = oa.InvoiceUnitRetailQty
	a.RetailUnit = oa.RetailUnit
	a.RetailUnitStockQty = oa.RetailUnitStockQty
	a.StockUnit = oa.StockUnit
	//a.Status = oa.Status
	//a.CTime = ""
	//a.UTime = ""
	//a.DTime = ""
}

func (a *Article) GetInvoicePrice(qty int) float64 {
	return float64(qty) * a.InvoiceUnitPrice
}

func (a *Article) GetRetailPrice(qty int) float64 {
	return a.GetInvoicePrice(qty) * a.InvoiceUnitRetailQty
}

func (a *Article) GetStockPrice(qty int) float64 {
	return a.GetRetailPrice(qty) * a.RetailUnitStockQty
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
	case articleconst.StatusValueOutOfStock:
		return articleconst.StatusLabelOutOfStock
	case articleconst.StatusValueUnavailable:
		return articleconst.StatusLabelUnavailable
	default:
		return articleconst.StatusLabelError
	}
}

func GetStatusClass(status string) string {
	switch status {
	case articleconst.StatusValueAvailable:
		return articleconst.StatusClassAvailable
	case articleconst.StatusValueOutOfStock:
		return articleconst.StatusClassOutOfStock
	case articleconst.StatusValueUnavailable:
		return articleconst.StatusClassUnavailable
	default:
		return articleconst.StatusClassError
	}
}

func (a *Article) ToggleInStock() {
	switch a.Status {
	// Avaliable article can not be removed from stock
	//case articleconst.StatusValueAvailable:
	//	a.Status = articleconst.StatusValueUnavailable
	case articleconst.StatusValueUnavailable:
		a.Status = articleconst.StatusValueOutOfStock
	case articleconst.StatusValueOutOfStock:
		a.Status = articleconst.StatusValueUnavailable
	}
}
