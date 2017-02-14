package packages

type Message struct {
	Client  string `json:"client"`
	Content string `json:"content"`
}

func (self *Message) String() string {
	return self.Client + " says " + self.Content
}