package models

import "time"

type ArticleResponse struct {
	ID          int           `json:"id"`
	Title       string        `json:"title"`
	Description string        `json:"description,omitempty"`
	Content     string        `json:"content"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	User        *UserResponse `json:"user,omitempty"`
}
