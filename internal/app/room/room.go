package room

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/pu4mane/goChat/internal/app/model"
)

type ChatRoom struct {
	clients  sync.Map
	messages chan model.Message
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		messages: make(chan model.Message),
	}
}

func (cs *ChatRoom) HandleClient(id string, conn net.Conn) {
	cs.clients.Store(id, conn)
	log.Println("User", id, "connected to the chat")

	input := bufio.NewScanner(conn)
	for input.Scan() {
		cs.messages <- model.Message{Text: id + ": " + input.Text(), ID: id}
	}

	cs.clients.Delete(id)
	log.Println("User", id, "left the chat")
	conn.Close()
}

func (cs *ChatRoom) Broadcaster() {
	for {
		msg := <-cs.messages
		cs.clients.Range(func(key, value interface{}) bool {
			clientID := key.(string)
			conn := value.(net.Conn)
			if msg.ID != clientID {
				_, err := fmt.Fprintln(conn, msg.Text)
				if err != nil {
					log.Println("Error sending message:", err)
				}
			}
			return true
		})
	}
}
