package server

import (
	"log"
	"net"
	"time"

	"github.com/pu4mane/goChat/internal/app/broker"
	"github.com/pu4mane/goChat/internal/app/room"
)

type Server struct {
	Addr         string
	IdleTimeout  time.Duration
	MaxReadBytes int64

	listener net.Listener
	room     *room.Room
}

func (srv *Server) ListenAndServe(broker broker.MessageBroker) error {
	addr := srv.Addr
	if addr == "" {
		addr = "localhost:9090"
	}
	log.Printf("Server is running on %v\n", addr)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	room := room.NewRoom("chat", broker)
	srv.room = room

	go room.Broadcast()

	srv.listener = listener
	for {
		newConn, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting connection %v", err)
			continue
		}
		log.Printf("accepted connection from %v", newConn.RemoteAddr())

		conn := &conn{
			Conn:          newConn,
			IdleTimeout:   srv.IdleTimeout,
			MaxReadBuffer: srv.MaxReadBytes,
		}

		conn.SetDeadline(time.Now().Add(conn.IdleTimeout))

		go room.HandleClient(conn.RemoteAddr().String(), conn)
	}
}
