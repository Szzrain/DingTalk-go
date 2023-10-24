package DingTalk_go

import (
	"context"
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

	WebHookCallbackMap map[string]string

	SyncEvents bool

	handlersMu   sync.RWMutex
	handlers     map[string][]*eventHandlerInstance
	onceHandlers map[string][]*eventHandlerInstance
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

type GroupJoinedEvent struct {
	EventId                string `json:"eventId"`
	OperateTime            int64  `json:"operateTime"`
	CoolAppCode            string `json:"coolAppCode"`
	OpenConversationId     string `json:"openConversationId"`
	RobotCode              string `json:"robotCode"`
	OpenConversationCorpId string `json:"openConversationCorpId"`
	Operator               string `json:"operator"`
}
