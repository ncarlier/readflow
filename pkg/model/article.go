package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/ncarlier/readflow/pkg/helper"
	"github.com/ncarlier/readflow/pkg/html"
)

// ArticleCreateForm structure definition
type ArticleCreateForm struct {
	Title       string     `json:"title,omitempty"`
	Text        *string    `json:"text,omitempty"`
	HTML        *string    `json:"html,omitempty"`
	URL         *string    `json:"url,omitempty"`
	Image       *string    `json:"image,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CategoryID  *uint      `json:"category,omitempty"`
	Origin      *string    `json:"origin,omitempty"`
	Tags        *string    `json:"tags,omitempty"`
}

// IsComplete test if the form is complete
func (form ArticleCreateForm) IsComplete() bool {
	return !helper.OneIsEmpty(form.Image, form.Text, form.HTML)
}

// Hash return form hash
func (form ArticleCreateForm) Hash() string {
	key := form.Title
	if form.URL != nil {
		key += *form.URL
	}
	if form.HTML != nil {
		key += *form.HTML
	}
	return helper.Hash(key)
}

// Payload return form payload (content without HTML tags)
func (form ArticleCreateForm) Payload() string {
	payload := ""
	if form.HTML != nil {
		payload = *form.HTML
		if form.Text != nil && len(*form.Text) > len(payload) {
			payload = *form.Text
		}
	} else if form.Text != nil {
		payload = *form.Text
	}
	if text, err := html.HTML2Text(payload); err != nil {
		payload = text
	}
	return payload
}

// TruncatedTitle return truncated title
func (form ArticleCreateForm) TruncatedTitle() string {
	return helper.Truncate(form.Title, 29)
}

// ArticleUpdateForm structure definition
type ArticleUpdateForm struct {
	ID         uint
	Title      *string
	Text       *string
	CategoryID *uint
	Status     *string
	Stars      *int
}

// CreatedArticlesResponse structure definition
type CreatedArticlesResponse struct {
	Articles []*Article
	Errors   []error
}

// Article structure definition
type Article struct {
	ID          uint       `json:"id,omitempty"`
	UserID      uint       `json:"user_id,omitempty"`
	CategoryID  *uint      `json:"category_id,omitempty"`
	Title       string     `json:"title,omitempty"`
	Text        *string    `json:"text,omitempty"`
	HTML        *string    `json:"html,omitempty"`
	URL         *string    `json:"url,omitempty"`
	Image       *string    `json:"image,omitempty"`
	Hash        string     `json:"hash,omitempty"`
	Status      string     `json:"status,omitempty"`
	Stars       uint       `json:"stars,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// String convert article to JSON string
func (a Article) String() string {
	result, _ := json.Marshal(a)
	return string(result)
}

// ToMap convert article to map
func (a Article) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":           a.ID,
		"title":        a.Title,
		"text":         helper.PtrValueOr(a.Text, ""),
		"html":         helper.PtrValueOr(a.HTML, ""),
		"url":          helper.PtrValueOr(a.URL, ""),
		"image":        helper.PtrValueOr(a.Image, ""),
		"published_at": helper.PtrValueOr(a.PublishedAt, time.Now()),
	}
}

// ArticlesPageRequest request structure for a paginated list of articles
type ArticlesPageRequest struct {
	Status      *string
	Starred     *bool
	Category    *uint
	Query       *string
	AfterCursor *uint
	SortOrder   *string
	SortBy      *string
	Limit       *int
}

// ArticlesPageResponse response structure for a paginated list of articles
type ArticlesPageResponse struct {
	TotalCount uint
	EndCursor  uint
	HasNext    bool
	Entries    []*Article
}

// ArticleCreateFormBuilder is a builder to create an Article create form
type ArticleCreateFormBuilder struct {
	form *ArticleCreateForm
}

// NewArticleCreateFormBuilder creates new Article create form builder instance
func NewArticleCreateFormBuilder() ArticleCreateFormBuilder {
	form := &ArticleCreateForm{}
	return ArticleCreateFormBuilder{form}
}

// Build creates the article create form
func (b *ArticleCreateFormBuilder) Build() *ArticleCreateForm {
	return b.form
}

// Random fill article create form with random data
func (b *ArticleCreateFormBuilder) Random() *ArticleCreateFormBuilder {
	b.form = &ArticleCreateForm{}
	gofakeit.Seed(0)
	b.form.Title = gofakeit.Sentence(3)
	text := gofakeit.Paragraph(2, 2, 5, ".")
	b.form.Text = &text
	_html := fmt.Sprintf("<p>%s</p>", *b.form.Text)
	b.form.HTML = &_html
	image := gofakeit.ImageURL(320, 200)
	b.form.Image = &image
	url := gofakeit.URL()
	b.form.URL = &url
	publishedAt := gofakeit.Date()
	b.form.PublishedAt = &publishedAt

	return b
}

// FromArticle fill article create form internal article
func (b *ArticleCreateFormBuilder) FromArticle(article Article) *ArticleCreateFormBuilder {
	b.form.HTML = article.HTML
	b.form.Image = article.Image
	b.form.PublishedAt = article.PublishedAt
	b.form.Text = article.Text
	b.form.Title = article.Title
	b.form.URL = article.URL
	return b
}

// CategoryID set article category ID
func (b *ArticleCreateFormBuilder) CategoryID(categoryID uint) *ArticleCreateFormBuilder {
	b.form.CategoryID = &categoryID
	return b
}

// Title set article title
func (b *ArticleCreateFormBuilder) Title(title string) *ArticleCreateFormBuilder {
	b.form.Title = title
	return b
}

// Text set article text
func (b *ArticleCreateFormBuilder) Text(text string) *ArticleCreateFormBuilder {
	b.form.Text = &text
	return b
}

// HTML set article HTML
func (b *ArticleCreateFormBuilder) HTML(_html string) *ArticleCreateFormBuilder {
	b.form.HTML = &_html
	return b
}

// Origin set article origin
func (b *ArticleCreateFormBuilder) Origin(origin string) *ArticleCreateFormBuilder {
	b.form.Origin = &origin
	return b
}

// Tags set article tags
func (b *ArticleCreateFormBuilder) Tags(tags string) *ArticleCreateFormBuilder {
	b.form.Tags = &tags
	return b
}
