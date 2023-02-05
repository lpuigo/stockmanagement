package femovement

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
)

type ArticleFlow struct {
	*js.Object

	ArtId int     `js:"ArtId"`
	Price float64 `js:"Price"`
	Qty   int     `js:"Qty"`
}

func NewArticleFlow() *ArticleFlow {
	af := &ArticleFlow{Object: tools.O()}
	af.ArtId = -1
	af.Price = 0
	af.Qty = 0
	return af
}

func ArticleFlowFromJS(o *js.Object) *ArticleFlow {
	return &ArticleFlow{Object: o}
}
