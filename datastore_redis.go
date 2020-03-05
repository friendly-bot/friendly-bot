package main

import (
	"time"

	"github.com/friendly-bot/friendly-bot/api"
	"github.com/go-redis/redis/v7"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const CacheRedis = "redis"

type (
	RedisDataStore struct {
		logger logrus.FieldLogger

		redis *redis.Client
	}
)

func NewRedisCache(cfg *viper.Viper, l logrus.FieldLogger) (GlobalDataStore, error) {
	r := &RedisDataStore{
		logger: l,
		redis: redis.NewClient(&redis.Options{
			Addr:         cfg.GetString("addr"),
			Password:     cfg.GetString("addr"),
			DialTimeout:  cfg.GetDuration("timeout.dial"),
			ReadTimeout:  cfg.GetDuration("timeout.read"),
			WriteTimeout: cfg.GetDuration("timeout.write"),
		}),
	}

	if cfg.GetBool("flush_on_start") {
		if err := r.redis.FlushAll(); err != nil {
			r.logger.WithField("context", "flush_all").Error(err)
		}
	}

	return r, r.redis.Ping().Err()
}

func (s *RedisDataStore) Scope(scope string) api.DataStore {
	return &DataStoreScoped{
		DataStore: s,
		scope:     scope,
	}
}

func (s *RedisDataStore) Get(key string) string {
	s.logger.WithField("key", key).Debug("get")

	return s.redis.Get(key).Val()
}

func (s *RedisDataStore) Set(key, value string, ttl time.Duration) {
	s.logger.WithFields(logrus.Fields{
		"key":   key,
		"value": value,
		"ttl":   ttl,
	}).Debug("set")

	s.redis.Set(key, value, ttl)
}

func (s *RedisDataStore) Exist(key string) bool {
	s.logger.WithField("key", key).Debug("exist")

	return s.redis.Exists(key).Val() > 0
}
