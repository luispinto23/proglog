package main

import (
	"log"

	"github.com/luispinto23/proglog/internal/server"
)

func main() {
	srv := server.NewHttpServer(":8080")
	log.Fatal(srv.ListenAndServe())
}
