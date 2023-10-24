package DingTalk_go

import (
	"flag"
	"fmt"
	"testing"
)

var clientID = flag.String("c", "", "ClientID")
var token = flag.String("t", "", "Token")

func Test_main(t *testing.T) {
	session := New(*clientID, *token)
	err := session.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	select {}
}
