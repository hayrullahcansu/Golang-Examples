package packages

type Message struct {
	Client string `json:"client"`
	Conent   string `json:"content"`
}

func (self *Message) String() string {
	return self.Client + " says " + self.Conent
}