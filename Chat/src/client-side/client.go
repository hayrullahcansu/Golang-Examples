package main

import (
	"log"
	"fmt"
	"golang.org/x/net/websocket"
)

var origin = "http://localhost/"
var url = "ws://localhost:8080/chatroom"

type Message struct {
	Client  string `json:"client"`
	Content string `json:"content"`
}

var nick string
/*func main() {
	flag.Parse()
	log.SetFlags(0)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	done := make(chan struct{})
	go func() {
		defer c.Close()
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case t := <-ticker.C:
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
		// To cleanly close a connection, a client should send a close

		// frame and wait for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
				select {
				case <-done:
				case <-time.After(time.Second):
				}
			c.Close()
			return
		}
	}
}*/

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
		var input string
		fmt.Scanln(&input)
		msg := Message{Client:nick, Content:input}
		err := websocket.JSON.Send(ws, msg)
		if (err != nil) {
			log.Println(err)
		}
	}
}

