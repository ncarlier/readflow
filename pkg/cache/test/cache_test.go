package test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ncarlier/readflow/pkg/cache"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	c, err := cache.NewDefault("readflow-tests")
	require.Nil(t, err)
	require.NotNil(t, c)

	defer c.Close()
	value, err := c.Get("foo")
	require.Nil(t, err)
	require.Empty(t, value)

	value = []byte("value")
	err = c.Put("foo", value)
	require.Nil(t, err)

	value, err = c.Get("foo")
	require.Nil(t, err)
	require.NotEmpty(t, value)
}

func TestCacheEviction(t *testing.T) {
	cacheFileName := filepath.ToSlash(filepath.Join(os.TempDir(), "readflow-tests.cache"))
	os.Remove(cacheFileName)
	conn := "boltdb://" + cacheFileName + "?maxEntries=5"
	c, err := cache.New(conn)
	require.Nil(t, err)
	require.NotNil(t, c)

	defer c.Close()

	value := []byte("bar")
	for i := 0; i < 6; i++ {
		err = c.Put(fmt.Sprintf("foo-%d", i), value)
		require.Nil(t, err)
	}

	value, err = c.Get("foo-0")
	require.Nil(t, err)
	require.Empty(t, value)

	value, err = c.Get("foo-1")
	require.Nil(t, err)
	require.NotEmpty(t, value)
}
