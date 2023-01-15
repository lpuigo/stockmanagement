package femovement

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
)

type ArticleFlow struct {
	*js.Object

	ArticleId int     `js:"ArticleId"`
	Price     float64 `js:"Price"`
	Amount    int     `js:"Amount"`
}

func NewArticleFlow() *ArticleFlow {
	af := &ArticleFlow{Object: tools.O()}
	af.ArticleId = -1
	af.Price = 0
	af.Amount = 0
	return af
}

func ArticleFlowFromJS(o *js.Object) *ArticleFlow {
	return &ArticleFlow{Object: o}
}
