package stock

import (
	"github.com/lpuig/batec/stockmanagement/src/backend/model/article"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/movement"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/timestamp"
)

type Stock struct {
	Id         int
	Ref        string
	Articles   []*article.Article
	Movements  []*movement.Movement
	Quantities map[int]int
	timestamp.TimeStamp
}
