package api

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var wsSet websocket.Upgrader

// Message is the vercel Handler entrypoint
//
// path: /api/v1/stream
// methods allowed: ANY
//
//	model: models.Message {
//		RetryCount: 0,
//		Action: models.Action{
//			Name:      "delete",        # get                    | delete   | create   | update
//			MessageID: "dlkfjsdfjdsfs", # Optional: Last Message | Required | Optional | Required
//		},
//		Message:   "Hello!",
//		SenderID:  "",                # Optional
//		Extra:     nil,               # Optional: ANY
//		SendedAt:  time.Time{},       # Optional Example: "2024-01-16 12:06:02.758586816-03:00" | RFC3339 full datetime
//		CreatedAt: time.Time{},       # Optional Example: "2024-01-16 12:06:02.758586816-03:00" | RFC3339 full datetime
//	}
//
//	content-type: json/application
//
//		methods: ANY
//	  body: |
//			 {
//				 "message": {
//					 "action": {
//						 "name": "",
//						 "messageID": ""
//					 }
//				 },
//			 }
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
