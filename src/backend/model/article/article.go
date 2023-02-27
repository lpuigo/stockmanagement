package article

import "github.com/lpuig/batec/stockmanagement/src/backend/model/timestamp"

type Article struct {
	Id           int
	Category     string
	SubCategory  string
	Designation  string
	Ref          string
	Manufacturer string
	PhotoId      string
	Location     string

	InvoiceUnit          string
	InvoiceUnitPrice     float64
	InvoiceUnitRetailQty float64
	RetailUnit           string
	RetailUnitStockQty   float64
	StockUnit            string

	Status string
	timestamp.TimeStamp
}
