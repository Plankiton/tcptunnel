package handler

import (
	"net/http"

	"github.com/plankiton/tcptunnel/pkg/api"
)

var (
	messagePath = "/api/v1/message"
	streamPath  = "/api/v1/stream"
)

// V1 is the vercel Handler entrypoint
//
// path: /api/v1
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
func V1(w http.ResponseWriter, r *http.Request) {
	if isSamePath(r, streamPath) {
		api.Stream(w, r)
		return
	}

	api.Message(w, r)
}

func isSamePath(r *http.Request, path string) bool {
	return r.URL.Path[:len(path)] == path
}
