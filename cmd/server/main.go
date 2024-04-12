package main

import (
	"log"
	"net"

	"github.com/pu4mane/goChat/internal/app/room"
)

func main() {
	cs := room.NewChatRoom()

	listen, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Server is running")
	}

	go cs.Broadcaster()

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go cs.HandleClient(conn.RemoteAddr().String(), conn)
	}
}
