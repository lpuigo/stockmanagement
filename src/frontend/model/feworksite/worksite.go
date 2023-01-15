package feworksite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/festatus"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
)

// type Worksite reflects backend/model/worksite.Worksite
type Worksite struct {
	*js.Object

	Id             int                `js:"Id"`
	Client         string             `js:"Client"`
	Ref            string             `js:"Ref"`
	DateBegin      string             `js:"DateBegin"`
	DateEnd        string             `js:"DateEnd"`
	Status_history []*festatus.Status `js:"Status_history"`
	CTime          string             `js:"CTime"`
	UTime          string             `js:"UTime"`
	DTime          string             `js:"DTime"`
}

func NewWorksite() *Worksite {
	w := &Worksite{Object: tools.O()}
	w.Id = -1
	w.Client = ""
	w.Ref = ""
	w.DateBegin = ""
	w.DateEnd = ""
	w.Status_history = []*festatus.Status{}
	w.CTime = ""
	w.UTime = ""
	w.DTime = ""
	return w
}

func WorksiteFromJS(o *js.Object) *Worksite {
	return &Worksite{Object: o}
}
