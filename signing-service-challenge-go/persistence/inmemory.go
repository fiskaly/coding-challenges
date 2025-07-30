package persistence

import (
	"fmt"
	"sync"
)

func MakeInMemoryDB() inMemoryDB {
	return inMemoryDB{
		mu: sync.RWMutex{},
		db: dbMap{},
	}
}

type dbMap map[ID]Device

type inMemoryDB struct {
	mu sync.RWMutex
	db dbMap
}

func (mdb *inMemoryDB) AddDevice(d Device) error {
	mdb.mu.Lock()
	defer mdb.mu.Unlock()
	if _, exists := mdb.db[d.ID]; exists {
		return ErrAlreadyExists
	}
	mdb.db[d.ID] = d
	return nil
}

func (mdb *inMemoryDB) ListDevices() ([]ID, error) {
	mdb.mu.RLock()
	defer mdb.mu.RUnlock()
	ids := make([]ID, 0, len(mdb.db))
	for id := range mdb.db {
		ids = append(ids, id)
	}
	return ids, nil
}

func (mdb *inMemoryDB) GetDevice(id ID) (*Device, error) {
	mdb.mu.RLock()
	defer mdb.mu.RUnlock()
	d, exists := mdb.db[id]
	if !exists {
		return nil, fmt.Errorf("device not found")
	}
	return &d, nil
}

func (mdb *inMemoryDB) CompareAndSwapSignature(id ID, o Signature, n Signature) error {
	mdb.mu.Lock()
	defer mdb.mu.Unlock()
	d, exists := mdb.db[id]
	if !exists {
		return fmt.Errorf("device not found")
	}
	if d.LastSignature != o {
		return ErrRace
	}
	d.LastSignature = n
	d.SignatureCount++
	mdb.db[id] = d
	return nil
}
