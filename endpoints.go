package DingTalk_go

var (
	EndPointDingTalk = "https://api.dingtalk.com"
	EndPointApi      = EndPointDingTalk + "/v1.0/"

	EndPointAccessToken = EndPointApi + "oauth2/accessToken"

	EndPointRobot = EndPointApi + "robot/"

	EndPointBatchSend   = EndPointRobot + "oToMessages/batchSend"
	EndPointGroupSend   = EndPointRobot + "groupMessages/send"
	EndPointPrivateSend = EndPointRobot + "privateChatMessages/send"
)
