package cache

// Cache represents an object capable of setting and getting data from a backing storage (RedisCache, a map...).
type Cache interface {
	Set(did string, resource string) error
	Get(did string) (string, error)
	Delete(did string)
}
