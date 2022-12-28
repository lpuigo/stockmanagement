package manager

import (
	"encoding/json"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/actor"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/user"
	"io"
	"sort"
)

func (m Manager) GetUsers(writer io.Writer) error {
	usrs := m.Users.GetUsers()
	return json.NewEncoder(writer).Encode(usrs)
}

func (m Manager) UpdateUsers(updatedUsers []*user.User) error {
	return m.Users.UpdateUsers(updatedUsers)
}

// GetCurrentUserActors returns slice of actor.Actor visible by current user
//
// Rules:
//
// - any user sees all actors
func (m Manager) GetCurrentUserActors() []*actor.Actor {
	actors := m.Actors.GetActors()
	sort.Slice(actors, func(i, j int) bool {
		return actors[i].Ref < actors[j].Ref
	})
	return actors
}
