package DingTalk_go

import "github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"

type EventHandler interface {
	Type() string

	Handle(*Session, interface{})
}

type eventHandlerInstance struct {
	eventHandler EventHandler
}

const botCallBackHandlerEventType = "BOT_CALLBACK"
const botJoinGroupEventType = "im_cool_app_install"

type botCallbackModelHandler func(s *Session, data *chatbot.BotCallbackDataModel)

// Type returns the event type for interface{} events.
func (eh botCallbackModelHandler) Type() string {
	return botCallBackHandlerEventType
}

func (eh botCallbackModelHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*chatbot.BotCallbackDataModel); ok {
		eh(s, t)
	}
}

type groupJoinedEventHandler func(s *Session, data *GroupJoinedEvent)

func (eh groupJoinedEventHandler) Type() string {
	return botJoinGroupEventType
}

func (eh groupJoinedEventHandler) Handle(s *Session, i interface{}) {
	if t, ok := i.(*GroupJoinedEvent); ok {
		eh(s, t)
	}
}
