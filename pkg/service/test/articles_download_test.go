package test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/service"
)

func TestDownloadArticle(t *testing.T) {
	teardownTestCase := setupTestCase(t)
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
	assert.Nil(t, err)
	assert.Equal(t, title, art.Title)
	assert.Equal(t, url, *art.URL)
	assert.NotNil(t, art.Image)
	assert.True(t, strings.HasPrefix(*art.Image, "https://repository-images.githubusercontent.com"), "unexpected image URL")
	// Download article
	asset, err := service.Lookup().DownloadArticle(testContext, art.ID, "html-single")
	assert.Nil(t, err)
	assert.Equal(t, "text/html; charset=utf-8", asset.ContentType)
	assert.NotNil(t, asset.Data)
	// htmlContent := string(content)
	// assert.True(t, strings.Contains(htmlContent, "<img src=\"/ncarlier/readflow/raw/master/readflow.svg\" alt=\"Logo\" style=\"max-width:100%;\">"))
}
