package main

import (
	"context"
	"log"

	"github.com/FlipTheBoard/Distributed-Backup-System/client/client"
	"github.com/FlipTheBoard/Distributed-Backup-System/client/config"
)

func main() {
	ctx := context.Background()

	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err = client.Run(ctx, cfg); err != nil {
		log.Fatal(err)
	}
}
