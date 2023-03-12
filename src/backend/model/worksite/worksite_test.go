package worksite

import (
	"github.com/lpuig/batec/stockmanagement/src/backend/model/timestamp"
	"testing"
)

func Test_AddNewWorksite(t *testing.T) {
	const (
		worksitesPersisterDir string = `C:\Users\Laurent\Golang\src\github.com\lpuig\batec\stockmanagement\Ressources\Worksites`
	)

	wp, err := NewWorksitesPersister(worksitesPersisterDir)
	if err != nil {
		t.Fatalf("NewWorksitesPersister returned unexpected %s", err.Error())
	}

	err = wp.LoadDirectory()
	if err != nil {
		t.Fatalf("LoadDirectory returned unexpected %s", err.Error())
	}
	t.Logf("Loaded %d worksites\n", wp.NbWorksites())
	wp.NoDelay()

	for i, ws := range []Worksite{
		{
			Client:      "OPHLM",
			City:        "Nancy",
			Ref:         "Chantier fini",
			Responsible: "Raoul",
			DateBegin:   "2022-01-01",
			DateEnd:     "2023-02-28",
		},
		{
			Client:      "Immeuble Dupont",
			City:        "Nancy",
			Ref:         "Lotissement 1",
			Responsible: "Julien",
			DateBegin:   "2023-01-01",
			DateEnd:     "2023-05-31",
		},
		{
			Client:      "Immeuble Dupont",
			City:        "Custine",
			Ref:         "Lotissement 2",
			Responsible: "Maurice",
			DateBegin:   "2023-03-01",
			DateEnd:     "2023-07-31",
		},
		{
			Client:      "OPHLM",
			City:        "Nancy",
			Ref:         "51 rue de la pompe",
			Responsible: "Julien",
			DateBegin:   "2023-02-01",
			DateEnd:     "2023-06-30",
		},
	} {
		nm := &Worksite{
			Id:          -1 - i,
			Client:      ws.Client,
			City:        ws.City,
			Ref:         ws.Ref,
			Responsible: ws.Responsible,
			DateBegin:   ws.DateBegin,
			DateEnd:     ws.DateEnd,
			TimeStamp:   timestamp.TimeStamp{},
		}
		mr := wp.Add(NewWorksiteRecordFromWorksite(nm))
		t.Logf("Add new worksite Id %d\n", mr.Id)
	}
}
