package packages

import (
	"golang.org/x/net/websocket"
	"fmt"
	"log"
	"io"
)

const FrameBufferSize = 200

var IDGenarator = 0

type Client struct {
	ID             int
	WS             *websocket.Conn
	server         *Server
	recieveMsgChan chan *Message
	doneChan       chan bool
}

func NewClient(ws *websocket.Conn, server *Server) *Client {
	if ws == nil {
		panic("websocket can not be nil")
	}
	if server == nil {
		panic("server can not be nil")
	}
	IDGenarator++
	messageChan := make(chan *Message, FrameBufferSize)
	doneCh := make(chan bool)
	return &Client{IDGenarator, ws, server, messageChan, doneCh}
}

func (c *Client) Write(msg *Message) {
	select {
	case c.recieveMsgChan <- msg:
	default:
		c.server.Del(c)
		err := fmt.Errorf("client %d is disconnected.", c.ID)
		c.server.Err(err)
	}
}

// Listen Write and Read request via channel
func (c *Client) Listen() {
	go c.RequestWrite()
	c.RequestRead()
}

func (c *Client)RequestWrite() {
	for {
		select {

		// send message to the client
		case msg := <-c.recieveMsgChan:
			log.Println("Send:", msg)
			websocket.JSON.Send(c.WS, msg)

		// receive done request
		case <-c.doneChan:
			c.server.Del(c)
			c.doneChan <- true // for listenRead method
			return
		}
	}
}

func (c *Client)RequestRead() {
	for {
		select {
		// receive done request
		case <-c.doneChan:
			c.server.Del(c)
			c.doneChan <- true // for listenWrite method
			return

		// read data from websocket client connection
		default:
			var msg Message
			err := websocket.JSON.Receive(c.WS, &msg)
			if err == io.EOF {
				c.doneChan <- true
			} else if err != nil {

				//added this code
				//when client closed the connection, server was printing error messages endlessly
				c.doneChan <- true

				c.server.Err(err)
			} else {
				c.server.SendAll(&msg)
			}
		}
	}
}
