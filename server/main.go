package main

import (
	"context"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"

	"github.com/FlipTheBoard/Distributed-Backup-System/server/config"
	"github.com/FlipTheBoard/Distributed-Backup-System/server/executor"
	"github.com/FlipTheBoard/Distributed-Backup-System/server/server"
)

func main() {
	cfg, err := config.ParseConfig()
	if err != nil {
		zlog.Fatal().Err(err).Send()
	}

	ctx := context.Background()

	log := zlog.
		Output(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) { w.TimeFormat = time.RFC3339Nano })).
		Level(cfg.LoggingLevel)

	err = config.Log(ctx, cfg)
	if err != nil {
		zlog.Fatal().Err(err).Send()
	}

	ctx = log.WithContext(ctx)

	if err = executor.Run(ctx, cfg); err != nil {
		log.Fatal().Err(err).Send()
	}

	if err = server.Run(ctx, cfg); err != nil {
		log.Fatal().Err(err).Send()
	}
}
