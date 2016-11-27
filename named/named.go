package named

import (
	"sync"
)

type Registry struct {
	sync.RWMutex
	entries map[string]interface{}
}

func NewRegistry() *Registry {
	return &Registry{
		entries: make(map[string]interface{}),
	}
}

func (r *Registry) Get(name string, factory func() interface{}) interface{} {
	if len(name) == 0 {
		return nil
	}
	r.RLock()
	entry, ok := r.entries[name]
	r.RUnlock()
	if ok {
		return entry
	}
	r.Lock()
	defer r.Unlock()
	if entry, ok := r.entries[name]; ok {
		return entry
	}
	if factory == nil {
		return nil
	}
	entry = factory()
	if entry == nil {
		return nil
	}
	r.entries[name] = entry
	return entry
}

func (r *Registry) Set(name string, entry interface{}) {
	r.Lock()
	defer r.Unlock()
	r.entries[name] = entry
}

func (r *Registry) Delete(name string) {
	r.Lock()
	defer r.Unlock()
	delete(r.entries, name)
}

func (r *Registry) DeleteAll() {
	r.Lock()
	defer r.Unlock()
	for k := range r.entries {
		delete(r.entries, k)
	}
}
