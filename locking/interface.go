package locking

// TODO(gabriel-samfira): needs owner attribute.
type Locker interface {
	TryLock(key string) bool
	Unlock(key string, remove bool)
	Delete(key string)
}
