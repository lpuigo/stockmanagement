package femovement

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/festatus"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
)

// type Movement reflects backend/model/movement.Movement
type Movement struct {
	*js.Object

	Id            int                `js:"Id"`
	StockId       int                `js:"StockId"`
	Type          string             `js:"Type"`
	Date          string             `js:"Date"`
	Actor         string             `js:"Actor"`
	Responsible   string             `js:"Responsible"`
	WorksiteId    int                `js:"WorksiteId"`
	StatusHistory []*festatus.Status `js:"StatusHistory"`
	ArticleFlows  []*ArticleFlow     `js:"ArticleFlows"`
	CTime         string             `js:"CTime"`
	UTime         string             `js:"UTime"`
	DTime         string             `js:"DTime"`
}

func NewMovement() *Movement {
	m := &Movement{Object: tools.O()}
	m.Id = -1
	m.StockId = -1
	m.Type = ""
	m.Date = ""
	m.Actor = ""
	m.Responsible = ""
	m.WorksiteId = -1
	m.StatusHistory = []*festatus.Status{}
	m.ArticleFlows = []*ArticleFlow{}
	m.CTime = ""
	m.UTime = ""
	m.DTime = ""
	return m
}

func MovementFromJS(o *js.Object) *Movement {
	return &Movement{Object: o}
}
