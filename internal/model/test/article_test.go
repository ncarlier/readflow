package test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ncarlier/readflow/internal/model"
)

func TestArticleCreateFormTags(t *testing.T) {
	title := "This is a #title with #hastags"
	builder := model.NewArticleCreateFormBuilder()
	form := builder.Random().Title(title).Build()

	require.Equal(t, form.Title, title)
	hashtags := form.Hashtags()
	require.Len(t, hashtags, 2)
	require.Contains(t, hashtags, "#title")
}
