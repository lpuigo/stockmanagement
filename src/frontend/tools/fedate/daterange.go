package date

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lpuig/batec/stockmanagement/src/frontend/tools"
)

// Type DateRange reflects stockmanagement/src/backend/model/date.DateRange
type DateRange struct {
	*js.Object

	Begin string `js:"Begin"`
	End   string `js:"End"`
}

func NewDateRange() *DateRange {
	dr := &DateRange{Object: tools.O()}
	dr.Begin = ""
	dr.End = ""
	return dr
}

func NewDateRangeFrom(beg, end string) *DateRange {
	dr := &DateRange{Object: tools.O()}
	dr.Begin = beg
	dr.End = end
	return dr
}

func (dr *DateRange) Overlap(odr *DateRange) bool {
	switch {
	case dr.Begin > odr.End:
		return false
	case dr.End < odr.Begin:
		return false
	default:
		return true
	}
}
