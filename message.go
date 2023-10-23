package DingTalk_go

type MessageType string

const (
	MessageTypeText        MessageType = "sampleText"
	MessageTypeMarkdown    MessageType = "sampleMarkdown"
	MessageTypeImage       MessageType = "sampleImageMsg"
	MessageTypeLink        MessageType = "sampleLink"
	MessageTypeActionCard  MessageType = "sampleActionCard"
	MessageTypeActionCard2 MessageType = "sampleActionCard2"
	MessageTypeActionCard3 MessageType = "sampleActionCard3"
	MessageTypeActionCard4 MessageType = "sampleActionCard4"
	MessageTypeActionCard5 MessageType = "sampleActionCard5"
	MessageTypeActionCard6 MessageType = "sampleActionCard6"
)

type Message interface {
	Type() MessageType
}

type MessageSampleText struct {
	Content string `json:"content"`
}

func (msg *MessageSampleText) Type() MessageType {
	return MessageTypeText
}

type MessageSampleMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

func (msg *MessageSampleMarkdown) Type() MessageType {
	return MessageTypeMarkdown
}

type MessageSampleImage struct {
	PhotoURL string `json:"photoURL"`
}

func (msg *MessageSampleImage) Type() MessageType {
	return MessageTypeImage
}

type MessageSampleLink struct {
	Title      string `json:"title"`
	Text       string `json:"text"`
	MessageURL string `json:"messageURL"`
	PicURL     string `json:"picURL"`
}

func (msg *MessageSampleLink) Type() MessageType {
	return MessageTypeLink
}
