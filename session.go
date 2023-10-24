package DingTalk_go

import (
	"context"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/client"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/event"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/payload"
	"net/http"
	"time"
)

func New(clientID string, token string) *Session {
	logger.SetLogger(logger.NewStdTestLogger())
	s := &Session{
		ClientID: clientID,
		Token:    token,
		Client:   &http.Client{Timeout: (20 * time.Second)},
	}
	cli := client.NewStreamClient(
		client.WithAppCredential(client.NewAppCredentialConfig(clientID, token)),
		client.WithUserAgent(client.NewDingtalkGoSDKUserAgent()),
	)
	cli.RegisterCallbackRouter("/v1.0/im/bot/messages/get", chatbot.NewDefaultChatBotFrameHandler(s.onSteamEventReceived).OnEventReceived)
	cli.RegisterAllEventRouter(s.OnEventReceived)
	s.StreamClient = cli
	return s
}

func (s *Session) onSteamEventReceived(_ context.Context, data *chatbot.BotCallbackDataModel) ([]byte, error) {
	logger.GetLogger().Infof("received event, conversationId=[%s] senderId=[%s] senderNick=[%s] msgType=[%s] text=[%s] webHook=[%s]",
		data.ConversationId,
		data.SenderId,
		data.SenderNick,
		data.Msgtype,
		data.Text.Content,
		data.SessionWebhook)
	s.WebHookCallbackMap[data.ConversationId] = data.SessionWebhook
	s.handle(botCallBackHandlerEventType, data)
	return nil, nil
}

func (s *Session) OnEventReceived(ctx context.Context, df *payload.DataFrame) (frameResp *payload.DataFrameResponse, err error) {
	eventHeader := event.NewEventHeaderFromDataFrame(df)

	logger.GetLogger().Infof("received event, eventId=[%s] eventBornTime=[%d] eventCorpId=[%s] eventType=[%s] eventUnifiedAppId=[%s] data=[%s]",
		eventHeader.EventId,
		eventHeader.EventBornTime,
		eventHeader.EventCorpId,
		eventHeader.EventType,
		eventHeader.EventUnifiedAppId,
		df.Data)
	switch eventHeader.EventType {
	case "im_cool_app_install":
		var joinEvent GroupJoinedEvent
		err = Unmarshal([]byte(df.Data), &joinEvent)
		if err != nil {
			return nil, err
		}
		s.handle(botJoinGroupEventType, &joinEvent)
		break
	}
	frameResp = payload.NewSuccessDataFrameResponse()
	if err = frameResp.SetJson(event.NewEventProcessResultSuccess()); err != nil {
		return nil, err
	}
	return
}
