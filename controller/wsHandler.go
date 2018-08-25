package controller

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
)

type ClientManager struct {
	clients    map[*Client]bool
	Broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	id     string
	socket *websocket.Conn
	send   chan []byte
}

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var Manager = ClientManager{
	Broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

func (Manager *ClientManager) Start() {
	for {
		select {
		case conn := <-Manager.register:
			Manager.clients[conn] = true
		case conn := <-Manager.unregister:
			if _, ok := Manager.clients[conn]; ok {
				close(conn.send)
				delete(Manager.clients, conn)
			}
		case message := <-Manager.Broadcast:
			for conn := range Manager.clients {
				conn.send <- message
			}
		}
	}
}

func (Manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range Manager.clients {
		if conn != ignore {
			conn.send <- message
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

func (c *Client) write() {
	defer func() {
		c.socket.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func WsPage(res http.ResponseWriter, req *http.Request) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}
	uuid,error:=uuid.NewV4()
	client := &Client{id: uuid.String(), socket: conn, send: make(chan []byte)}

	Manager.register <- client

	//go client.read()
	go client.write()
}
func ReceiveApplication(msg string)  {
	jsonMessage, _ := json.Marshal(&Message{Content: msg})
	//fmt.Println(jsonMessage)
	Manager.Broadcast <- jsonMessage
}
