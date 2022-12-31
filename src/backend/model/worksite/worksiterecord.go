package worksite

import (
	"encoding/json"
	"github.com/lpuig/batec/stockmanagement/src/backend/persist"
	"io"
	"os"
)

type WorksiteRecord struct {
	*persist.Record
	*Worksite
}

// NewWorksiteRecord returns a new WorksiteRecord
func NewWorksiteRecord() *WorksiteRecord {
	ar := &WorksiteRecord{}
	ar.Record = persist.NewRecord(func(w io.Writer) error {
		return json.NewEncoder(w).Encode(ar.Worksite)
	})
	return ar
}

// NewWorksiteRecordFrom returns a WorksiteRecord populated from the given reader
func NewWorksiteRecordFrom(r io.Reader) (wr *WorksiteRecord, err error) {
	wr = NewWorksiteRecord()
	err = json.NewDecoder(r).Decode(wr)
	if err != nil {
		wr = nil
		return
	}
	wr.SetId(wr.Id)
	return
}

// NewWorksiteRecordFromFile returns a WorksiteRecord populated from the given file
func NewWorksiteRecordFromFile(file string) (wr *WorksiteRecord, err error) {
	f, err := os.Open(file)
	if err != nil {
		return
	}
	defer f.Close()
	wr, err = NewWorksiteRecordFrom(f)
	if err != nil {
		wr = nil
		return
	}
	return
}

// NewWorksiteRecordFromWorksite returns a WorksiteRecord populated from given Worksite
func NewWorksiteRecordFromWorksite(ws *Worksite) *WorksiteRecord {
	wr := NewWorksiteRecord()
	wr.Worksite = ws
	wr.SetId(wr.Id)
	return wr
}

func (ws *WorksiteRecord) SetId(id int) {
	ws.Record.SetId(id)
	ws.Id = ws.GetId()
}
