package persist

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Record struct {
	id       int
	dirty    bool
	marshall func(w io.Writer) error
}

// GetId returns the inner record id
func (r Record) GetId() int {
	return r.id
}

// SetId sets the inner record id
func (r *Record) SetId(id int) {
	r.id = id
}

// Dirty marks the record as dirty (need to be persisted in a file)
func (r *Record) Dirty() {
	r.dirty = true
}

// Persist writes receiver to its named file within the given path
func (r *Record) Persist(path string) error {
	file := r.GetFilePath(path)
	f, err := os.Create(file)
	defer f.Close()
	if err != nil {
		return err
	}
	err = r.marshall(f)
	if err != nil {
		return fmt.Errorf("error marshalling: %v", err)
	}
	r.dirty = false
	return nil
}

// Remove moves the receiver record to deleted directory. This record won't be loaded anymore with persister init / reload
func (r Record) Remove(path string) error {
	dpath := filepath.Join(path, "deleted")
	file := r.GetFilePath(path)
	dfile := r.GetFilePath(dpath)
	return os.Rename(file, dfile)
}

// GetFilePath returns receiver's full file path name by appening given path and receiver file name
func (r Record) GetFilePath(path string) string {
	return filepath.Join(path, r.GetFileName())
}

// GetFileName returns receiver file name (zero padded id with json extension)
func (r Record) GetFileName() string {
	return fmt.Sprintf("%06d.json", r.id)
}

// NewRecord create a new record with given marshalling function
func NewRecord(marshall func(w io.Writer) error) *Record {
	return &Record{marshall: marshall}
}

// SetIdFromFile sets receiver's id based on given file name (must be zero padded decimal digit)
func (r *Record) SetIdFromFile(file string) error {
	_, err := fmt.Sscanf(filepath.Base(file), "%d.json", &r.id)
	return err
}

// Marshall writes marshalled receiver to given writer
func (r Record) Marshall(w io.Writer) error {
	return r.marshall(w)
}
