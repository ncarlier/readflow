package test

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ncarlier/readflow/pkg/scraper"
)

func TestInternalScraper(t *testing.T) {
	ctx := context.TODO()
	page, err := scraper.NewInternalWebScraper(&scraper.WebScraperConfiguration{}).Scrap(ctx, "https://about.readflow.app/")
	assert.Nil(t, err)
	assert.NotNil(t, page)
	assert.Equal(t, "https://about.readflow.app/", page.URL)
	assert.Equal(t, "readflow", page.Title)
	assert.Equal(t, "read your Internet article flow in one place with complete peace of mind and freedom", page.Text)
	assert.Contains(t, page.HTML, "relax.png")
	assert.Equal(t, "https://about.readflow.app/images/readflow.png", page.Image)
	assert.Equal(t, "https://about.readflow.app/favicon.png", page.Favicon)
}

func TestInternalScraperTimeout(t *testing.T) {
	ctx := context.TODO()
	_, err := scraper.NewInternalWebScraper(&scraper.WebScraperConfiguration{
		HttpClient: &http.Client{Timeout: time.Second},
	}).Scrap(ctx, "https://httpstat.us/200?sleep=2000")
	require.NotNil(t, err)
	timeoutErr, ok := err.(net.Error)
	require.True(t, ok)
	require.True(t, timeoutErr.Timeout())
}
