package movement

import (
	"github.com/lpuig/batec/stockmanagement/src/backend/model/date"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/status"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/status/statusconst"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/timestamp"
	"github.com/lpuig/batec/stockmanagement/src/frontend/model/femovement/movementconst"
	"testing"
)

func Test_AddNewMovement(t *testing.T) {
	const (
		movementsPersisterDir string = `C:\Users\Laurent\Golang\src\github.com\lpuig\batec\stockmanagement\Ressources\Movements`
	)

	sp, err := NewMovementsPersister(movementsPersisterDir)
	if err != nil {
		t.Fatalf("NewMovementsPersister returned unexpected %s", err.Error())
	}

	err = sp.LoadDirectory()
	if err != nil {
		t.Fatalf("LoadDirectory returned unexpected %s", err.Error())
	}
	t.Logf("Loaded %d movements\n", sp.NbMovements())
	sp.NoDelay()

	for i, name := range []string{"Pierre", "Paul", "Jacques"} {
		nm := &Movement{
			Id:          -1 - i,
			StockId:     0,
			Type:        movementconst.TypeValueWithdrawal,
			Date:        date.Now().AddDays(i).String(),
			Actor:       name,
			Responsible: "Eric",
			WorksiteId:  0,
			StatusHistory: []*status.Status{
				&status.Status{
					Time:   date.Now().AddDays(i).TimeStampShort(),
					Actor:  name,
					Status: statusconst.ValueToBeValidated,
				},
			},
			ArticleFlows: []*ArticleFlow{
				&ArticleFlow{
					ArtId: i,
					Price: float64((i + 5) * 3),
					Qty:   (i + 1) * 4,
				},
				&ArticleFlow{
					ArtId: i + 1,
					Price: float64((i + 6) * 3),
					Qty:   (i + 1) * 4,
				},
			},
			TimeStamp: timestamp.TimeStamp{},
		}
		mr := sp.Add(NewMovementRecordFromMovement(nm))
		t.Logf("Add new movement Id %d\n", mr.Id)
	}
}
