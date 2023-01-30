package femovement

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/femovement/movementconst"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/festatus"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/elements"
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

func (m *Movement) GetCurrentStatus() *festatus.Status {
	if len(m.StatusHistory) > 0 {
		return m.StatusHistory[0]
	}
	return nil
}

func (a *Movement) SearchString(filter string) string {
	searchItem := func(prefix, typ, value string) string {
		if value == "" {
			return ""
		}
		if filter != movementconst.FilterValueAll && filter != typ {
			return ""
		}
		return prefix + typ + value
	}

	res := searchItem("", movementconst.FilterValueType, a.Type)
	res += searchItem("", movementconst.FilterValueActor, a.Actor)
	res += searchItem("", movementconst.FilterValueResponsible, a.Responsible)
	res += searchItem("", movementconst.FilterValueStatus, a.GetCurrentStatus().GetLabel())
	//res += searchItem("", movementconst.FilterValueArticle, a.Category)
	return res
}

func GetFilterTypeValueLabel() []*elements.ValueLabel {
	return []*elements.ValueLabel{
		elements.NewValueLabel(movementconst.FilterValueAll, movementconst.FilterLabelAll),
		elements.NewValueLabel(movementconst.FilterValueType, movementconst.FilterLabelType),
		elements.NewValueLabel(movementconst.FilterValueActor, movementconst.FilterLabelActor),
		elements.NewValueLabel(movementconst.FilterValueResponsible, movementconst.FilterLabelResponsible),
		elements.NewValueLabel(movementconst.FilterValueStatus, movementconst.FilterLabelStatus),
		//elements.NewValueLabel(movementconst.FilterValueArticle, movementconst.FilterLabelArticle),
	}
}
