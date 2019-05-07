package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/ncarlier/readflow/pkg/tooling"
)

// ArticleForm API structure definition
type ArticleForm struct {
	Title       string     `json:"title,omitempty"`
	Text        *string    `json:"text,omitempty"`
	HTML        *string    `json:"html,omitempty"`
	URL         *string    `json:"url,omitempty"`
	Image       *string    `json:"image,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CategoryID  *uint      `json:"category,omitempty"`
	Tags        *string    `json:"tags,omitempty"`
}

// Articles structure definition
type Articles struct {
	Articles []*Article
	Errors   []error
}

// Article structure definition
type Article struct {
	ID          *uint      `json:"id,omitempty"`
	UserID      uint       `json:"user_id,omitempty"`
	CategoryID  *uint      `json:"category_id,omitempty"`
	Title       string     `json:"title,omitempty"`
	Text        *string    `json:"text,omitempty"`
	HTML        *string    `json:"html,omitempty"`
	URL         *string    `json:"url,omitempty"`
	Image       *string    `json:"image,omitempty"`
	Hash        string     `json:"hash,omitempty"`
	Status      string     `json:"status,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

func (a Article) String() string {
	result, _ := json.Marshal(a)
	return string(result)
}

// ArticlesPageRequest request structure for a paginated list of articles
type ArticlesPageRequest struct {
	Status      *string
	Limit       uint
	AfterCursor *uint
	Category    *uint
	SortOrder   string
}

// ArticlesPageResponse response structure for a paginated list of articles
type ArticlesPageResponse struct {
	TotalCount uint
	EndCursor  uint
	HasNext    bool
	Entries    []*Article
}

// ArticleBuilder is a builder to create an Article
type ArticleBuilder struct {
	article *Article
}

// NewArticleBuilder creates new Article builder instance
func NewArticleBuilder() ArticleBuilder {
	article := &Article{}
	return ArticleBuilder{article}
}

// Build creates the article
func (ab *ArticleBuilder) Build() *Article {
	if ab.article.Status == "" {
		ab.article.Status = "unread"
	}
	payload := ab.article.Title
	if ab.article.URL != nil {
		payload += *ab.article.URL
	}
	if ab.article.HTML != nil {
		payload += *ab.article.HTML
	}
	ab.article.Hash = tooling.Hash(payload)
	// log.Println(ab.article)
	return ab.article
}

// BuildForm creates an article form
func (ab *ArticleBuilder) BuildForm(tags string) *ArticleForm {
	ab.Build()
	return &ArticleForm{
		Title:       ab.article.Title,
		Text:        ab.article.Text,
		HTML:        ab.article.HTML,
		URL:         ab.article.URL,
		Image:       ab.article.Image,
		PublishedAt: ab.article.PublishedAt,
		CategoryID:  ab.article.CategoryID,
		Tags:        &tags,
	}
}

// Random fill article with random data
func (ab *ArticleBuilder) Random() *ArticleBuilder {
	gofakeit.Seed(0)
	ab.article.Title = gofakeit.Sentence(3)
	text := gofakeit.Paragraph(2, 2, 5, ".")
	ab.article.Text = &text
	html := fmt.Sprintf("<p>%s</p>", *ab.article.Text)
	ab.article.HTML = &html
	image := gofakeit.ImageURL(320, 200)
	ab.article.Image = &image
	url := gofakeit.URL()
	ab.article.URL = &url
	publishedAt := gofakeit.Date()
	ab.article.PublishedAt = &publishedAt

	return ab
}

// UserID set article user ID
func (ab *ArticleBuilder) UserID(userID uint) *ArticleBuilder {
	ab.article.UserID = userID
	return ab
}

// CategoryID set article category ID
func (ab *ArticleBuilder) CategoryID(categoryID uint) *ArticleBuilder {
	ab.article.CategoryID = &categoryID
	return ab
}

// Title set article title
func (ab *ArticleBuilder) Title(title string) *ArticleBuilder {
	ab.article.Title = title
	return ab
}

// Text set article text
func (ab *ArticleBuilder) Text(text string) *ArticleBuilder {
	ab.article.Text = &text
	return ab
}

// Form set article content using Form object
func (ab *ArticleBuilder) Form(form *ArticleForm) *ArticleBuilder {
	ab.article.Title = form.Title
	ab.article.Text = form.Text
	ab.article.HTML = form.HTML
	ab.article.URL = form.URL
	ab.article.Image = form.Image
	ab.article.PublishedAt = form.PublishedAt
	ab.article.CategoryID = form.CategoryID
	return ab
}
