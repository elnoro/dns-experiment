package main

import (
	"context"
	"fmt"
	"github.com/elnoro/foxylock/m/v2/admin"
	"github.com/elnoro/foxylock/m/v2/db"
	"log"
)

func main() {
	fmt.Println("test")

	s := admin.NewHttpServer(db.NewInMemory(), "localhost:8081")

	go log.Fatal(s.Run(context.Background()))

	select {}
}
