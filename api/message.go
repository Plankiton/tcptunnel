package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/plankiton/tcptunnel/pkg/models"
)

// Message is the vercel Handler entrypoint
//
// path: /api/message
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
//		 methods: ANY
//	  body: |
//			 {
//				 "message": {
//					 "action": {
//						 "name": "",
//						 "messageID": ""
//					 }
//				 },
//			 }
//
//	content-type: plain/text
//
//	# models.Message.Action.Name:   update | update | create
//	methods:                        "PUT    | PATCH  | POST"
//
//	body: "Hello!"                  # models.Message.Message
//	headers:
//	  senderID: ""                  # models.Message.SenderID
//		extra: null                   # models.Message.Extra
func Message(w http.ResponseWriter, r *http.Request) {
	var message models.Message

	encoder := json.NewEncoder(w)
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&message); err != nil {
		encoder.Encode(responseError{
			Message: message,
			Error:   fmt.Sprint("error decoding request body:", err),
		})

		w.WriteHeader(400)
		return
	}

	responseMessage, err := sendMessage(message, conn, isConnClosed)
	if err != nil {
		encoder.Encode(responseError{
			Message: message,
			Error:   fmt.Sprint("error tunneling message request:", err),
			Extra: map[string]interface{}{
				"response": responseMessage,
			},
		})
		return
	}

	encoder.Encode(response{
		Info: "message successful sended",
		Message: models.Message{
			Message: fmt.Sprintln(message.Action.Name, "operation processed with success"),
		},
	})
}
