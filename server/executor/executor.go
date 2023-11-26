package executor

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	zlog "github.com/rs/zerolog/log"

	"github.com/FlipTheBoard/Distributed-Backup-System/server/config"
)

var cmdMutex sync.Mutex

func Run(ctx context.Context, config *config.Config) error {
	for name, backup := range config.Backups {
		go startBackupRunner(ctx, name, *backup, config.BackupsDir)
	}
	return nil
}

func startBackupRunner(ctx context.Context, name string, backup config.Backup, dir string) {
	log := zlog.Ctx(ctx)
	log = ptr(log.With().
		Str("backup_name", name).
		Dur("duration", backup.Interval).
		Logger())
	log.Info().Msg("starting backup runner...")

	t := time.NewTicker(backup.Interval)
	defer t.Stop()

	for ; true; <-t.C {
		log.Debug().Msg("running backup...")
		cmdMutex.Lock()

		for _, command := range backup.Commands {
			cmd := formatCommand(command, dir, name)

			stdout, err := exec.Command("/bin/bash", "-c", cmd).Output()
			if err != nil {
				log.Err(err).Send()
			} else {
				log.Debug().
					Bytes("stdout", stdout).
					Msg("success")
			}
		}

		cmdMutex.Unlock()
	}
}

func formatCommand(cmd string, dir string, name string) string {
	match := map[string]func() string{
		"{dir}":  func() string { return dir },
		"{ts}":   func() string { return fmt.Sprintf("%v", time.Now().Unix()) },
		"{dt}":   func() string { return time.Now().Format("2006-01-02_15:04:05") },
		"{name}": func() string { return name },
	}

	for key, fn := range match {
		cmd = strings.ReplaceAll(cmd, key, fn())
	}

	return cmd
}

func ptr[T any](v T) *T {
	return &v
}
