package cache

import (
	redisClient "github.com/go-redis/redis/v7"
)

type RedisCache struct {
	rc *redisClient.Client
}

// NewRedis returns a new instance of RedisCache with ru as RedisCache host address.
func NewRedis(ru string) Cache {
	rc := redisClient.NewClient(&redisClient.Options{
		Addr: ru,
	})

	return Cache(RedisCache{rc})
}

// Set implements the Cache interface for RedisCache.
func (r RedisCache) Set(did string, resource string) error {
	return r.rc.Set(did, resource, 0).Err()
}

// Get implements the Cache interface for RedisCache.
func (r RedisCache) Get(did string) (string, error) {
	b, err := r.rc.Get(did).Bytes()
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (r RedisCache) Delete(did string) {
	_ = r.rc.Del(did).Err()
}
