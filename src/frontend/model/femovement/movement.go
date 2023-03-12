package femovement

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/status/statusconst"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/femovement/movementconst"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/festatus"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/elements"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/fedate"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/json"
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

func CreateNewMovement(stockId int, mvtType string) *Movement {
	nm := NewMovement()
	nm.StockId = stockId
	nm.Type = mvtType
	nm.Date = fedate.TodayAfter(0)
	nm.Actor = "A Saisir"
	nm.Responsible = "A Saisir"

	return nm
}

func MovementFromJS(o *js.Object) *Movement {
	return &Movement{Object: o}
}

// Copy returns a deep copy of receiver
func (m *Movement) Copy() *Movement {
	return MovementFromJS(json.Parse(json.Stringify(m.Object)))
}

// Clone updates all receiver attributes with given Movement one's
func (m *Movement) Clone(om *Movement) {
	m.Id = om.Id
	m.StockId = om.StockId
	m.Type = om.Type
	m.Date = om.Date
	m.Actor = om.Actor
	m.Responsible = om.Responsible
	m.WorksiteId = om.WorksiteId
	sh := []*festatus.Status{}
	for _, status := range om.StatusHistory {
		sh = append(sh, status.Copy())
	}
	m.StatusHistory = sh
	afs := []*ArticleFlow{}
	for _, flow := range om.ArticleFlows {
		afs = append(afs, flow.Copy())
	}
	m.ArticleFlows = afs
	m.CTime = om.CTime
	m.UTime = om.UTime
	m.DTime = om.DTime
}

func (m *Movement) GetCurrentStatus() *festatus.Status {
	if len(m.StatusHistory) > 0 {
		return m.StatusHistory[0]
	}
	return nil
}

func (m *Movement) AddStatus(name string, validated bool) {
	s := festatus.NewStatus()
	s.Actor = name
	if validated {
		s.Status = statusconst.ValueValidated
	} else {
		s.Status = statusconst.ValueToBeValidated
	}
	m.StatusHistory = append([]*festatus.Status{s}, m.StatusHistory...)
}

func GetTypeLabel(t string) string {
	switch t {
	case movementconst.TypeValueInventory:
		return movementconst.TypeLabelInventory
	case movementconst.TypeValueSupply:
		return movementconst.TypeLabelSupply
	case movementconst.TypeValueReturn:
		return movementconst.TypeLabelReturn
	case movementconst.TypeValueWithdrawal:
		return movementconst.TypeLabelWithdrawal
	default:
		return "undefined '" + t + "'"
	}
}

func (m *Movement) GetTypeLabel() string {
	return GetTypeLabel(m.Type)
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
