package test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/scraper"
)

func TestInternalScraper(t *testing.T) {
	ctx := context.TODO()
	page, err := scraper.NewInternalWebScraper(constant.DefaultClient).Scrap(ctx, "https://about.readflow.app/")
	assert.Nil(t, err)
	assert.NotNil(t, page)
	assert.Equal(t, "https://about.readflow.app/", page.URL)
	assert.Equal(t, "readflow", page.Title)
	assert.Equal(t, "Read your Internet article flow in one place with complete peace of mind and freedom", page.Text)
	assert.Contains(t, page.HTML, "relax.png")
	assert.Equal(t, "https://about.readflow.app/images/readflow.png", page.Image)
	assert.Equal(t, "https://about.readflow.app/favicon.png", page.Favicon)
}
