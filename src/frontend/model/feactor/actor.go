package feactor

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools/fedate"
)

// type Actor reflects backend/model/actor.Actor
type Actor struct {
	*js.Object

	Id        int               `js:"Id"`
	Ref       string            `js:"Ref"`
	FirstName string            `js:"FirstName"`
	LastName  string            `js:"LastName"`
	State     string            `js:"State"`
	Period    *fedate.DateRange `js:"Period"`
	Company   string            `js:"Company"`
	Comment   string            `js:"Comment"`
	CTime     string            `js:"CTime"`
	UTime     string            `js:"UTime"`
	DTime     string            `js:"DTime"`
}

func NewActor() *Actor {
	w := &Actor{Object: tools.O()}
	w.Id = -1
	w.Ref = ""
	w.FirstName = ""
	w.LastName = ""
	w.State = ""
	w.Period = fedate.NewDateRange()
	w.Company = ""
	w.Comment = ""
	w.CTime = ""
	w.UTime = ""
	w.DTime = ""
	return w
}

func ActorFromJS(o *js.Object) *Actor {
	return &Actor{Object: o}
}
