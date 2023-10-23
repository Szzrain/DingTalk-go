package DingTalk_go

import (
	"context"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
)

func (s *Session) addEventHandler(eventHandler EventHandler) func() {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()

	if s.handlers == nil {
		s.handlers = map[string][]*eventHandlerInstance{}
	}

	ehi := &eventHandlerInstance{eventHandler}
	s.handlers[eventHandler.Type()] = append(s.handlers[eventHandler.Type()], ehi)

	return func() {
		s.removeEventHandlerInstance("event", ehi)
	}
}

func (s *Session) AddEventHandler(handler interface{}) func() {
	eh := handlerForInterface(handler)
	return s.addEventHandler(eh)
}

func handlerForInterface(handler interface{}) EventHandler {
	switch v := handler.(type) {
	case func(s *Session, data *chatbot.BotCallbackDataModel):
		return botCallbackModelHandler(v)
	}
	return nil
}

func (s *Session) removeEventHandlerInstance(t string, ehi *eventHandlerInstance) {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()

	handlers := s.handlers[t]
	for i := range handlers {
		if handlers[i] == ehi {
			s.handlers[t] = append(handlers[:i], handlers[i+1:]...)
		}
	}

	onceHandlers := s.onceHandlers[t]
	for i := range onceHandlers {
		if onceHandlers[i] == ehi {
			s.onceHandlers[t] = append(onceHandlers[:i], onceHandlers[i+1:]...)
		}
	}
}

func (s *Session) onSteamEventReceived(c context.Context, data *chatbot.BotCallbackDataModel) ([]byte, error) {
	s.handle(botCallBackHandlerEventType, data)
	return nil, nil
}

func (s *Session) handle(t string, i interface{}) {
	for _, eh := range s.handlers[t] {
		if s.SyncEvents {
			eh.eventHandler.Handle(s, i)
		} else {
			go eh.eventHandler.Handle(s, i)
		}
	}

	if len(s.onceHandlers[t]) > 0 {
		for _, eh := range s.onceHandlers[t] {
			if s.SyncEvents {
				eh.eventHandler.Handle(s, i)
			} else {
				go eh.eventHandler.Handle(s, i)
			}
		}
		s.onceHandlers[t] = nil
	}
}
