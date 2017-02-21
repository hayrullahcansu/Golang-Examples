package main

import (
	"./packages"
	"log"
	"net/http"
	"flag"
	"io/ioutil"
	"fmt"
	"os"
)

var addr = flag.String("addr", "127.0.0.1:8080", "http service address")

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	flag.Parse()
	server := packages.NewServer()
	go server.Run()
	http.HandleFunc("/chatroom", func(w http.ResponseWriter, r *http.Request) {

		htmlData, err := ioutil.ReadAll(r.Body) //<--- here!
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// print out
		fmt.Println(string(htmlData)) //<-- here !
		packages.ServeWs(server, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
