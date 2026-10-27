package main

import (
	"context"
	"crcs/event/contenders_ingest_v1"
	"crcs/platform/broker/listener"
	"crcs/platform/kafka"
	"crcs/platform/logging"
	"crcs/platform/postgrex"

	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/charmingruby/crcs/services/worker/internal/ingest"
	"github.com/charmingruby/crcs/services/worker/internal/platform/config"
	"github.com/jmoiron/sqlx"
)

func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer stop()

	log := logging.InitLogger()

	cfg, err := config.Load()
	if err != nil {
		log.Error("config: error loading config", "error", err)
		return err
	}

	log.Info("postgres: connecting...")

	db, err := postgrex.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("postgres: connection error", "error", err)
		return err
	}

	log.Info("postgres: connected")

	listeners, err := ingest.New(db)
	if err != nil {
		log.Error("ingest: module error", "error", err)
		return err
	}

	shutdownErrCh := make(chan error, 1)

	go func() {
		log.Info("listener: started consuming messages",
			"topic", contenders_ingest_v1.Topic,
		)

		if err := listener.NewKafka[contenders_ingest_v1.Message](kafka.NewReader(
			cfg.KafkaBrokers,
			contenders_ingest_v1.Topic,
			cfg.ServiceName,
		)).Listen(ctx, listeners.OnContendersIngestV1.Handle); err != nil {
			log.Error("listener: error on consuming messages",
				"error", err,
				"topic", contenders_ingest_v1.Topic,
			)

			shutdownErrCh <- err
		}
	}()

	go shutdown(ctx, log, shutdownErrCh, db)

	if err := <-shutdownErrCh; err != nil {
		log.Error("shutdown: shutdown error", "error", err)
		return err
	}

	log.Info("shutdown: gracefully shutdown")
	return nil
}

func shutdown(
	ctx context.Context,
	log *slog.Logger,
	errCh chan error,
	db *sqlx.DB,

) {
	<-ctx.Done()
	log.Info("shutdown: signal received, starting graceful shutdown...")

	var err error

	if closeErr := db.Close(); closeErr != nil {
		err = errors.Join(err, closeErr)
	}

	errCh <- err
}
