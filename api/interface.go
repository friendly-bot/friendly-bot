package api

import (
	"time"

	"github.com/slack-go/slack"
	"github.com/spf13/viper"
)

type (
	DataStore interface {
		Get(string) string
		Set(string, string, time.Duration)
		Exist(string) bool
	}

	NewOnMessage = func(*viper.Viper) (OnMessage, error)
	OnMessage    interface {
		OnMessage(*slack.MessageEvent, Context) error
	}

	NewOnReactionAdded = func(*viper.Viper) (OnReactionAdded, error)
	OnReactionAdded interface {
		OnReactionAdded(*slack.ReactionAddedEvent, Context) error
	}

	NewJob = func(*viper.Viper) (Runner, error)
	Runner interface {
		Run(Context) error
	}
)
