package handler

import (
	"net/http"

	"github.com/plankiton/tcptunnel/pkg/api"
)

var (
	messagePath = "/api/v1/message"
	streamPath  = "/api/v1/stream"
)

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
