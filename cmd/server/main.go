package main

import (
	"log"
	"net"

	"github.com/pu4mane/goChat/internal/app/broker"
	"github.com/pu4mane/goChat/internal/app/room"
)

func main() {
	listen, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Server is running")
	}

	ns, err := broker.NewNATS("localhost:4222")
	if err != nil {
		log.Fatal()
	}

	r := room.NewRoom("chat", ns)

	go r.Broadcast()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go r.HandleClient(conn.RemoteAddr().String(), conn)
	}
}
