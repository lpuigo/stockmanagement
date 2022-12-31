package movement

import (
	"fmt"
	"github.com/lpuig/batec/stockmanagement/src/backend/persist"
)

type MovementsPersister struct {
	*persist.Persister
}

func NewMovementsPersister(dir string) (*MovementsPersister, error) {
	up := &MovementsPersister{
		Persister: persist.NewPersister("Movements", dir),
	}
	err := up.CheckDirectory()
	if err != nil {
		return nil, err
	}
	return up, nil
}

func (ap MovementsPersister) NbMovements() int {
	return ap.NbRecords()
}

// LoadDirectory loads all persisted Movements Records
func (ap *MovementsPersister) LoadDirectory() error {
	return ap.Persister.LoadDirectory("movement", func(file string) (persist.Recorder, error) {
		return NewMovementRecordFromFile(file)
	})
}

// Add adds the given MovementRecord to the USersPersister and return its (updated with new id) MovementRecord
func (ap *MovementsPersister) Add(nar *MovementRecord) *MovementRecord {
	nar.SetCreateDate()
	ap.Persister.Add(nar)
	return nar
}

// Update updates the given MovementRecord
func (ap *MovementsPersister) Update(uar *MovementRecord) error {
	uar.SetUpdateDate()
	return ap.Persister.Update(uar)
}

// Remove removes the given MovementRecord from the MovementsPersister (pertaining file is moved to deleted dir)
func (ap *MovementsPersister) Remove(rar *MovementRecord) error {
	rar.SetDeleteDate()
	return ap.Persister.Remove(rar)
}

// GetById returns the MovementRecord with given Id (or nil if Id not found)
func (ap *MovementsPersister) GetById(id int) *MovementRecord {
	return ap.Persister.GetById(id).(*MovementRecord)
}

func (ap *MovementsPersister) GetMovements() []*Movement {
	res := []*Movement{}
	for _, r := range ap.GetRecords() {
		if ar, ok := r.(*MovementRecord); ok {
			res = append(res, ar.Movement)
		}
	}
	return res
}

func (ap *MovementsPersister) UpdateMovements(updatedMovements []*Movement) error {
	for _, updAct := range updatedMovements {
		updMvtRec := NewMovementRecordFromMovement(updAct)
		if updMvtRec.Id < 0 { // New Movement, add it instead of update
			ap.Add(updMvtRec)
			continue
		}
		err := ap.Update(updMvtRec)
		if err != nil {
			fmt.Errorf("could not update movement '%s' (id: %d)", updMvtRec.Date, updMvtRec.Id)
		}
	}
	return nil
}
