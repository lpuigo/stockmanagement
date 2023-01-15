package festatus

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
)

// type Status reflects backend/model/status.Status
type Status struct {
	*js.Object
	Time   string
	Actor  string
	Status string
}

func NewStatus() *Status {
	s := &Status{Object: tools.O()}
	s.Time = ""
	s.Actor = ""
	s.Status = ""
	return s
}

func StatusFromJS(o *js.Object) *Status {
	return &Status{Object: o}
}
