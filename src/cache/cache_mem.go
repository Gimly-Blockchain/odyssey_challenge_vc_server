package cache

type Mem struct {
	store map[string]string
}

// NewMem returns a new instance of Mem with an in-memory map as backing store, typically used for testing.
func NewMem() Cache {
	return Cache(Mem{store: make(map[string]string)})
}

// Set implements the Cache interface for RedisCache.
func (m Mem) Set(did string, resource string) error {
	m.store[did] = resource
	return nil
}

// Get implements the Cache interface for RedisCache.
func (m Mem) Get(did string) (string, error) {
	return m.store[did], nil
}

func (m Mem) Delete(did string) {
	delete(m.store, did)
}
