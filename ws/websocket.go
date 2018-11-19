package ws

import (
	"github.com/gorilla/websocket"
	"encoding/json"
)

type ClientManager struct {
	clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
}

type Client struct {
	Id     string
	Socket *websocket.Conn
	Send   chan []byte
}

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var Manager = ClientManager{
	Broadcast:  make(chan []byte),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

func (Manager *ClientManager) Start() {
	for {
		select {
		case conn := <-Manager.Register:
			Manager.clients[conn] = true
		case conn := <-Manager.Unregister:
			if _, ok := Manager.clients[conn]; ok {
				close(conn.Send)
				delete(Manager.clients, conn)
			}
		case message := <-Manager.Broadcast:
			for conn := range Manager.clients {
				conn.Send <- message
			}
		}
	}
}

func (Manager *ClientManager) Send(message []byte, ignore *Client) {
	for conn := range Manager.clients {
		if conn != ignore {
			conn.Send <- message
		}
	}
}

//func (c *Client) read() {
//	defer func() {
//		Manager.unregister <- c
//		c.socket.Close()
//	}()
//
//	for {
//		_, message, err := c.socket.ReadMessage()
//		if err != nil {
//			Manager.unregister <- c
//			c.socket.Close()
//			break
//		}
//		jsonMessage, _ := json.Marshal(&Message{Content: string(message)})
//		Manager.Broadcast <- jsonMessage
//	}
//}

func (c *Client) Write() {
	defer func() {
		c.Socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}
func ReceiveApplication(msg string)  {
	jsonMessage, _ := json.Marshal(&Message{Content: msg})
	//fmt.Println(jsonMessage)
	Manager.Broadcast <- jsonMessage
}