package DingTalk_go

type EventHandler interface {
	Type() string

	Handle(*Session, interface{})
}

type eventHandlerInstance struct {
	eventHandler EventHandler
}
