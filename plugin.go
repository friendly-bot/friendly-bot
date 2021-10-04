package main

import (
	"fmt"
	"plugin"
	"time"

	"github.com/friendly-bot/friendly-bot/api"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

const (
	NewOnMessage  = "NewOnMessage"
	ReactionAdded = "NewReactionAdded"
	NewJob        = "NewJob"
)

type (
	PluginStore struct {
		cron *cron.Cron

		onMessage map[string]api.OnMessage
		onReactionAdded map[string]api.OnReactionAdded
	}

	PluginConfiguration struct {
		Path     string
		Schedule string
		Disable  bool
	}
)

func (p PluginStore) count() int {
	return len(p.cron.Entries()) + len(p.onMessage)
}

func (b *Bot) loadPluginStore() error {
	b.logger.WithField("context", "init").Info("plugin_store")

	if err := b.loadPluginsOnMessage(); err != nil {
		return err
	}

	if err := b.loadPluginsOnReactionAdded(); err != nil {
		return err
	}

	if err := b.loadPluginsCronjob(); err != nil {
		return err
	}

	return nil
}

func (b *Bot) loadPluginsOnMessage() error {
	b.logger.WithField("context", "init").Info("plugin on_message")
	b.plugins.onMessage = make(map[string]api.OnMessage)

	var cfgs map[string]PluginConfiguration

	if err := b.config.UnmarshalKey("bot.plugin.event.on_message", &cfgs); err != nil {
		return fmt.Errorf("can't unmarshal 'bot.plugin.event.on_message': %w", err)
	}

	for name, cfg := range cfgs {
		if cfg.Disable {
			continue
		}

		logger := b.logger.WithFields(logrus.Fields{"path": cfg.Path, "name": name})
		logger.Info("load plugin")

		sym, err := lookup(NewOnMessage, cfg)
		if err != nil {
			logger.WithField("context", "lookup").Error(err)
			continue
		}

		n, ok := sym.(api.NewOnMessage)
		if !ok {
			logger.WithField("context", "assertion").Error("can't cast to 'api.NewOnMessage'")
			continue
		}

		onMessage, err := n(b.config.Sub(fmt.Sprintf("bot.plugin.event.on_message.%s.configuration", name)))
		if err != nil {
			logger.WithField("context", "new").Error(err)
			continue
		}

		b.plugins.onMessage[name] = onMessage
	}

	return nil
}

func (b *Bot) loadPluginsOnReactionAdded() error {
	b.logger.WithField("context", "init").Info("plugin on_reaction")
	b.plugins.onReactionAdded = make(map[string]api.OnReactionAdded)

	var cfgs map[string]PluginConfiguration

	if err := b.config.UnmarshalKey("bot.plugin.event.on_reaction_added", &cfgs); err != nil {
		return fmt.Errorf("can't unmarshal 'bot.plugin.event.on_reaction_added': %w", err)
	}

	for name, cfg := range cfgs {
		if cfg.Disable {
			continue
		}

		logger := b.logger.WithFields(logrus.Fields{"path": cfg.Path, "name": name})
		logger.Info("load plugin")

		sym, err := lookup(ReactionAdded, cfg)
		if err != nil {
			logger.WithField("context", "lookup").Error(err)
			continue
		}

		n, ok := sym.(api.NewOnReactionAdded)
		if !ok {
			logger.WithField("context", "assertion").Error("can't cast to 'api.OnReactionAdded'")
			continue
		}

		onReactionAdded, err := n(b.config.Sub(fmt.Sprintf("bot.plugin.event.on_reaction_added.%s.configuration", name)))
		if err != nil {
			logger.WithField("context", "new").Error(err)
			continue
		}

		b.plugins.onReactionAdded[name] = onReactionAdded
	}

	return nil
}

func (b *Bot) loadPluginsCronjob() error {
	b.logger.WithField("context", "init").Info("plugin cronjob")

	loc, err := time.LoadLocation(b.config.GetString("app.timezone"))
	if err != nil {
		return err
	}
	b.plugins.cron = cron.New(cron.WithLocation(loc))

	var cfgs map[string]PluginConfiguration

	if err := b.config.UnmarshalKey("bot.plugin.event.cronjob", &cfgs); err != nil {
		return fmt.Errorf("can't unmarshal 'bot.plugin.event.cronjob': %w", err)
	}

	for name, cfg := range cfgs {
		if cfg.Disable {
			continue
		}

		logger := b.logger.WithFields(logrus.Fields{"path": cfg.Path, "name": name, "schedule": cfg.Schedule})
		logger.Info("load plugin")

		sym, err := lookup(NewJob, cfg)
		if err != nil {
			logger.WithField("context", "lookup").Error(err)
			continue
		}

		n, ok := sym.(api.NewJob)
		if !ok {
			logger.WithField("context", "assertion").Error("can't cast to 'api.NewJob'")
			continue
		}

		job, err := n(b.config.Sub(fmt.Sprintf("bot.plugin.event.cronjob.%s.configuration", name)))
		if err != nil {
			logger.WithField("context", "new").Error(err)
			continue
		}

		_, err = b.plugins.cron.AddFunc(cfg.Schedule, b.wrapJob(name, job))
		if err != nil {
			logger.WithField("context", "add_func").Error(err)
			continue
		}

		logger.Info("plugin loaded")
	}

	return nil
}

func lookup(f string, cfg PluginConfiguration) (plugin.Symbol, error) {
	p, err := plugin.Open(cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "can't open", err)
	}

	sym, err := p.Lookup(f)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "can't lookup", err)
	}

	return sym, nil
}

func (b *Bot) wrapJob(name string, r api.Runner) func() {
	return func() {
		if err := r.Run(b.newContext(name)); err != nil {
			b.logger.WithFields(logrus.Fields{
				"name":    name,
				"context": "job",
			}).Error(err)
		}
	}
}
