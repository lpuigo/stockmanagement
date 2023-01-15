package movement

import (
	"github.com/lpuig/batec/stockmanagement/src/backend/model/status"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/timestamp"
)

type Movement struct {
	Id            int
	StockId       int
	Type          string
	Date          string
	Actor         string
	Responsible   string
	WorksiteId    int
	StatusHistory []*status.Status
	ArticleFlows  []*ArticleFlow
	timestamp.TimeStamp
}
