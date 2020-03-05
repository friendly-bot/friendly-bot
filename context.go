package main

import "github.com/friendly-bot/friendly-bot/api"

func (b *Bot) newContext(name string) api.Context {
	ctx, ok := b.contextCache[name]
	if ok {
		return ctx
	}

	b.contextCache[name] = api.Context{
		RTM:    b.rtm,
		Logger: b.logger.WithField("plugin", name),
		Cache:  b.cache.Scope(name),
	}

	return b.contextCache[name]
}
