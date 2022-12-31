package manager

import (
	"encoding/json"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/stock"
	"io"
)

func (m Manager) GetStocks(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(m.Stocks.GetStocks())
}

func (m Manager) UpdateStocks(updatedStocks []*stock.Stock) error {
	return m.Stocks.UpdateStocks(updatedStocks)
}
