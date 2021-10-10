package oembed

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/scraper"
)

func TestOEmbedContentProvider(t *testing.T) {
	ctx := context.TODO()

	rawurl := "https://www.youtube.com/watch?v=ee-LhNZPZ1U"

	provider := scraper.GetContentProvider(rawurl)
	assert.NotNil(t, provider, "content provider not found")
	assert.True(t, provider.Match(rawurl))

	page, err := provider.Get(ctx, rawurl)
	assert.NoError(t, err)
	assert.NotNil(t, page)
	assert.Equal(t, "YouTube", page.SiteName)
	assert.Equal(t, "Les Circuits Intégrés", page.Title)
	assert.Equal(t, "YouTube video from Deus Ex Silicium", page.Text)
}
