package article

import "github.com/lpuig/batec/stockmanagement/src/backend/model/timestamp"

type Article struct {
	Id                  int
	Category            string
	SubCategory         string
	Designation         string
	Ref                 string
	Manufacturer        string
	PhotoId             string
	Location            string
	UnitStock           string
	UnitAccounting      string
	ConvStockAccounting float64
	Status              string
	timestamp.TimeStamp
}
