package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
)

func main() {
	filename := flag.String("config", "config.yaml", "Configuration file")
	flag.Parse()

	cfg := LoadConfiguration(*filename)
	logger := LoadLogger(cfg)

	bot, err := New(cfg, logger)
	if err != nil {
		logger.WithField("context", "run").Fatal(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		logger.Info("bot started")
		defer logger.Info("bot stopped")

		if err := bot.Run(); err != nil {
			logger.WithField("context", "run").Fatal(err)
		}
	}()

	<-stop

	timeout := cfg.GetDuration("app.graceful_shutdown_timeout")
	logger.WithField("timeout", timeout).Info("graceful shutdown...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	go func() {
		defer cancel()

		if err := bot.Stop(); err != nil {
			logger.WithField("context", "stop").Fatal(err)
		}
	}()

	<-ctx.Done()
	if ctx.Err() == context.DeadlineExceeded {
		logger.WithField("context", "graceful_shutdown").Error("context timeout exceeded")
	}
}
