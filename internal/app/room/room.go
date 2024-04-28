package room

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/pu4mane/goChat/internal/app/model"
)

const MAX_MESSAGE_LENGTH = 1000

type ChatRoom struct {
	clients  sync.Map
	messages chan model.Message
}

func NewChatRoom() *ChatRoom {
	return &ChatRoom{
		messages: make(chan model.Message),
	}
}

func (cs *ChatRoom) AddClient(id string, conn net.Conn) {
	cs.clients.Store(id, conn)
	log.Println("User", id, "connected to the chat")
}

func (cs *ChatRoom) RemoveClient(id string) {
	cs.clients.Delete(id)
	log.Println("User", id, "left the chat")
}

func (cs *ChatRoom) HandleClient(id string, conn net.Conn) {
	defer conn.Close()

	cs.AddClient(id, conn)

	input := bufio.NewScanner(conn)
	for input.Scan() {
		message := input.Text()
		//добавил ограничение символов
		if len(message) > MAX_MESSAGE_LENGTH {
			fmt.Fprintf(conn, "message is too long! (%d / %d maximum allowed characters)\n",
				len(message),
				MAX_MESSAGE_LENGTH)
			continue
		}
		//добавил выход
		if message == "/q" {
			break
		}

		cs.messages <- model.Message{Text: id + ": " + message, ID: id}
	}

	cs.RemoveClient(id)
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
