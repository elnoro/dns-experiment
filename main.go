package main

import (
	"context"
	"github.com/elnoro/foxylock/m/v2/admin"
	"github.com/elnoro/foxylock/m/v2/db"
	"log"
)

func main() {

	inMemoryDb := db.NewInMemory()

	rs := admin.NewRedisLikeServer(inMemoryDb, ":6379", "testpass")
	startServer(rs)
	gs := admin.NewHttpServer(inMemoryDb, ":8081")
	startServer(gs)

	select {}
}

func startServer(s admin.DbServer) {
	go func(s admin.DbServer) {
		log.Fatal(s.Run(context.Background()))
	}(s)
}
