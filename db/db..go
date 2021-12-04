package db

import "sync"

type HostDb interface {
	Save(host string) error
	Delete(host string) error
	Get(host string) (bool, error)
}

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

func NewInMemory() *inMemoryDb {
	return &inMemoryDb{hosts: make(map[string]bool)}
}
