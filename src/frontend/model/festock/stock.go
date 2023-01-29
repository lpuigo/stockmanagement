package festock

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/femovement"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
	"sort"
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
	s := &Stock{Object: tools.O()}
	s.Id = -1
	s.Ref = ""
	s.Ref = ""
	s.Articles = []int{}
	s.Movements = []*femovement.Movement{}
	s.Quantities = make(map[int]int)
	s.CTime = ""
	s.UTime = ""
	s.DTime = ""
	return s
}

func StockFromJS(o *js.Object) *Stock {
	return &Stock{Object: o}
}

// GetArticleAvailability returns a map with article id as key. If id exist (value = true), article is avaliable in stock receiver
func (s *Stock) GetArticleAvailability() map[int]bool {
	dict := make(map[int]bool)
	for _, articleId := range s.Articles {
		dict[articleId] = true
	}
	return dict
}

// UpdateArticleAvailability updates receiver Articles with given articleInStock map keys
func (s *Stock) UpdateArticleAvailability(articleInStock map[int]bool) {
	artIds := []int{}
	for id, _ := range articleInStock {
		artIds = append(artIds, id)
	}
	sort.Ints(artIds)

	// s.Articles = artIds // => causes JS to use  int32Array instead of Array
	// thus causing issue with JSON unmarshalling on GO Back-End side
	// this method uses a workaround to force attribute Articles as an Array
	res := []interface{}{}
	for _, id := range artIds {
		res = append(res, id)
	}
	s.Set("Articles", res)
}
