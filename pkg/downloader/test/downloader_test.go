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
	defer downloadCache.Close()

	url := "https://about.readflow.app/img/logo.svg"
	dl := downloader.NewDefaultDownloader(downloadCache)
	asset, err := dl.Download(context.Background(), url)
	require.Nil(t, err)
	require.NotNil(t, asset)
	require.Equal(t, url, asset.Name)
	require.Equal(t, "image/svg+xml", asset.ContentType)
}
