package manager

import (
	"encoding/json"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/worksite"
	"io"
)

func (m Manager) GetWorksites(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(m.Worksites.GetWorksites())
}

func (m Manager) UpdateWorksites(updatedWorksites []*worksite.Worksite) error {
	return m.Worksites.UpdateWorksites(updatedWorksites)
}
