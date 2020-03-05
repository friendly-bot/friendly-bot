package main

import (
	"net/http"

	"github.com/friendly-bot/friendly-bot/api"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/spf13/viper"
)

type (
	Bot struct {
		logger logrus.FieldLogger
		config *viper.Viper

		client *slack.Client
		rtm    *slack.RTM

		quit chan bool
		done chan bool

		plugins PluginStore
		cache   GlobalDataStore

		contextCache map[string]api.Context
	}
)

func New(cfg *viper.Viper, l logrus.FieldLogger) (*Bot, error) {
	client := slack.New(
		cfg.GetString("bot.slack.token"),
		slack.OptionDebug(cfg.GetBool("bot.slack.debug")),
		slack.OptionHTTPClient(&http.Client{
			Timeout: cfg.GetDuration("bot.slack.http.timeout"),
		}),
	)

	b := &Bot{
		logger: l,
		config: cfg,
		client: client,
		quit:   make(chan bool, 1),
		done:   make(chan bool, 1),
		rtm:    client.NewRTM(),
	}

	return b, b.init()
}

func (b *Bot) init() error {
	b.logger.WithField("context", "init").Info("initializing...")

	if err := b.loadCache(); err != nil {
		return err
	}

	if err := b.loadPluginStore(); err != nil {
		return err
	}

	b.contextCache = make(map[string]api.Context, b.plugins.count())

	return nil
}

func (b *Bot) Run() error {
	go b.rtm.ManageConnection()
	b.plugins.cron.Start()

	for {
		select {
		case e := <-b.rtm.IncomingEvents:
			b.broadcast(e)
		case <-b.quit:
			err := b.rtm.Disconnect()
			<-b.plugins.cron.Stop().Done()
			b.done <- true
			return err
		}
	}
}

func (b *Bot) broadcast(e slack.RTMEvent) {
	switch ev := e.Data.(type) {

	case *slack.MessageEvent:
		for name, plg := range b.plugins.onMessage {
			if err := plg.OnMessage(ev, b.newContext(name)); err != nil {
				b.logger.WithField("plugin_name", name).Error(err)
			}
		}

	case *slack.RTMError:
		b.logger.WithField("event", "rtm_error").Error(ev.Error())
	}
}

func (b *Bot) Stop() error {
	b.quit <- true
	<-b.done

	return nil
}
