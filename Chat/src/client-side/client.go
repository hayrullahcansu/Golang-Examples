package main

import (
	"log"
	"fmt"
	"golang.org/x/net/websocket"
	"bufio"
	"os"
	"strings"
)

var origin = "http://localhost/"
var url = "ws://localhost:8080/chatroom"

type Message struct {
	Client      string `json:"client"`
	ContentCode int `json:"content_code"`
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
			ws.Close()
			return
		}

		//client got a frame from server
		//firstly, check content code
		//message code equals is 30 means
		// server request something client
		if msg.ContentCode == 30 {
			switch msg.Content {
			case "UserName":
				msg.Client = nick
				msg.ContentCode = 31
				msg.Content = "UserName"
				//answer to server for request
				//tell our username
				err = websocket.JSON.Send(ws, msg)
				if (err != nil) {
					log.Println(err)
				}
			}

		} else if msg.ContentCode >= 20 && msg.ContentCode <= 29 {
			//content code is equals 20-29 means
			// about special request's answer from server
			fmt.Println(msg.ToString())
		} else if msg.ContentCode == 1 {
			//content code is equals 1 means
			// server send frame to print to screen
			fmt.Println(msg.ToString())
		}
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
		msg := Message{Client:nick}
		//clients checks message

		input = strings.Replace(input, "\n", "", -1)
		msg.Content = input
		if input == "-help" || input == "-list" {
			msg.ContentCode = 20
		} else {
			msg.ContentCode = 1
		}
		err = websocket.JSON.Send(ws, msg)
		if (err != nil) {
			log.Println(err)
		}
	}
}

