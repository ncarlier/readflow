package cache

import (
	"fmt"
	"net/url"
	"os"

	boltcache "github.com/ncarlier/readflow/pkg/cache/bolt"
	"github.com/rs/zerolog/log"
)

// DefaultCacheSize is the maximum number of items
const DefaultCacheSize = 100

// Cache interface
type Cache interface {
	Put(key string, data []byte) error
	Get(key string) ([]byte, error)
	Close() error
}

// New create new cache provider regarding the datasource URI
func New(conn string, size int) (Cache, error) {
	u, err := url.ParseRequestURI(conn)
	if err != nil {
		return nil, fmt.Errorf("invalid connection URL: %s", conn)
	}
	provider := u.Scheme
	var cache Cache

	switch provider {
	case "boltdb":
		cache, err = boltcache.New(size, u.Path)
		if err != nil {
			return nil, err
		}
		log.Info().Str("component", "cache").Str("uri", u.Redacted()).Msg("using BoltDB cache")
	default:
		return nil, fmt.Errorf("unsupported cache provider: %s", provider)
	}
	return cache, nil
}

// NewDefault return default cache
func NewDefault() (Cache, error) {
	conn := "boltdb://" + os.TempDir() + string(os.PathSeparator) + "readflow.cache"
	return New(conn, DefaultCacheSize)
}
