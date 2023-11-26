package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/gorilla/mux"
	zlog "github.com/rs/zerolog/log"

	"github.com/FlipTheBoard/Distributed-Backup-System/server/config"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Привет. Это главная страница")
}

func Run(ctx context.Context, config *config.Config) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	log := zlog.Ctx(ctx)

	router := mux.NewRouter()
	handler := http.StripPrefix("/files/", http.FileServer(http.Dir(config.BackupsDir)))

	router.HandleFunc("/", Index)
	router.PathPrefix("/files/").Handler(handler)

	server := &http.Server{Addr: config.ListenAddr, Handler: router}

	log.Info().Msg(fmt.Sprintf("starting HTTP server on %s ...", config.ListenAddr))

	go func() {
		<-ctx.Done()

		log.Info().Msgf("gracefully stoping HTTP server...")
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal().Err(err).Send()
		}
	}()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal().Err(err).Send()
	}

	return nil
}
