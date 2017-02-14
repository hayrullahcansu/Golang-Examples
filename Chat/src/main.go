package main

import (
	"./packages"
	"log"
	"net/http"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	server := packages.NewServer("/chatroom")
	go server.WorkToListen()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
