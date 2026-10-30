package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"crcs/platform/messaging/broker/kafka"
	"crcs/platform/messaging/publisher"
	"crcs/platform/o11y/logging"
	"crcs/platform/store/postgrex"

	"github.com/charmingruby/crcs/services/outbox/internal/outbox"
	"github.com/charmingruby/crcs/services/outbox/internal/platform/config"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	log := logging.InitLogger()

	cfg, err := config.Load()
	if err != nil {
		log.Error("config: loading failed", "error", err)
		return err
	}

	db, err := postgrex.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("postgres: connection failed", "error", err)
		return err
	}
	defer db.Close()

	writer := kafka.NewWriter(cfg.KafkaBrokers)
	defer writer.Close()

	log.Info("outbox: relay started")

	relay := outbox.New(db, publisher.NewKafka(writer))

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		for ctx.Err() == nil {
			published, err := relay.PublishNext(ctx)

			if err != nil {
				if ctx.Err() == nil {
					log.Error("outbox: publish failed; retrying next poll", "error", err)
				}

				break
			}

			if !published {
				break
			}
		}

		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}
