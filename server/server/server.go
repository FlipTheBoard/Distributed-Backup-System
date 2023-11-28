package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"

	"github.com/gorilla/mux"
	zlog "github.com/rs/zerolog/log"

	"github.com/FlipTheBoard/Distributed-Backup-System/server/config"
)

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Привет. Это главная страница")
}

func GetBackupFiles(ctx context.Context, config *config.Config) http.HandlerFunc {
	cmd := fmt.Sprintf(
		"(cd %s && find . -type f | cut -c 2- | sort)",
		config.BackupsDir,
	)

	return func(w http.ResponseWriter, r *http.Request) {
		stdout, err := exec.Command("/bin/bash", "-c", cmd).Output()
		if err != nil {
			fmt.Fprintf(w, "err: %v", err)
		} else {
			fmt.Fprint(w, string(stdout))
		}

	}
}

func GetBackupDirs(ctx context.Context, config *config.Config) http.HandlerFunc {
	cmd := fmt.Sprintf(
		"(cd %s && find . -type d | cut -c 2- | sort | tail -n +2)",
		config.BackupsDir,
	)

	return func(w http.ResponseWriter, r *http.Request) {
		stdout, err := exec.Command("/bin/bash", "-c", cmd).Output()
		if err != nil {
			fmt.Fprintf(w, "err: %v", err)
		} else {
			fmt.Fprint(w, string(stdout))
		}

	}
}

func Run(ctx context.Context, config *config.Config) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt)
	defer cancel()

	log := zlog.Ctx(ctx)

	router := mux.NewRouter()
	handler := http.StripPrefix("/backups/", http.FileServer(http.Dir(config.BackupsDir)))

	router.HandleFunc("/", Index)
	router.HandleFunc("/files/", GetBackupFiles(ctx, config))
	router.HandleFunc("/dirs/", GetBackupDirs(ctx, config))
	router.PathPrefix("/backups/").Handler(handler)

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
