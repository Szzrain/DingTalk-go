package DingTalk_go

var (
	EndPointDingTalk = "https://api.dingtalk.com"
	EndPointApi      = EndPointDingTalk + "/v1.0/robot/"

	EndPointAccessToken = EndPointApi + "oauth2/accessToken"

	EndPointBatchSend   = EndPointApi + "oToMessages/batchSend"
	EndPointGroupSend   = EndPointApi + "groupMessages/send"
	EndPointPrivateSend = EndPointApi + "privateChatMessages/send"
)
