package readability

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	nurl "net/url"
	"strings"
	"time"

	read "github.com/go-shiori/go-readability"
	"github.com/ncarlier/readflow/pkg/model"
	"github.com/ncarlier/readflow/pkg/tooling"
	"golang.org/x/net/html/charset"
)

func getContentType(ctx context.Context, url string) (string, error) {
	ctx, _ = context.WithTimeout(ctx, 10*time.Second)
	req, _ := http.NewRequest("HEAD", url, nil)
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	return res.Header.Get("Content-type"), nil
}

func get(ctx context.Context, url string) (*http.Response, error) {
	ctx, _ = context.WithTimeout(ctx, 10*time.Second)
	req, _ := http.NewRequest("GET", url, nil)
	req = req.WithContext(ctx)
	return http.DefaultClient.Do(req)
}

// FetchArticle fetch article from an URL
func FetchArticle(ctx context.Context, url string) (*model.Article, error) {
	// Validate URL
	_, err := nurl.ParseRequestURI(url)
	if err != nil {
		return nil, fmt.Errorf("invalid URL: %v", err)
	}

	// Get URL content type
	contentType, err := getContentType(ctx, url)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(contentType, "text/html") {
		return nil, fmt.Errorf("invalid content-type: %s", contentType)
	}

	// Get URL content
	res, err := get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := charset.NewReader(res.Body, contentType)
	if err != nil {
		return nil, err
	}

	// Extract metas
	metas, err := ExtractMetas(body)
	if err != nil {
		return nil, err
	}

	// Create article with Open Graph atributes
	result := &model.Article{
		Text:  metas.GetContent("og:description", "twitter:description", "description"),
		Image: metas.GetContent("og:image", "twitter:image"),
	}
	title := metas.GetContent("og:title")
	if title != nil {
		result.Title = *title
	}

	var buffer bytes.Buffer
	tee := io.TeeReader(body, &buffer)

	// Test if the HTML page is readable by Shiori readability
	if !read.IsReadable(tee) {
		return result, fmt.Errorf("unable to extract content from HTML page")
	}

	// Extract content from the HTML page
	article, err := read.FromReader(&buffer, url)
	if err != nil {
		return result, err
	}

	// Complete result with extracted properties
	result.HTML = &article.Content
	if result.Title == "" {
		result.Title = article.Title
	}
	if result.Text == nil {
		// FIXME: readability excerpt don't well support UTF8
		text := tooling.ToUTF8(article.Excerpt)
		result.Text = &text
	}
	if result.Image == nil {
		result.Image = &article.Image
	}

	// TODO: add other properties to the result
	// article.Favicon
	// article.Length
	// article.SiteName

	return result, nil
}
