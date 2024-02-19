package main

import (
	"context"
	"github.com/g82411/gd/cmd"
	"github.com/g82411/gd/utils/googleDrive"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Context(context.Background())
	srv, err := googleDrive.GetService(context.Background())
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}
	srvContext := context.WithValue(ctx, "service", srv)
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := cmd.Main(srvContext, os.Args); err != nil {
		os.Exit(1)
	}
}
