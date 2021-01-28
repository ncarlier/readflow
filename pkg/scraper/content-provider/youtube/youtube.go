package youtube

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"regexp"

	"github.com/ncarlier/readflow/pkg/constant"
	"github.com/ncarlier/readflow/pkg/scraper"
)

var ytWRe = regexp.MustCompile(`https://.+\.youtube.com/watch.+`)
var ytLRe = regexp.MustCompile(`https://.+\.youtube.com/v/.+`)
var ytSRe = regexp.MustCompile(`https://youtu.be/.+`)

const ytOEmbedBaseURL = "https://www.youtube.com/oembed?maxheight=600&maxwidth=800&format=json&url="

type oEmbedResponse struct {
	Title        string `json:"title,omitempty"`
	HTML         string `json:"html,omitempty"`
	ThumbnailURL string `json:"thumbnail_url,omitempty"`
	AuthorName   string `json:"author_name,omitempty"`
}

type youtubeContentProvider struct {
	httpClient *http.Client
}

func newYoutubeContentProvider() *youtubeContentProvider {
	return &youtubeContentProvider{
		httpClient: &http.Client{
			Timeout: constant.DefaultTimeout,
		},
	}
}

func (ycp youtubeContentProvider) Get(ctx context.Context, rawurl string) (*scraper.WebPage, error) {
	oembedURL := ytOEmbedBaseURL + rawurl

	req, err := http.NewRequest("GET", oembedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", constant.UserAgent)
	req = req.WithContext(ctx)
	res, err := ycp.httpClient.Do(req)
	if err != nil || res.StatusCode >= 300 {
		if err == nil {
			err = fmt.Errorf("bad status code: %d", res.StatusCode)
		}
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	oembed := oEmbedResponse{}
	err = json.Unmarshal(body, &oembed)
	if err != nil {
		return nil, err
	}

	return &scraper.WebPage{
		Title:    oembed.Title,
		HTML:     oembed.HTML,
		Image:    oembed.ThumbnailURL,
		URL:      rawurl,
		Text:     "Youtube video from " + oembed.AuthorName,
		SiteName: "Youtube",
	}, nil
}

func (ycp youtubeContentProvider) Match(url string) bool {
	return ytWRe.MatchString(url) || ytLRe.MatchString(url) || ytSRe.MatchString(url)
}

func getYoutubeVideoID(rawurl string) string {
	u, err := url.Parse(rawurl)
	if err != nil {
		return ""
	}
	if ytWRe.MatchString(rawurl) {
		q := u.Query()
		return q.Get("v")
	}
	return path.Base(u.Path)
}

func init() {
	scraper.ContentProviders["youtube"] = newYoutubeContentProvider()
}
