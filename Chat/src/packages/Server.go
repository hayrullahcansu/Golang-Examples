package packages

import "strconv"

type Server struct {
	clients    map[*Client]bool
	broadcast  chan Message
	register   chan *Client
	unregister chan *Client
}

func NewServer() *Server {
	return &Server{
		broadcast:  make(chan Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}
func (h *Server) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
func (s *Server) SpecialRequestFromClient(c *Client, m *Message) {
	//content code equals is 20-29 means
	// about special request from client
	msg := Message{Client:"Server", ContentCode:20}
	switch m.Content {
	case "-help":
		msg.Content = s.GetHelp()
		c.send <- msg
	case "-list":
		msg.Content = s.GetUserList()
		c.send <- msg
	}
}
func (s *Server) SpecialRequestFromServer(c *Client, m *Message) {
	//content code equals is 30-39 means
	// about special request from server
	if m.ContentCode == 31 {
		switch m.Content {
		case "UserName":
			//save usernick to map of clients
			c.UserName = m.Client
			msg := Message{Client:"Server", ContentCode:1}
			msg.Content = m.Client + ", Wellcome to our Chatroom. Have a nice chats!"
			c.send <- msg
		}
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
	for client := range s.clients {
		h += strconv.Itoa(i) + ". " + client.UserName + "\n"
		i++
	}
	return h
}