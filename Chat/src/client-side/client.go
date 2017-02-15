package main

import (
	"log"
	"fmt"
	"golang.org/x/net/websocket"
	"bufio"
	"os"
)

var origin = "http://localhost/"
var url = "ws://localhost:8080/chatroom"

type Message struct {
	Client      string `json:"client"`
	ContentCode string `json:"content_code"`
	Content     string `json:"content"`
}

var nick string

func (msg *Message)ToString() string {
	return msg.Client + " 'says " + msg.Content
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	var input string
	fmt.Println("Press 1 to connect chatroom")
	fmt.Scanln(&input)
	if (input == "1") {
		fmt.Print("Write a nickname : ")
		fmt.Scanln(&input)
		if (len(input) > 2) {
			nick = input
			LoginChatRoom()
		} else {
			return
		}

	}
}
func LoginChatRoom() {
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	//reacieve messages
	go RecieveMessage(ws)
	//send messages
	WriteMessage(ws)
}
func RecieveMessage(ws *websocket.Conn) {
	for {
		var msg Message
		err := websocket.JSON.Receive(ws, &msg)
		if (err != nil) {
			log.Println(err)
		}
		fmt.Println(msg.ToString())
	}
}
func WriteMessage(ws *websocket.Conn) {
	for {
		var err error
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')

		if err != nil {
			log.Println(err)
		}

		msg := Message{Client:nick, Content:input}
		err = websocket.JSON.Send(ws, msg)
		if (err != nil) {
			log.Println(err)
		}
	}
}

