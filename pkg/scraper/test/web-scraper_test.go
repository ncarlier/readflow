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

func TestSimpleWebScraping(t *testing.T) {
	ctx := context.TODO()
	page, err := scraper.NewWebScraper(&scraper.WebScraperConfiguration{}).Scrap(ctx, "https://about.readflow.app/")
	assert.Nil(t, err)
	assert.NotNil(t, page)
	assert.Equal(t, "https://about.readflow.app/", page.URL)
	assert.Equal(t, "readflow", page.Title)
	assert.Equal(t, "read your Internet article flow in one place with complete peace of mind and freedom", page.Text)
	assert.Contains(t, page.HTML, "relax.png")
	assert.Equal(t, "https://about.readflow.app/img/readflow.png", page.Image)
	assert.Equal(t, "https://about.readflow.app/favicon.png", page.Favicon)
}

func TestWebScrapingTimeout(t *testing.T) {
	ctx := context.TODO()
	_, err := scraper.NewWebScraper(&scraper.WebScraperConfiguration{
		HttpClient: &http.Client{Timeout: time.Second},
	}).Scrap(ctx, "https://httpstat.us/200?sleep=2000")
	require.NotNil(t, err)
	timeoutErr, ok := err.(net.Error)
	require.True(t, ok)
	require.True(t, timeoutErr.Timeout())
}

func TestMediumWebScraping(t *testing.T) {
	ctx := context.TODO()
	page, err := scraper.NewWebScraper(&scraper.WebScraperConfiguration{}).Scrap(ctx, "https://blog.medium.com/state-of-medium-c54d1706a9b4")
	assert.Nil(t, err)
	assert.NotNil(t, page)
	assert.Equal(t, "State of Medium", page.Title)
	assert.Contains(t, page.HTML, "src=\"https://miro.medium.com/v2/resize:fit:640/1*3hgCec7DtnccJd0vQhP2Kw.jpeg\"")
}
