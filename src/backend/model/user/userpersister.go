package user

import (
	"fmt"
	"github.com/lpuig/batec/stockmanagement/src/backend/persist"
)

type UsersPersister struct {
	*persist.Persister
}

func NewUsersPersister(dir string) (*UsersPersister, error) {
	up := &UsersPersister{
		Persister: persist.NewPersister("Users", dir),
	}
	err := up.CheckDirectory()
	if err != nil {
		return nil, err
	}
	return up, nil
}

func (up UsersPersister) NbUsers() int {
	return up.NbRecords()
}

// LoadDirectory loads all persisted Users Records
func (up *UsersPersister) LoadDirectory() error {
	return up.Persister.LoadDirectory("user", func(file string) (persist.Recorder, error) {
		return NewUserRecordFromFile(file)
	})
}

// Add adds the given UserRecord to the USersPersister and return its (updated with new id) UserRecord
func (up *UsersPersister) Add(nur *UserRecord) *UserRecord {
	nur.SetCreateDate()
	up.Persister.Add(nur)
	return nur
}

// Update updates the given UserRecord
func (up *UsersPersister) Update(uur *UserRecord) error {
	uur.SetUpdateDate()
	return up.Persister.Update(uur)
}

// Remove removes the given UserRecord from the UsersPersister (pertaining file is moved to deleted dir)
func (up *UsersPersister) Remove(rur *UserRecord) error {
	rur.SetDeleteDate()
	return up.Persister.Remove(rur)
}

// GetById returns the UserRecord with given Id (or nil if Id not found)
func (up *UsersPersister) GetById(id int) *UserRecord {
	return up.Persister.GetById(id).(*UserRecord)
}

// GetByRef returns the UserRecord with given Name (or nil if Id not found)
func (up *UsersPersister) GetByName(name string) *UserRecord {
	for _, r := range up.GetRecords() {
		if ur, ok := r.(*UserRecord); ok && ur.Name == name {
			return ur
		}
	}
	return nil
}

func (up *UsersPersister) GetUsers() []*User {
	res := []*User{}
	for _, r := range up.GetRecords() {
		if ur, ok := r.(*UserRecord); ok {
			res = append(res, ur.User)
		}
	}
	return res
}

func (up *UsersPersister) UpdateUsers(updatedUsers []*User) error {
	for _, updUsr := range updatedUsers {
		updUsrRec := NewUserRecordFromUser(updUsr)
		if updUsrRec.Id < 0 { // New User, add it instead of update
			up.Add(updUsrRec)
			continue
		}
		err := up.Update(updUsrRec)
		if err != nil {
			fmt.Errorf("could not update user '%s' (id: %d)", updUsrRec.Name, updUsrRec.Id)
		}
	}
	return nil
}
