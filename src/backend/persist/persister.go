package persist

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Recorder interface {
	GetId() int
	SetId(id int)
	Dirty()
	Persist(path string) error
	Remove(path string) error
	GetFileName() string
	Marshall(writer io.Writer) error
}

const (
	DefaultPersistDelay = 1 * time.Second
	ParallelPersister   = 10
)

type Persister struct {
	name      string
	directory string
	delay     time.Duration
	records   map[int]Recorder
	nextId    int

	mut          sync.RWMutex
	persistDone  *sync.Cond
	dirtyIds     []int
	persistTimer *time.Timer
}

// NewPersister creates a new persister with given name and storing its record in given dir directory
func NewPersister(name, dir string) *Persister {
	p := &Persister{
		name:      name,
		directory: dir,
		delay:     DefaultPersistDelay,
	}
	p.persistDone = sync.NewCond(&p.mut)
	p.Reinit()
	return p
}

func (p *Persister) NbRecords() int {
	return len(p.records)
}

// Reinit waits persister mechanism to finish (if any) and reset the persister (empty record and id counter reset to 0)
func (p *Persister) Reinit() {
	p.WaitPersistDone()
	p.mut.Lock()
	defer p.mut.Unlock()
	p.records = make(map[int]Recorder)
	p.nextId = 0
}

// SetPersistDelay sets the Pesistance Delay of the Persister
//
// if persistDelay is set to 0, dirty records will be synchronously persisted (writen to disk)
func (p *Persister) SetPersistDelay(persistDelay time.Duration) {
	p.delay = persistDelay
}

// NoDelay suppresses receiver persist delay : any record marked as dirt will be persisted synchronously
func (p *Persister) NoDelay() {
	p.SetPersistDelay(0)
}

// CheckDirectory checks if Persister directory exists and create deleted dir if not exists (returns nil error if ok)
func (p *Persister) CheckDirectory() error {
	fi, err := os.Stat(p.directory)
	if err != nil {
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("not a proper directory: %s\n", p.directory)
	}
	dpath := filepath.Join(p.directory, "deleted")
	if _, err := os.Stat(dpath); os.IsNotExist(err) {
		return os.Mkdir(dpath, os.ModePerm)
	}
	return nil
}

// LoadDirectory loads all persisted Records
func (p *Persister) LoadDirectory(recordType string, recordFactory func(string) (Recorder, error)) error {
	p.WaitPersistDone() // first wait for all ongoing operations to end
	p.mut.Lock()
	defer p.mut.Unlock()

	p.records = make(map[int]Recorder)
	p.nextId = 0

	files, err := p.GetFilesList("deleted")
	if err != nil {
		return fmt.Errorf("could not get files from %s persister: %v", p.name, err)
	}

	for _, file := range files {
		ar, err := recordFactory(file)
		if err != nil {
			return fmt.Errorf("could not instantiate %s from '%s': %v", recordType, filepath.Base(file), err)
		}
		err = p.load(ar)
		if err != nil {
			return fmt.Errorf("error while loading %s: %s", file, err.Error())
		}
	}
	return nil
}

// HasId returns true if persister contains a record with given id, false otherwise
func (p *Persister) HasId(id int) bool {
	if _, ok := p.records[id]; ok {
		return true
	}
	return false
}

// GetFilesList returns all the record files contained in persister directory (User class is responsible to Load the record)
func (p *Persister) GetFilesList(skipdir string) (list []string, err error) {
	err = filepath.Walk(p.directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == skipdir {
				return filepath.SkipDir
			}
			return nil
		}
		list = append(list, path)
		return nil
	})
	if err != nil {
		return
	}
	return
}

// Load adds the given Record to the Persister
func (p *Persister) GetRecords() []Recorder {
	p.mut.RLock()
	defer p.mut.RUnlock()
	res := make([]Recorder, len(p.records))
	i := 0
	for _, recorder := range p.records {
		res[i] = recorder
		i++
	}
	return res
}

// GetById returns a recorder with given Id (or nil if Id not found)
func (p *Persister) GetById(id int) Recorder {
	p.mut.RLock()
	defer p.mut.RUnlock()
	return p.records[id]
}

// Add adds the given Record to the Persister, assigns it a new id, triggers Persit mechanism and returns its (new) id
func (p *Persister) Add(r Recorder) int {
	p.mut.Lock()
	defer p.mut.Unlock()
	defer func() { p.nextId++ }()

	r.SetId(p.nextId)
	p.records[p.nextId] = r
	p.markDirty(r)

	return p.nextId
}

// Update the given Record to the Persister and triggers Persit mechanism
func (p *Persister) Update(r Recorder) error {
	rId := r.GetId()
	if !p.HasId(rId) {
		return fmt.Errorf("record with id %d not found", rId)
	}
	p.mut.Lock()
	defer p.mut.Unlock()
	p.records[rId] = r
	p.markDirty(r)

	return nil
}

// Load adds the given Record to the Persister
func (p *Persister) Load(r Recorder) error {
	p.mut.Lock()
	defer p.mut.Unlock()

	return p.load(r)
}

// load adds the given Recorder to the Persister
func (p *Persister) load(r Recorder) error {
	rId := r.GetId()
	if p.HasId(rId) {
		return fmt.Errorf("persister already contains a record with Id %d", rId)
	}
	p.records[rId] = r
	if p.nextId <= rId {
		p.nextId = rId + 1
	}
	return nil
}

// markDirty marks the given recorder as dirty and triggers the persistence mechanism
func (p *Persister) MarkDirty(r Recorder) {
	p.mut.Lock()
	defer p.mut.Unlock()
	p.markDirty(r)
}

func (p *Persister) markDirty(r Recorder) {
	if _, ok := p.records[r.GetId()]; !ok {
		return
	}
	r.Dirty()
	p.dirtyIds = append(p.dirtyIds, r.GetId())
	p.triggerPersist()
}

// Remove removes the given recorder from the persister (pertaining persisted file is deleted)
func (p *Persister) Remove(r Recorder) error {
	id := r.GetId()
	if _, ok := p.records[id]; !ok {
		return fmt.Errorf("persister does not contains given record")
	}
	p.mut.Lock()
	defer p.mut.Unlock()
	go func(dr Recorder) {
		err := dr.Remove(p.directory)
		if err != nil {
			log.Printf("error removing record GetId %d: %v\n", id, err)
		}
	}(r)
	delete(p.records, id)
	return nil
}

// PersistAll immediatly persist all contained recorder(persistance delay is ignored)
func (p *Persister) PersistAll() {
	p.mut.Lock()
	defer p.mut.Unlock()
	// desactivate persistMechanism if activated
	if p.persistTimer != nil {
		p.persistTimer.Stop()
		p.persistTimer = nil
		p.dirtyIds = []int{}
	}

	token := make(chan struct{}, ParallelPersister)
	defer close(token)
	for _, r := range p.records {
		token <- struct{}{}
		go func(pr Recorder) {
			err := r.Persist(p.directory)
			if err != nil {
				log.Printf("error persisting record ID %d: %v\n", r.GetId(), err)
			}
			_ = <-token
		}(r)
	}

	for i := 0; i < ParallelPersister; i++ {
		token <- struct{}{}
	}
}

func (p *Persister) triggerPersist() {
	if p.delay == 0 {
		//p.mut.Lock()
		//defer p.mut.Unlock()
		if p.persistTimer != nil {
			p.persistTimer.Stop()
			p.persistTimer = nil
		}
		p.persist()
		return
	}
	if p.persistTimer != nil {
		return
	}
	p.persistTimer = time.AfterFunc(p.delay, func() {
		p.mut.Lock()
		defer p.mut.Unlock()
		p.persistTimer = nil
		p.persist()
	})
}

func (p *Persister) persist() {
	token := make(chan struct{}, ParallelPersister)
	for _, id := range p.dirtyIds {
		r, found := p.records[id]
		if !found { // can happen if record was remove before persistence was triggered
			continue
		}
		token <- struct{}{}
		go func(pr Recorder) {
			err := pr.Persist(p.directory)
			if err != nil {
				log.Printf("error persisting record ID %d: %v\n", pr.GetId(), err)
			}
			_ = <-token
		}(r)
	}
	for i := 0; i < ParallelPersister; i++ {
		token <- struct{}{}
	}
	p.dirtyIds = []int{}
	p.persistDone.Broadcast()
}

// WaitPersistDone waits for current persisting mechanism to end and return (return instantly if no persist in progress)
func (p *Persister) WaitPersistDone() {
	if p.persistTimer == nil && len(p.dirtyIds) == 0 {
		return
	}
	p.persistDone.L.Lock()
	p.persistDone.Wait()
	p.persistDone.L.Unlock()
}

// GetName returns persister's name
func (p *Persister) GetName() string {
	return p.name
}

func (p *Persister) GetArchivableRecords() []ArchivableRecord {
	p.mut.RLock()
	defer p.mut.RUnlock()
	res := make([]ArchivableRecord, p.NbRecords())
	i := 0
	for _, record := range p.records {
		res[i] = record
		i++
	}
	return res
}
