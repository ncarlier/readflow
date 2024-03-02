package test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ncarlier/readflow/internal/model"
	"github.com/ncarlier/readflow/internal/service"
)

func TestDownloadArticle(t *testing.T) {
	teardownTestCase := SetupTestCase(t)
	defer teardownTestCase(t)

	// Create new article
	url := "https://github.com/ncarlier/readflow"
	title := "article download test"
	req := model.ArticleCreateForm{
		URL:   &url,
		Title: title,
	}
	opts := service.ArticleCreationOptions{}
	art, err := service.Lookup().CreateArticle(testContext, req, opts)
	require.Nil(t, err)
	require.Equal(t, title, art.Title)
	require.Equal(t, url, *art.URL)
	require.NotNil(t, art.Image)
	require.True(t, strings.HasPrefix(*art.Image, "https://repository-images.githubusercontent.com"), "unexpected image URL")
	// Download article
	asset, err := service.Lookup().DownloadArticle(testContext, art.ID, "html-single")
	require.Nil(t, err)
	require.Equal(t, "text/html; charset=utf-8", asset.ContentType)
	require.NotNil(t, asset.Data)
	// htmlContent := string(content)
	// require.True(t, strings.Contains(htmlContent, "<img src=\"/ncarlier/readflow/raw/master/readflow.svg\" alt=\"Logo\" style=\"max-width:100%;\">"))
}
