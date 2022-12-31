package worksite

import (
	"github.com/lpuig/batec/stockmanagement/src/backend/model/status"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/timestamp"
)

type Worksite struct {
	Id             int
	Client         string
	Ref            string
	DateBegin      string
	DateEnd        string
	Status_history []*status.Status
	timestamp.TimeStamp
}
