package main

import (
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:9090")
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan struct{})

	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			log.Fatal("Copying error:", err)
		} else {
			log.Fatal("Copying complete")
		}
		done <- struct{}{}
	}()

	go func() {
		_, err := io.Copy(conn, os.Stdin)
		if err != err {
			log.Fatal("Error copying:", err)
		}
	}()

	<-done
	conn.Close()
}
