package models

import "fmt"

type CommitReport struct {
	CommitID string `json:"commit_id"`
	Report   string `json:"report"`
}

func ToCommitReport(config *CachedConfig) *CommitReport {
	return &CommitReport{
		CommitID: config.Key,
		Report:   config.Value,
	}
}

func (c *CommitReport) ToCachedConfig() *CachedConfig {
	return &CachedConfig{
		Key:   fmt.Sprintf("COMMITS_%s", c.CommitID),
		Value: c.Report,
	}
}
