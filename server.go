package hx

type Message struct {
	Includes    map[string]any
	Request     string `json:"HX-Request"`
	Trigger     string `json:"HX-Trigger"`
	TriggerName string `json:"HX-Trigger-Name"`
	Target      string `json:"HX-Target"`
	CurrentURL  string `json:"HX-Current-URL"`
}

type IServer interface {
	Mount(mountpoint string)
	//Listen(event string, listener func(*IClient, *Message))
	Broadcast(event string, message []byte) error
}

type IClient interface {
	Send(message []byte) error
	SendStr(message string) error
}
