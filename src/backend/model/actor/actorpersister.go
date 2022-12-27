package actor

import (
	"fmt"
	"github.com/lpuig/batec/stockmanagement/src/backend/persist"
)

type ActorsPersister struct {
	*persist.Persister
}

func NewActorsPersister(dir string) (*ActorsPersister, error) {
	up := &ActorsPersister{
		Persister: persist.NewPersister("Actors", dir),
	}
	err := up.CheckDirectory()
	if err != nil {
		return nil, err
	}
	return up, nil
}

func (ap ActorsPersister) NbActors() int {
	return ap.NbRecords()
}

// LoadDirectory loads all persisted Actors Records
func (ap *ActorsPersister) LoadDirectory() error {
	return ap.Persister.LoadDirectory("actor", func(file string) (persist.Recorder, error) {
		return NewActorRecordFromFile(file)
	})
}

// Add adds the given ActorRecord to the USersPersister and return its (updated with new id) ActorRecord
func (ap *ActorsPersister) Add(nar *ActorRecord) *ActorRecord {
	nar.SetCreateDate()
	ap.Persister.Add(nar)
	return nar
}

// Update updates the given ActorRecord
func (ap *ActorsPersister) Update(uar *ActorRecord) error {
	uar.SetUpdateDate()
	return ap.Persister.Update(uar)
}

// Remove removes the given ActorRecord from the ActorsPersister (pertaining file is moved to deleted dir)
func (ap *ActorsPersister) Remove(rar *ActorRecord) error {
	rar.SetDeleteDate()
	return ap.Persister.Remove(rar)
}

// GetById returns the ActorRecord with given Id (or nil if Id not found)
func (ap *ActorsPersister) GetById(id int) *ActorRecord {
	return ap.Persister.GetById(id).(*ActorRecord)
}

// GetByRef returns the ActorRecord with given Ref (or nil if Id not found)
func (ap *ActorsPersister) GetByRef(ref string) *ActorRecord {
	for _, r := range ap.GetRecords() {
		if ar, ok := r.(*ActorRecord); ok && ar.Ref == ref {
			return ar
		}
	}
	return nil
}

func (ap *ActorsPersister) GetActors() []*Actor {
	res := []*Actor{}
	for _, r := range ap.GetRecords() {
		if ar, ok := r.(*ActorRecord); ok {
			res = append(res, ar.Actor)
		}
	}
	return res
}

func (ap *ActorsPersister) UpdateActors(updatedActors []*Actor) error {
	for _, updAct := range updatedActors {
		updActRec := NewActorRecordFromActor(updAct)
		if updActRec.Id < 0 { // New Actor, add it instead of update
			ap.Add(updActRec)
			continue
		}
		err := ap.Update(updActRec)
		if err != nil {
			fmt.Errorf("could not update actor '%s' (id: %d)", updActRec.Ref, updActRec.Id)
		}
	}
	return nil
}
