package main

import (
	"./packages"
	"log"
	"net/http"
)

func main() {
	server := packages.NewServer("/entry")
	go server.WorkToListen()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
