package DingTalk_go

import (
	"context"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
	"net/http"
	"sync"
)

type Session struct {
	sync.RWMutex

	ClientID              string
	Token                 string
	AccessToken           string
	AccessTokenTimeStamps int64

	StreamClient *client.StreamClient
	Client       *http.Client

	ShouldRetryOnRateLimit bool
	MaxRestRetries         int

	SyncEvents bool

	handlersMu   sync.RWMutex
	handlers     map[string][]*eventHandlerInstance
	onceHandlers map[string][]*eventHandlerInstance
}

func (s *Session) addEventHandler(eventHandler EventHandler) func() {
	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()

	if s.handlers == nil {
		s.handlers = map[string][]*eventHandlerInstance{}
	}

	ehi := &eventHandlerInstance{eventHandler}
	s.handlers["event"] = append(s.handlers["event"], ehi)

	return func() {
		s.removeEventHandlerInstance("event", ehi)
	}
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
	s.handle("event", data)
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

func (s *Session) Close() (err error) {
	s.Lock()
	defer s.Unlock()
	s.StreamClient.Close()
	return
}

func (s *Session) Open() (err error) {
	s.Lock()
	defer s.Unlock()
	err = s.StreamClient.Start(context.Background())
	if err != nil {
		return err
	}
	return
}
