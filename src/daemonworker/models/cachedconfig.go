package models

import (
	"log"
	"sync"
)

type CachedConfig struct {
	Key   string `json:"cached_key"`
	Value string `json:"cached_val"`
}

type CachedStore struct {
	cache map[string]*CachedConfig
	mutex sync.Mutex
}

func NewCachedStore() *CachedStore {
	return &CachedStore{
		cache: make(map[string]*CachedConfig),
	}
}

func (c *CachedStore) Add(config *CachedConfig) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.cache[config.Key] = config
	log.Printf("Stored config into cache: %+v\n", config)
}

func (c *CachedStore) Get(key string) (config *CachedConfig, found bool) {
	config, found = c.cache[key]
	return
}

func (c *CachedStore) All() map[string]*CachedConfig {
	return c.cache
}
