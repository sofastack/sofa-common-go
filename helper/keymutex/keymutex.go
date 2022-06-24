package keymutex

import (
    "runtime"
    "sync"

    "github.com/sofastack/sofa-common-go/helper/fnv"
)

// KeyMutex represents a mutex.
type KeyMutex struct {
    mutexes []sync.Mutex
}

// New returns a new KeyMutex.
func New(size int) *KeyMutex {
    if size <= 0 {
        size = runtime.NumCPU()
    }

    return &KeyMutex{
        mutexes: make([]sync.Mutex, size),
    }
}

// Acquires a lock associated with the specified ID.
func (km *KeyMutex) LockKey(id string) {
    km.mutexes[km.hash(id)%uint64(len(km.mutexes))].Lock()
}

// Releases the lock associated with the specified ID.
func (km *KeyMutex) UnlockKey(id string) {
    km.mutexes[km.hash(id)%uint64(len(km.mutexes))].Unlock()
}

func (km *KeyMutex) hash(id string) uint64 {
    return fnv.HashAdd(fnv.HashNew(), id)
}
