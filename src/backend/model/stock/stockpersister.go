package stock

import (
	"fmt"
	"github.com/lpuig/batec/stockmanagement/src/backend/persist"
)

type StocksPersister struct {
	*persist.Persister
}

func NewStocksPersister(dir string) (*StocksPersister, error) {
	up := &StocksPersister{
		Persister: persist.NewPersister("Stocks", dir),
	}
	err := up.CheckDirectory()
	if err != nil {
		return nil, err
	}
	return up, nil
}

func (sp StocksPersister) NbStocks() int {
	return sp.NbRecords()
}

// LoadDirectory loads all persisted Stocks Records
func (sp *StocksPersister) LoadDirectory() error {
	return sp.Persister.LoadDirectory("stock", func(file string) (persist.Recorder, error) {
		return NewStockRecordFromFile(file)
	})
}

// Add adds the given StockRecord to the USersPersister and return its (updated with new id) StockRecord
func (sp *StocksPersister) Add(nar *StockRecord) *StockRecord {
	nar.SetCreateDate()
	sp.Persister.Add(nar)
	return nar
}

// Update updates the given StockRecord
func (sp *StocksPersister) Update(uar *StockRecord) error {
	uar.SetUpdateDate()
	return sp.Persister.Update(uar)
}

// Remove removes the given StockRecord from the StocksPersister (pertaining file is moved to deleted dir)
func (sp *StocksPersister) Remove(rar *StockRecord) error {
	rar.SetDeleteDate()
	return sp.Persister.Remove(rar)
}

// GetById returns the StockRecord with given Id (or nil if Id not found)
func (sp *StocksPersister) GetById(id int) *StockRecord {
	return sp.Persister.GetById(id).(*StockRecord)
}

// GetByRef returns the StockRecord with given Ref (or nil if Ref not found)
func (sp *StocksPersister) GetByRef(ref string) *StockRecord {
	for _, r := range sp.GetRecords() {
		if ar, ok := r.(*StockRecord); ok && ar.Ref == ref {
			return ar
		}
	}
	return nil
}

func (sp *StocksPersister) GetStocks() []*Stock {
	res := []*Stock{}
	for _, r := range sp.GetRecords() {
		if ar, ok := r.(*StockRecord); ok {
			res = append(res, ar.Stock)
		}
	}
	return res
}

func (sp *StocksPersister) UpdateStocks(updatedStocks []*Stock) error {
	for _, updStock := range updatedStocks {
		updStockRec := NewStockRecordFromStock(updStock)
		if updStockRec.Id < 0 { // New Stock, add it instead of update
			sp.Add(updStockRec)
			continue
		}
		err := sp.Update(updStockRec)
		if err != nil {
			fmt.Errorf("could not update stock '%s' (id: %d)", updStockRec.Ref, updStockRec.Id)
		}
	}
	return nil
}
