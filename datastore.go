package main

import (
	"fmt"
	"time"

	"github.com/friendly-bot/friendly-bot/api"
)

type (
	GlobalDataStore interface {
		Scope(string) api.DataStore
	}

	DataStoreScoped struct {
		api.DataStore
		scope string
	}
)

func (b *Bot) loadCache() (err error) {
	b.logger.WithField("context", "init").Info("cache")

	t := b.config.GetString("bot.cache.type")
	cfg := b.config.Sub("bot.cache.configuration")
	logger := b.logger.WithField("cache", t)

	switch t {
	case CacheInMemory:
		b.cache, err = NewInMemoryCache(cfg, logger)
	case CacheRedis:
		b.cache, err = NewRedisCache(cfg, logger)
	default:
		err = fmt.Errorf("%s: %w", t, UnknownDatastoreErr)
	}

	return
}

func (s *DataStoreScoped) Get(key string) string {
	return s.DataStore.Get(fmt.Sprintf("%s.%s", s.scope, key))
}

func (s *DataStoreScoped) Set(key, value string, ttl time.Duration) {
	s.DataStore.Set(fmt.Sprintf("%s.%s", s.scope, key), value, ttl)
}

func (s *DataStoreScoped) Exist(key string) bool {
	return s.DataStore.Exist(fmt.Sprintf("%s.%s", s.scope, key))
}
