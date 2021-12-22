package cache

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	boltcache "github.com/ncarlier/readflow/pkg/cache/bolt"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/rs/zerolog/log"
)

// DefaultCacheSize is the maximum number of items
const DefaultCacheSize = 256

// Cache interface
type Cache interface {
	Put(key string, data *model.FileAsset) error
	Get(key string) (*model.FileAsset, error)
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
	cacheFileName := filepath.ToSlash(os.TempDir() + string(os.PathSeparator) + "readflow.cache")
	os.Remove(cacheFileName)
	conn := "boltdb://" + cacheFileName
	return New(conn, DefaultCacheSize)
}
