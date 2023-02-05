package manager

import (
	"encoding/json"
	"github.com/lpuig/batec/stockmanagement/src/backend/model/movement"
	"io"
)

func (m Manager) GetMovements(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(m.Movements.GetMovements())
}

func (m Manager) GetMovementsForStockId(writer io.Writer, sid int) error {
	return json.NewEncoder(writer).Encode(m.Movements.GetMovementsForStockId(sid))
}

func (m Manager) UpdateMovements(updatedMovements []*movement.Movement) error {
	return m.Movements.UpdateMovements(updatedMovements)
}
