package main

import (
	"context"
	"github.com/elnoro/foxylock/m/v2/admin"
	"github.com/elnoro/foxylock/m/v2/config"
	"github.com/elnoro/foxylock/m/v2/db"
	coredns_integration "github.com/elnoro/foxylock/m/v2/dns/coredns-integration"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	inMemoryDb := db.NewInMemory()

	ctx, cancel := context.WithCancel(context.Background())

	rs := admin.NewRedisLikeServer(inMemoryDb, c.RedisAddr, c.RedisPass)
	startServer(rs, ctx)

	if c.HttpAddr != "" {
		gs := admin.NewHttpServer(inMemoryDb, c.HttpAddr)
		startServer(gs, ctx)
	}

	err = coredns_integration.
		NewCoreDns(inMemoryDb).
		Start()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Running...")

	waitForExit(cancel)
}

func waitForExit(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	cancel()
	log.Println("Exiting...")
}

func startServer(s admin.DbServer, ctx context.Context) {
	go func(s admin.DbServer) {
		log.Fatal(s.Run(ctx))
	}(s)
}
