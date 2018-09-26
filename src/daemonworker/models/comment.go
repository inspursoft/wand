package models

import "time"

// Comment represents a comment in commit and issue page.
type Comment struct {
	ID      int64     `json:"id"`
	HTMLURL string    `json:"html_url"`
	Poster  *User     `json:"user"`
	Body    string    `json:"body"`
	Created time.Time `json:"created_at"`
	Updated time.Time `json:"updated_at"`
}
