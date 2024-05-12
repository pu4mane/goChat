package room

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/pu4mane/goChat/internal/app/broker"
	"github.com/pu4mane/goChat/internal/app/model"
)

const MAX_MESSAGE_LENGTH = 1000

type Room struct {
	name     string
	mu       sync.RWMutex
	clients  map[string]net.Conn
	messages chan *model.Message
	broker   broker.MessageBroker
}

func NewRoom(name string, broker broker.MessageBroker) *Room {
	return &Room{
		name:     name,
		mu:       sync.RWMutex{},
		clients:  make(map[string]net.Conn),
		messages: make(chan *model.Message),
		broker:   broker,
	}
}

func (r *Room) addClient(ID string, conn net.Conn) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.clients[ID] = conn
}

func (r *Room) removeClient(ID string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.clients, ID)
}

func (r *Room) handleBrokerMessage(msg *model.Message) {
	r.messages <- msg
}

func (r *Room) readMessages(ID string, conn net.Conn) {
	var (
		message string
		msg     *model.Message
	)
	input := bufio.NewScanner(conn)
	for input.Scan() {
		message = input.Text()
		if len(message) > MAX_MESSAGE_LENGTH {
			fmt.Fprintf(conn, "message is too long! (%d / %d maximum allowed characters)\n",
				len(message),
				MAX_MESSAGE_LENGTH)
			continue
		}

		if message == "/q" {
			break
		}
		msg = &model.Message{ID: ID, Text: ID + ": " + message}
		r.broker.Publish(r.name, msg)
	}
}

func (r *Room) HandleClient(ID string, conn net.Conn) {
	defer func() {
		log.Printf("closing connection from %v", conn.RemoteAddr())
		conn.Close()
		r.removeClient(ID)
	}()

	r.addClient(ID, conn)

	r.readMessages(ID, conn)
}

func (r *Room) Broadcast() {
	subscription, err := r.broker.Subscribe(r.name, r.handleBrokerMessage)
	if err != nil {
		log.Println("Error subscribing:", err)
		return
	}

	defer r.broker.Unsubscribe(subscription)

	for msg := range r.messages {
		r.mu.RLock()
		for ID, conn := range r.clients {
			if msg.ID != ID {
				_, err := fmt.Fprintln(conn, msg.Text)
				if err != nil {
					log.Println("Error sending message:", err)
				}
			}
		}
		r.mu.RUnlock()
	}
}
