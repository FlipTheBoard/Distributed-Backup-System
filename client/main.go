package main

import (
	"context"
	"log"
	"os"

	"github.com/FlipTheBoard/Distributed-Backup-System/client/client"
)

func main() {
	ctx := context.Background()

	backupDir, ok := os.LookupEnv("FTB_BACKUP_DIR")
	if !ok {
		log.Fatal("failed to get env variable FTB_BACKUP_DIR")
	}

	serverAddr, ok := os.LookupEnv("FTB_SERVER_ADDR")
	if !ok {
		log.Fatal("failed to get env variable FTB_SERVER_ADDR")
	}

	if err := client.Run(ctx, backupDir, serverAddr); err != nil {
		log.Fatal(err)
	}
}
