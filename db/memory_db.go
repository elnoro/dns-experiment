package db

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
)

func (i *inMemoryDb) Save(host string) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	i.hosts[host] = true

	return nil
}

func (i *inMemoryDb) Delete(host string) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	delete(i.hosts, host)

	return nil
}

func (i *inMemoryDb) Get(host string) (bool, error) {
	i.mutex.RLock()
	defer i.mutex.RUnlock()

	_, ok := i.hosts[host]

	return ok, nil
}

type inMemoryDb struct {
	hosts map[string]bool
	mutex sync.RWMutex
}

func newEmptyInMemory() *inMemoryDb {
	return &inMemoryDb{hosts: make(map[string]bool)}
}

func NewInMemoryFromFile(path string) (*inMemoryDb, error) {
	db := newEmptyInMemory()
	if path == "" {
		return db, nil
	}

	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %s, %w", path, err)
	}
	hosts := strings.Split(string(f), "\n")

	for _, host := range hosts {
		err = db.Save(host)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
