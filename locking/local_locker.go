package locking

import (
	"context"
	"sync"

	"github.com/cloudbase/garm/database/common"
)

func NewLocalLocker(_ context.Context, _ common.Store) (Locker, error) {
	return &keyMutex{}, nil
}

type keyMutex struct {
	muxes sync.Map
}

var _ Locker = &keyMutex{}

func (k *keyMutex) TryLock(key string) bool {
	mux, _ := k.muxes.LoadOrStore(key, &sync.Mutex{})
	keyMux := mux.(*sync.Mutex)
	return keyMux.TryLock()
}

func (k *keyMutex) Unlock(key string, remove bool) {
	mux, ok := k.muxes.Load(key)
	if !ok {
		return
	}
	keyMux := mux.(*sync.Mutex)
	if remove {
		k.Delete(key)
	}
	keyMux.Unlock()
}

func (k *keyMutex) Delete(key string) {
	k.muxes.Delete(key)
}
