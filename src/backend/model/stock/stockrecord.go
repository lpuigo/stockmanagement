package stock

import (
	"encoding/json"
	"github.com/lpuig/batec/stockmanagement/src/backend/persist"
	"io"
	"os"
)

type StockRecord struct {
	*persist.Record
	*Stock
}

// NewStockRecord returns a new StockRecord
func NewStockRecord() *StockRecord {
	ar := &StockRecord{}
	ar.Record = persist.NewRecord(func(w io.Writer) error {
		return json.NewEncoder(w).Encode(ar.Stock)
	})
	return ar
}

// NewStockRecordFrom returns a StockRecord populated from the given reader
func NewStockRecordFrom(r io.Reader) (ar *StockRecord, err error) {
	ar = NewStockRecord()
	err = json.NewDecoder(r).Decode(ar)
	if err != nil {
		ar = nil
		return
	}
	ar.SetId(ar.Id)
	return
}

// NewStockRecordFromFile returns a StockRecord populated from the given file
func NewStockRecordFromFile(file string) (ar *StockRecord, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	ar, err = NewStockRecordFrom(f)
	if err != nil {
		ar = nil
		return
	}
	return
}

// NewStockRecordFromStock returns a StockRecord populated from given Stock
func NewStockRecordFromStock(act *Stock) *StockRecord {
	ar := NewStockRecord()
	ar.Stock = act
	ar.SetId(ar.Id)
	return ar
}

func (ar *StockRecord) SetId(id int) {
	ar.Record.SetId(id)
	ar.Id = ar.GetId()
}
