package mqhandler

import (
	"encoding/json"
)

// HandlerMessage mq message
type HandlerMessage struct {
	UUID string `json:"uuid"`
}

// HandlerTransport mq transport
type HandlerTransport interface {
	Decode(data []byte) (UUID string, err error)
	Encode(uuid string) (message []byte, err error)
}

type handlerTransport struct{}

// Decode message from mq
func (t *handlerTransport) Decode(data []byte) (UUID string, err error) {
	var req HandlerMessage
	err = json.Unmarshal(data, &req)
	return req.UUID, err
}

// Encode message to mq
func (t *handlerTransport) Encode(uuid string) (message []byte, err error) {
	req := HandlerMessage{
		UUID: uuid,
	}
	return json.Marshal(req)
}

// NewHandlerTransport ...
func NewHandlerTransport() HandlerTransport {
	return &handlerTransport{}
}
