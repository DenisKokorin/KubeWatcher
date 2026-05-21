package models

import "time"

type ExternalPost struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Title       string    `json:"title"`
	Summary     string    `json:"summary"`
	Content     string    `json:"content"`
	Source      string    `json:"source"`
	PublishedAt time.Time `json:"published_at"`
}
