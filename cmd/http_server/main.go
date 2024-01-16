package main

import (
	"log"

	handler "github.com/plankiton/tcptunnel/api"
)

func main() {
	server := handler.NewHandler()
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("server dies: " + err.Error())
	}
}
