package festatus

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/status/statusconst"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
)

// type Status reflects backend/model/status.Status
type Status struct {
	*js.Object
	Time   string `js:"Time"`
	Actor  string `js:"Actor"`
	Status string `js:"Status"`
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

func GetLabel(s string) string {
	switch s {
	case statusconst.ValueToBeValidated:
		return statusconst.LabelToBeValidated
	case statusconst.ValueValidated:
		return statusconst.LabelValidated
	default:
		return "Undefined '" + s + "'"
	}

}

func (s *Status) GetLabel() string {
	if s == nil {
		return "No Status"
	}
	return GetLabel(s.Status)
}
