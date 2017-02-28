package main

import (
	"log"
	"fmt"
	"golang.org/x/net/websocket2"
	"bufio"
	"os"
	"strings"
	"net/url"
	"flag"
	"encoding/json"
)

type Message struct {
	Client      string `json:"client"`
	ContentCode int    `json:"content_code"`
	Content     string `json:"content"`
}

var addr = flag.String("addr", "127.0.0.1:8080", "http service address")
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

	flag.Parse()
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/chatroom"}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {

		log.Fatal("dial:", err)

	}
	defer c.Close()
	//reacieve messages
	go RecieveMessage(c)
	//send messages
	WriteMessage(c)
}
func RecieveMessage(ws *websocket.Conn) {
	for {
		var msg Message
		_, message, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}
		log.Printf(string(message))
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Printf("error: %v", err)
			break
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
				err = websocket.WriteJSON(ws, msg)
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
		err = websocket.WriteJSON(ws, msg)
		if (err != nil) {
			log.Println(err)
		}
	}
}

