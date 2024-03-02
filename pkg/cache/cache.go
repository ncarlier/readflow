package cache

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	boltcache "github.com/ncarlier/readflow/pkg/cache/bolt"
	"github.com/ncarlier/readflow/pkg/logger"
)

// Cache interface
type Cache interface {
	Put(key string, value []byte) error
	Get(key string) ([]byte, error)
	Clear() error
	Close() error
}

// New create new cache provider regarding the datasource URI
func New(conn string) (Cache, error) {
	u, err := url.ParseRequestURI(conn)
	if err != nil {
		return nil, fmt.Errorf("invalid connection URL: %s", conn)
	}
	provider := u.Scheme
	var cache Cache

	switch provider {
	case "boltdb":
		cache, err = boltcache.New(filepath.Clean(u.Path), u.Query())
		if err != nil {
			return nil, err
		}
		logger.Info().Str("component", "cache").Str("uri", u.Redacted()).Msg("using BoltDB cache")
	default:
		return nil, fmt.Errorf("unsupported cache provider: %s", provider)
	}
	return cache, nil
}

// NewDefault return default cache
func NewDefault(name string) (Cache, error) {
	cacheFileName := filepath.Join(os.TempDir(), name+".cache")
	os.Remove(cacheFileName)
	conn := "boltdb://" + cacheFileName
	return New(conn)
}
