package feworksite

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
)

// type Worksite reflects backend/model/worksite.Worksite
type Worksite struct {
	*js.Object

	Id          int    `js:"Id"`
	Client      string `js:"Client"`
	City        string `js:"City"`
	Ref         string `js:"Ref"`
	Responsible string `js:"Responsible"`
	DateBegin   string `js:"DateBegin"`
	DateEnd     string `js:"DateEnd"`
	CTime       string `js:"CTime"`
	UTime       string `js:"UTime"`
	DTime       string `js:"DTime"`
}

func NewWorksite() *Worksite {
	w := &Worksite{Object: tools.O()}
	w.Id = -1
	w.Client = ""
	w.City = ""
	w.Ref = ""
	w.Responsible = ""
	w.DateBegin = ""
	w.DateEnd = ""
	w.CTime = ""
	w.UTime = ""
	w.DTime = ""
	return w
}

func WorksiteFromJS(o *js.Object) *Worksite {
	return &Worksite{Object: o}
}

func (w Worksite) IsActiveOn(today string) bool {
	return today >= w.DateBegin && today <= w.DateEnd
}

func (w Worksite) GetLabel() string {
	if w.Id == -1 {
		return "Non défini"
	}
	return w.Client + " - " + w.City + " / " + w.Ref
}
