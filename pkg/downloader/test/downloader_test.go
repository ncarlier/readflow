package test

import (
	"context"
	"testing"

	"github.com/ncarlier/readflow/pkg/cache"
	"github.com/ncarlier/readflow/pkg/downloader"
	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	// create download cache
	downloadCache, err := cache.NewDefault("readflow-tests")
	require.Nil(t, err)
	require.NotNil(t, downloadCache)
	err = downloadCache.Clear()
	require.Nil(t, err)
	defer downloadCache.Close()

	url := "https://about.readflow.app/img/logo.svg"
	dl := downloader.NewInternalDownloader(&downloader.InternalDownloaderConfig{
		Cache: downloadCache,
	})

	asset, resp, err := dl.Get(context.Background(), url, nil)
	require.Nil(t, err)
	require.NotNil(t, resp)
	require.Equal(t, 200, resp.StatusCode)
	require.NotNil(t, asset)
	require.Equal(t, url, asset.Name)
	require.Equal(t, "image/svg+xml", asset.ContentType)

	asset, resp, err = dl.Get(context.Background(), url, nil)
	require.Nil(t, err)
	require.Nil(t, resp)
	require.NotNil(t, asset)
}
