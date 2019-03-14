package model

import "time"

// Article structure definition
type Article struct {
	ID          uint32     `json:"id,omitempty"`
	CategoryID  *uint32    `json:"category_id,omitempty"`
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
