package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/scraper"
)

func TestInternalScraper(t *testing.T) {
	ctx := context.TODO()
	page, err := scraper.NewInternalWebScraper().Scrap(ctx, "https://about.readflow.app/")
	assert.NotNil(t, err)
	assert.Equal(t, "unable to extract content from HTML page", err.Error())
	assert.NotNil(t, page)
	assert.Equal(t, "readflow", page.Title)
	assert.Equal(t, "Read your Internet article flow in one place with complete peace of mind and freedom", page.Text)
	assert.Equal(t, "https://about.readflow.app/images/readflow.png", page.Image)
}
