package env

import (
	"errors"

	"github.com/commercionetwork/dsb/src/cache"
)

const (
	defaultStoragePath  = "./dsb-storage"
	defaultLogPath      = "./dsb.log"
	defaultListenAddr   = "localhost:9999"
	defaultRedisAddr    = "127.0.0.1:6379"
	defaultCommercioLCD = "http://127.0.0.1:1317"
)

type CacheType int

const (
	CacheTypeRedis CacheType = iota
	CacheTypeMemory
)

// Variables holds information about the Variables tumbler is running in, including those needed to
// communicate with the commercio.network LCD server, the tumbler private key.
type Variables struct {
	StoragePath   string
	CommercioLCD  string
	LogPath       string
	RedisAddr     string
	RedisPass     string
	ListenAddr    string
	JWTSecret     string
	CacheType     CacheType
	Debug         bool
	cacheInstance cache.Cache
}

func (v *Variables) CacheInstance() cache.Cache {
	return v.cacheInstance
}

func (v *Variables) Validate() error {
	if v.StoragePath == "" {
		v.StoragePath = defaultStoragePath
	}

	if v.LogPath == "" {
		v.LogPath = defaultLogPath
	}

	if v.RedisAddr == "" {
		v.RedisAddr = defaultRedisAddr
	}

	if v.ListenAddr == "" {
		v.ListenAddr = defaultListenAddr
	}

	if v.CommercioLCD == "" {
		v.CommercioLCD = defaultCommercioLCD
	}

	if v.JWTSecret == "" {
		return errors.New("must provide jwt secret")
	}

	switch v.CacheType {
	case CacheTypeRedis:
		v.cacheInstance = cache.NewRedis(v.RedisAddr)
	case CacheTypeMemory:
		v.cacheInstance = cache.NewMem()
	default:
		v.CacheType = CacheTypeRedis
		v.cacheInstance = cache.NewRedis(v.RedisAddr)
	}

	return nil
}

var evMapping = map[string]string{
	"DSB_STORAGE_PATH":  "StoragePath",
	"DSB_LOG_PATH":      "LogPath",
	"DSB_DEBUG":         "Debug",
	"DSB_REDIS_ADDR":    "RedisAddr",
	"DSB_REDIS_PASS":    "RedisPass",
	"DSB_CACHE_TYPE":    "CacheType",
	"DSB_JWT_SECRET":    "JWTSecret",
	"DSB_COMMERCIO_LCD": "CommercioLCD",
}
