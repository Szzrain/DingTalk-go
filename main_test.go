package DingTalk_go

import (
	"flag"
	"fmt"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/chatbot"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
	"testing"
)

var clientID = flag.String("c", "", "ClientID")
var token = flag.String("t", "", "Token")

func Test_main(t *testing.T) {
	logger.SetLogger(logger.NewStdTestLogger())
	session := New(*clientID, *token)
	session.AddEventHandler(func(s *Session, data *chatbot.BotCallbackDataModel) {
		logger.GetLogger().Infof("BotCallbackDataModel: %s", data)

	})
	session.AddEventHandler(func(s *Session, data *GroupJoinedEvent) {
		logger.GetLogger().Infof("GroupJoinedEvent: %s", data)
		_, err := s.MessageGroupSend(data.OpenConversationId, data.RobotCode, data.CoolAppCode, &MessageSampleText{Content: "欢迎加入"})
		if err != nil {
			logger.GetLogger().Errorf("MessageGroupSend failed, err=[%s]", err.Error())
			return
		}
	})
	err := session.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	select {}
}
