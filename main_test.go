package DingTalk_go

import (
	"flag"
	"fmt"
	"github.com/open-dingtalk/dingtalk-stream-sdk-go/logger"
	"testing"
)

var clientID = flag.String("c", "", "ClientID")
var token = flag.String("t", "", "Token")

func Test_main(t *testing.T) {
	session := New(*clientID, *token)
	session.AddEventHandler(func(s *Session, data *GroupJoinedEvent) {
		logger.GetLogger().Infof("GroupJoinedEvent: %s", data)
	})
	err := session.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	select {}
}
