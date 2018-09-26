package models

import "time"

// Release represents a release API object.
type Release struct {
	ID              int64     `json:"id"`
	TagName         string    `json:"tag_name"`
	TargetCommitish string    `json:"target_commitish"`
	Name            string    `json:"name"`
	Body            string    `json:"body"`
	Draft           bool      `json:"draft"`
	Prerelease      bool      `json:"prerelease"`
	Author          *User     `json:"author"`
	Created         time.Time `json:"created_at"`
}
