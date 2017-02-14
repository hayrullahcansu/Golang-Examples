package packages

import (
	"golang.org/x/net/websocket"
	"net/http"
	"log"
)

type Server struct {
	url              string
	messages         []*Message
	clients          map[int]*Client
	addClientCh      chan *Client
	deleteClientCh   chan *Client
	sendAllClientsCh chan *Message
	doneCh           chan bool
	errorCh          chan error
}

// Create new chat server.
func NewServer(pattern string) *Server {
	messages := []*Message{}
	clients := make(map[int]*Client)
	addCh := make(chan *Client)
	delCh := make(chan *Client)
	sendAllCh := make(chan *Message)
	doneCh := make(chan bool)
	errCh := make(chan error)

	return &Server{
		pattern,
		messages,
		clients,
		addCh,
		delCh,
		sendAllCh,
		doneCh,
		errCh,
	}
}


func (s *Server) Add(c *Client) {
	s.addClientCh <- c
}

func (s *Server) Del(c *Client) {
	s.deleteClientCh <- c
}

func (s *Server) SendAll(msg *Message) {
	s.sendAllClientsCh <- msg
}

func (s *Server) Done() {
	s.doneCh <- true
}

func (s *Server) Err(err error) {
	s.errorCh <- err
}

func (s *Server) sendPastMessages(c *Client) {
	for _, msg := range s.messages {
		c.Write(msg)
	}
}

func (s *Server) sendAll(msg *Message) {
	for _, c := range s.clients {
		c.Write(msg)
	}
}

func (s *Server) WorkToListen(){
	onConnected := func(ws *websocket.Conn) {
		defer func() {
			err := ws.Close()
			if err != nil {
				s.errorCh <- err
			}
		}()

		client := NewClient(ws, s)
		s.Add(client)
		client.Listen()
	}
	http.Handle(s.url, websocket.Handler(onConnected))
	for {
		select {

		// Add new a client
		case c := <-s.addClientCh:
			log.Println("Added new client")
			s.clients[c.ID] = c
			log.Println("Now", len(s.clients), "clients connected.")
			s.sendPastMessages(c)

		// del a client
		case c := <-s.deleteClientCh:
			delete(s.clients, c.ID)

		// broadcast message for all clients
		case msg := <-s.sendAllClientsCh:
			log.Println("Send all:", msg)
			s.messages = append(s.messages, msg)
			s.sendAll(msg)

		case err := <-s.errorCh:
			log.Println("Error:", err.Error())

		case <-s.doneCh:
			return
		}
	}

}