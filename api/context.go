package api

import (
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

type (
	Context struct {
		RTM    *slack.RTM
		Logger logrus.FieldLogger

		Cache DataStore
	}
)
