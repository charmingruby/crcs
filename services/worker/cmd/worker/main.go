package main

import (
	"context"
	"crcs/event/contenders_submitted_v1"
	"crcs/event/season_matchmaking_requested_v1"

	"crcs/platform/messaging/broker/kafka"
	"crcs/platform/messaging/listener"
	"crcs/platform/o11y/logging"
	"crcs/platform/store/postgrex"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/charmingruby/crcs/services/worker/internal/ingest"
	"github.com/charmingruby/crcs/services/worker/internal/matchmaking"
	"github.com/charmingruby/crcs/services/worker/internal/platform/config"
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

	ingestListeners, err := ingest.New(db)
	if err != nil {
		log.Error("ingest: module error", "error", err)
		return err
	}

	matchmakingListeners, err := matchmaking.New(db)
	if err != nil {
		log.Error("matchmaking: module error", "error", err)
		return err
	}

	var wg sync.WaitGroup

	shutdownErrCh := make(chan error, 1)

	readers := []*kafka.Reader{
		listen(ctx, log, cfg, &wg, shutdownErrCh,
			contenders_submitted_v1.Topic,
			ingestListeners.OnContendersSubmittedV1,
		),
		listen(ctx, log, cfg, &wg, shutdownErrCh,
			season_matchmaking_requested_v1.Topic,
			matchmakingListeners.OnSeasonMatchmakingRequestedV1,
		),
	}

	go shutdown(ctx, log, shutdownErrCh, db, readers, &wg)

	if err := <-shutdownErrCh; err != nil {
		log.Error("shutdown: shutdown error", "error", err)
		return err
	}

	log.Info("shutdown: gracefully shutdown")
	return nil
}

func listen[T any](
	ctx context.Context,
	log *slog.Logger,
	cfg *config.Config,
	wg *sync.WaitGroup,
	errCh chan<- error,
	topic string,
	handle func(context.Context, T) error,
) *kafka.Reader {
	reader := kafka.NewReader(cfg.KafkaBrokers, topic, cfg.ServiceName)

	wg.Go(func() {
		log.Info("listener: started consuming messages",
			"topic", topic,
		)

		if err := listener.NewKafka[T](reader).Listen(ctx, handle); err != nil {
			log.Error("listener: error on consuming messages",
				"error", err,
				"topic", topic,
			)

			errCh <- err
		}
	})

	return reader
}

func shutdown(
	ctx context.Context,
	log *slog.Logger,
	errCh chan<- error,
	db *sqlx.DB,
	readers []*kafka.Reader,
	wg *sync.WaitGroup,
) {
	<-ctx.Done()
	log.Info("shutdown: signal received, starting graceful shutdown...")

	ctxTimeout, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	var err error

	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctxTimeout.Done():
		err = errors.Join(err, errors.New("listeners: deadline exceeded, forcing shutdown"))
	}

	for _, reader := range readers {
		if closeErr := reader.Close(); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}

	if closeErr := db.Close(); closeErr != nil {
		err = errors.Join(err, closeErr)
	}

	errCh <- err
}
