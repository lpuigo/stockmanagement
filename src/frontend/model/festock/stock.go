package festock

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/femovement"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
)

// type Stock reflects backend/model/stock.Stock
type Stock struct {
	*js.Object

	Id         int                    `js:"Id"`
	Ref        string                 `js:"Ref"`
	Articles   []int                  `js:"Articles"`
	Movements  []*femovement.Movement `js:"Movements"`
	Quantities map[int]int            `js:"Quantities"`
	CTime      string                 `js:"CTime"`
	UTime      string                 `js:"UTime"`
	DTime      string                 `js:"DTime"`
}

func NewStock() *Stock {
	w := &Stock{Object: tools.O()}
	w.Id = -1
	w.Ref = ""
	w.Ref = ""
	w.Articles = []int{}
	w.Movements = []*femovement.Movement{}
	w.Quantities = make(map[int]int)
	w.CTime = ""
	w.UTime = ""
	w.DTime = ""
	return w
}

func StockFromJS(o *js.Object) *Stock {
	return &Stock{Object: o}
}
