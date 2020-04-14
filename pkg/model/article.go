package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/brianvoe/gofakeit"
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
	Tags        *string    `json:"tags,omitempty"`
}

// ArticleUpdateForm structure definition
type ArticleUpdateForm struct {
	ID     uint
	Status *string
}

// Articles structure definition
type Articles struct {
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

// ArticleStatusResponse is the response structure of an article status modification
type ArticleStatusResponse struct {
	Article *Article  `json:"article,omitempty"`
	All     *Category `json:"_all,omitempty"`
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
	html := fmt.Sprintf("<p>%s</p>", *b.form.Text)
	b.form.HTML = &html
	image := gofakeit.ImageURL(320, 200)
	b.form.Image = &image
	url := gofakeit.URL()
	b.form.URL = &url
	publishedAt := gofakeit.Date()
	b.form.PublishedAt = &publishedAt

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

// Tags set article tags
func (b *ArticleCreateFormBuilder) Tags(tags string) *ArticleCreateFormBuilder {
	b.form.Tags = &tags
	return b
}
