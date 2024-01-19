package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var wsSet websocket.Upgrader

func Stream(w http.ResponseWriter, r *http.Request) {
	webSocketConn, err := wsSet.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	conn := webSocketConn.NetConn()
	handleWebsocketConnection(&conn)
	waitWebsocketConnectionFinish(&conn)
}
