package worksite

import (
	"fmt"
	"github.com/lpuig/batec/stockmanagement/src/backend/persist"
)

type WorksitesPersister struct {
	*persist.Persister
}

func NewWorksitesPersister(dir string) (*WorksitesPersister, error) {
	up := &WorksitesPersister{
		Persister: persist.NewPersister("Worksites", dir),
	}
	err := up.CheckDirectory()
	if err != nil {
		return nil, err
	}
	return up, nil
}

func (wp WorksitesPersister) NbWorksites() int {
	return wp.NbRecords()
}

// LoadDirectory loads all persisted Worksites Records
func (wp *WorksitesPersister) LoadDirectory() error {
	return wp.Persister.LoadDirectory("worksite", func(file string) (persist.Recorder, error) {
		return NewWorksiteRecordFromFile(file)
	})
}

// Add adds the given WorksiteRecord to the USersPersister and return its (updated with new id) WorksiteRecord
func (wp *WorksitesPersister) Add(nar *WorksiteRecord) *WorksiteRecord {
	nar.SetCreateDate()
	wp.Persister.Add(nar)
	return nar
}

// Update updates the given WorksiteRecord
func (wp *WorksitesPersister) Update(uar *WorksiteRecord) error {
	uar.SetUpdateDate()
	return wp.Persister.Update(uar)
}

// Remove removes the given WorksiteRecord from the WorksitesPersister (pertaining file is moved to deleted dir)
func (wp *WorksitesPersister) Remove(rar *WorksiteRecord) error {
	rar.SetDeleteDate()
	return wp.Persister.Remove(rar)
}

// GetById returns the WorksiteRecord with given Id (or nil if Id not found)
func (wp *WorksitesPersister) GetById(id int) *WorksiteRecord {
	return wp.Persister.GetById(id).(*WorksiteRecord)
}

// GetByRef returns the ArticleRecord with given Ref (or nil if Ref not found)
func (wp *WorksitesPersister) GetByRef(ref string) *WorksiteRecord {
	for _, r := range wp.GetRecords() {
		if ar, ok := r.(*WorksiteRecord); ok && ar.Ref == ref {
			return ar
		}
	}
	return nil
}

func (wp *WorksitesPersister) GetWorksites() []*Worksite {
	res := []*Worksite{}
	for _, r := range wp.GetRecords() {
		if ar, ok := r.(*WorksiteRecord); ok {
			res = append(res, ar.Worksite)
		}
	}
	return res
}

func (wp *WorksitesPersister) UpdateWorksites(updatedWorksites []*Worksite) error {
	for _, updAct := range updatedWorksites {
		updMvtRec := NewWorksiteRecordFromWorksite(updAct)
		if updMvtRec.Id < 0 { // New Worksite, add it instead of update
			wp.Add(updMvtRec)
			continue
		}
		err := wp.Update(updMvtRec)
		if err != nil {
			fmt.Errorf("could not update Worksite '%s' (id: %d)", updMvtRec.Ref, updMvtRec.Id)
		}
	}
	return nil
}
