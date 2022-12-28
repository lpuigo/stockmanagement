package manager

import (
	"encoding/json"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/actor"
	"io"
)

func (m Manager) GetActors(writer io.Writer) error {
	actors := m.GetCurrentUserActors()
	return json.NewEncoder(writer).Encode(actors)
}

func (m Manager) UpdateActors(updatedActors []*actor.Actor) error {
	return m.Actors.UpdateActors(updatedActors)
}
