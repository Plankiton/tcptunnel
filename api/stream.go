package handler

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var wsSet websocket.Upgrader

func Stream(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Origin") != "http://"+r.Host {
		http.Error(w, "Origin not allowed", 403)
		return
	}

	webSocketConn, err := wsSet.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	handleWebsocketConnection(webSocketConn.NetConn())
	waitWebsocketConnectionFinish(webSocketConn.NetConn())
}
