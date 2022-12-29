package timestamp

import (
	"github.com/lpuig/batec/stockmanagement/src/backend/model/date"
	"time"
)

type TimeStamp struct {
	CTime string
	UTime string
	DTime string
}

// SetCreateDate set CreateDate for receiver, based on Today' date
func (a *TimeStamp) SetCreateDate() {
	a.CTime = date.Now().TimeStampShort()
	a.UTime = a.CTime
}

// SetUpdateDate set UpdateDate for receiver, based on Today' date
func (a *TimeStamp) SetUpdateDate() {
	a.UTime = date.Now().TimeStampShort()
}

// SetUpdateDateFrom set UpdateDate for receiver, based on given time.Time
func (a *TimeStamp) SetUpdateDateFrom(t time.Time) {
	a.UTime = date.Date(t).TimeStampShort()
}

// SetCreateDate set CreateDate for receiver, based on Today' date
func (a *TimeStamp) SetDeleteDate() {
	a.DTime = date.Now().TimeStampShort()
}
