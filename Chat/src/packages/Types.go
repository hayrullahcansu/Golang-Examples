package packages

type Message struct {
	Client      string `json:"client"`
	ContentCode int `json:"content_code"`
	Content     string `json:"content"`
}

func (self *Message) String() string {
	return self.Client + " says " + self.Content
}