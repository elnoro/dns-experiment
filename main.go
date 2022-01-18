package main

import (
	"context"
	"github.com/elnoro/foxylock/m/v2/cmd"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	app, err := cmd.NewProduction()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	err = app.Start(ctx)
	if err != nil {
		cancel()
		log.Fatal(err)
	}

	waitForExit(cancel)
}

func waitForExit(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	cancel()
	log.Println("Exiting...")
}
