package main

import (
	"context"
	"crcs/platform/httpx"
	"crcs/platform/messaging/broker/kafka"
	"crcs/platform/messaging/publisher"

	"crcs/platform/o11y/logging"
	"crcs/platform/store/postgrex"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/charmingruby/crcs/services/api/internal/identity"
	"github.com/charmingruby/crcs/services/api/internal/matchmaking"
	"github.com/charmingruby/crcs/services/api/internal/platform/auth"
	"github.com/charmingruby/crcs/services/api/internal/platform/config"

	"github.com/charmingruby/crcs/services/api/internal/roster"
	"github.com/charmingruby/crcs/services/api/internal/task"
)

const shutdownTimeout = 30 * time.Second

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

	log.Info("kafka: creating producer...")

	kafkaWriter := kafka.NewWriter(cfg.KafkaBrokers)

	log.Info("kafka: producer created")

	pub := publisher.NewKafka(kafkaWriter)

	tokenMgr := auth.NewTokenManager(cfg.JWTSecret, cfg.JWTTTL)

	srv, router := httpx.NewServer(cfg.Port)

	if err := identity.New(router, db, tokenMgr); err != nil {
		log.Error("auth: module error", "error", err)
		return err
	}

	if err := matchmaking.New(router, db, pub, tokenMgr); err != nil {
		log.Error("matchmaking: module error", "error", err)
		return err
	}

	seasonReader, err := matchmaking.NewSeasonReader(db)
	if err != nil {
		log.Error("matchmaking: season reader error", "error", err)
		return err
	}

	if err := roster.New(router, db, pub, seasonReader, tokenMgr); err != nil {
		log.Error("roster: module error", "error", err)
		return err
	}

	if err := task.New(router, db, tokenMgr); err != nil {
		log.Error("task: module error", "error", err)
		return err
	}

	shutdownErrCh := make(chan error, 1)
	go shutdown(ctx, log, shutdownErrCh, srv, db, kafkaWriter)

	log.Info("server: running...", "port", cfg.Port)
	if err := srv.Start(); err != nil {
		log.Error("server: start error", "error", err)

		stop()
		<-shutdownErrCh
		return err
	}

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
	srv *httpx.Server,
	db *sqlx.DB,
	writer *kafka.Writer,
) {
	<-ctx.Done()
	log.Info("shutdown: signal received, starting graceful shutdown...")

	ctxTimeout, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	var err error

	if closeErr := srv.Close(ctxTimeout); closeErr != nil {
		if errors.Is(closeErr, context.DeadlineExceeded) {
			closeErr = errors.New("http server: deadline exceeded, forcing shutdown")
		}

		err = errors.Join(err, closeErr)
	}

	if closeErr := db.Close(); closeErr != nil {
		err = errors.Join(err, closeErr)
	}

	if closeErr := writer.Close(); closeErr != nil {
		err = errors.Join(err, closeErr)
	}

	errCh <- err
}
