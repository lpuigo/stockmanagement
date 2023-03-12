package worksite

import (
	"github.com/lpuig/batec/stockmanagement/src/backend/model/timestamp"
)

type Worksite struct {
	Id          int
	Client      string
	City        string
	Ref         string
	Responsible string
	DateBegin   string
	DateEnd     string
	timestamp.TimeStamp
}
