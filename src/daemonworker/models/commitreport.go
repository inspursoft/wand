package models

import (
	"log"
	"sync"
)

type CommitReport struct {
	CommitID string `json:"commit_id"`
	Report   string `json:"report"`
}

type CachedReport struct {
	cache map[string]*CommitReport
	mutex sync.Mutex
}

func NewCachedReport() *CachedReport {
	return &CachedReport{
		cache: make(map[string]*CommitReport),
	}
}

func (c *CachedReport) Add(commitReport *CommitReport) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache[commitReport.CommitID] = commitReport
	log.Printf("Stored commit report into cache: %+v\n", commitReport)
}

func (c *CachedReport) Get(commitID string) (report *CommitReport, found bool) {
	report, found = c.cache[commitID]
	return
}
