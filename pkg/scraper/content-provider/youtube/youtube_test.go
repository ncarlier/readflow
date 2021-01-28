package youtube

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/scraper"
)

func TestYoutubeContentProvider(t *testing.T) {
	ctx := context.TODO()

	rawurl := "https://www.youtube.com/watch?v=ee-LhNZPZ1U"

	provider := scraper.GetContentProvider(rawurl)
	assert.NotNil(t, provider, "content provider not found")
	assert.True(t, provider.Match(rawurl))

	page, err := provider.Get(ctx, rawurl)
	assert.NoError(t, err)
	assert.NotNil(t, page)
	assert.Equal(t, "Youtube", page.SiteName)
	assert.Equal(t, "Deus Ex Silicium : Les Circuits Intégrés", page.Title)
}
