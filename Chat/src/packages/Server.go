package packages

import (
	"golang.org/x/net/websocket"
	"net/http"
	"log"
	"strconv"
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
	msg := Message{Client:"Server", ContentCode:30, Content:"UserName"}
	defer c.Write(&msg)
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

func (s *Server) GetHelp() string {
	h := "Writable special request codes are below\n\n" +
		"-help     helps\n" +
		"-list     gets online user lists\n"
	return h
}
func (s *Server) GetUserList() string {
	h := strconv.Itoa(len(s.clients)) + " Users are Online now\n\n"
	i := 1
	for _, c := range s.clients {
		h += strconv.Itoa(i) + ". " + c.UserName + "\n"
		i++
	}
	return h
}
func (s *Server) specialRequestFromClient(c *Client, m *Message) {

	//content code equals is 20-29 means
	// about special request from client
	msg := Message{Client:"Server", ContentCode:20}
	switch m.Content {
	case "-help":
		msg.Content = s.GetHelp()
		c.Write(&msg)
	case "-list":
		msg.Content = s.GetUserList()
		c.Write(&msg)
	}
}
func (s *Server) specialRequestFromServer(c *Client, m *Message) {

	//content code equals is 30-39 means
	// about special request from server
	if m.ContentCode == 31 {
		switch m.Content {
		case "UserName":
			//save usernick to map of clients
			c.UserName = m.Client
		}
	}
}

func (s *Server) WorkToListen() {
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