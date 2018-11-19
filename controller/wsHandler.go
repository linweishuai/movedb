package controller

import (
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"movedb/ws"
)

func WsPage(res http.ResponseWriter, req *http.Request) {
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		http.NotFound(res, req)
		return
	}
	uuid,error:=uuid.NewV4()
	client := &ws.Client{Id: uuid.String(), Socket: conn, Send: make(chan []byte)}

	ws.Manager.Register <- client

	//go client.read()
	go client.Write()
}
