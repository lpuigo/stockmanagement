package stock

import (
	"github.com/lpuig/batec/stockmanagement/src/backend/model/movement"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/timestamp"
	"testing"
)

func Test_AddNEwStock(t *testing.T) {
	const (
		stocksPersisterDir string = `C:\Users\Laurent\Golang\src\github.com\lpuig\batec\stockmanagement\Ressources\Stocks`
	)

	sp, err := NewStocksPersister(stocksPersisterDir)
	if err != nil {
		t.Fatalf("NewStocksPersister returned unexpected %s", err.Error())
	}

	err = sp.LoadDirectory()
	if err != nil {
		t.Fatalf("LoadDirectory returned unexpected %s", err.Error())
	}
	t.Logf("Loaded %d stocks\n", sp.NbStocks())
	sp.NoDelay()

	ns := &Stock{
		Id:         -1,
		Ref:        "Stock secondaire",
		Articles:   []int{},
		Movements:  []*movement.Movement{},
		Quantities: make(map[int]int),
		TimeStamp:  timestamp.TimeStamp{},
	}
	sr := sp.Add(NewStockRecordFromStock(ns))
	t.Logf("Add new stock Id %d\n", sr.Id)
}
