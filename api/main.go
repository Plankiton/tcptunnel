package handler

import (
	"net/http"

	"github.com/plankiton/tcptunnel/pkg/api"
)

func Message(w http.ResponseWriter, r *http.Request) {
	api.Message(w, r)
}

func Stream(w http.ResponseWriter, r *http.Request) {
	api.Stream(w, r)
}
