package DingTalk_go

import (
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/utils"
	"net/http"
	"time"
)

func New(clientID string, token string) *Session {
	s := &Session{
		ClientID: clientID,
		Token:    token,
		Client:   &http.Client{Timeout: (20 * time.Second)},
	}
	cli := client.NewStreamClient(
		client.WithAppCredential(client.NewAppCredentialConfig(clientID, token)),
		client.WithUserAgent(client.NewDingtalkGoSDKUserAgent()),
		client.WithSubscription(utils.SubscriptionTypeKCallback, "/v1.0/im/bot/messages/get", chatbot.NewDefaultChatBotFrameHandler(s.onSteamEventReceived).OnEventReceived),
	)
	s.StreamClient = cli
	return s
}
