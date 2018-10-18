package models

type CommitReport struct {
	CommitID string `json:"commit_id"`
	Report   string `json:"report"`
}
