package movement

import (
	"encoding/json"
	"github.com/lpuig/batec/stockmanagement/src/backend/persist"
	"io"
	"os"
)

type MovementRecord struct {
	*persist.Record
	*Movement
}

// NewMovementRecord returns a new MovementRecord
func NewMovementRecord() *MovementRecord {
	ar := &MovementRecord{}
	ar.Record = persist.NewRecord(func(w io.Writer) error {
		return json.NewEncoder(w).Encode(ar.Movement)
	})
	return ar
}

// NewMovementRecordFrom returns a MovementRecord populated from the given reader
func NewMovementRecordFrom(r io.Reader) (ar *MovementRecord, err error) {
	ar = NewMovementRecord()
	err = json.NewDecoder(r).Decode(ar)
	if err != nil {
		ar = nil
		return
	}
	ar.SetId(ar.Id)
	return
}

// NewMovementRecordFromFile returns a MovementRecord populated from the given file
func NewMovementRecordFromFile(file string) (ar *MovementRecord, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	ar, err = NewMovementRecordFrom(f)
	if err != nil {
		ar = nil
		return
	}
	return
}

// NewMovementRecordFromMovement returns a MovementRecord populated from given Movement
func NewMovementRecordFromMovement(act *Movement) *MovementRecord {
	ar := NewMovementRecord()
	ar.Movement = act
	ar.SetId(ar.Id)
	return ar
}

func (ar *MovementRecord) SetId(id int) {
	ar.Record.SetId(id)
	ar.Id = ar.GetId()
}
