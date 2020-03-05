package main

import (
	"time"

	"github.com/friendly-bot/friendly-bot/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const CacheInMemory = "in_memory"

type (
	MemoryDataStore struct {
		logger logrus.FieldLogger

		data map[string]string
	}
)

func NewInMemoryCache(_ *viper.Viper, l logrus.FieldLogger) (GlobalDataStore, error) {
	return &MemoryDataStore{
		logger: l,
		data:   map[string]string{},
	}, nil
}

func (s *MemoryDataStore) Scope(scope string) api.DataStore {
	return &DataStoreScoped{
		DataStore: s,
		scope:     scope,
	}
}

func (s *MemoryDataStore) Get(key string) string {
	s.logger.WithField("key", key).Debug("get")

	return s.data[key]
}

func (s *MemoryDataStore) Set(key, value string, _ time.Duration) {
	s.logger.WithFields(logrus.Fields{
		"key":   key,
		"value": value,
	}).Debug("set")

	s.data[key] = value
}

func (s *MemoryDataStore) Exist(key string) bool {
	s.logger.WithField("key", key).Debug("exist")

	_, ok := s.data[key]

	return ok
}
