package main

import (
	"log"
	"net/http"

	"github.com/plankiton/tcptunnel/pkg/api"
)

func main() {
	http.HandleFunc("/message", api.Message)
	http.HandleFunc("/stream", api.Stream)

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("server dies: " + err.Error())
	}
}
