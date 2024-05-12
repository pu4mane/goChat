package main

import (
	"log"
	"runtime"
	"time"

	"github.com/pu4mane/goChat/internal/app/broker"
	"github.com/pu4mane/goChat/internal/app/server"
)

func main() {
	ns, err := broker.NewNATS("localhost:4222")
	if err != nil {
		log.Fatal()
	}

	srv := server.Server{
		Addr:         "localhost:9090",
		IdleTimeout:  30 * time.Minute,
		MaxReadBytes: 1000,
	}
	go srv.ListenAndServe(ns)
	for {
		runtime.Gosched()
	}
}
