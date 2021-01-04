package httprouter

import (
	"encoding/json"
	"net/http"

	"github.com/valyala/fasthttp"
)

// HandlerResponse response to http
type HandlerResponse struct {
	UUID string `json:"uuid"`
}

// HandlerTransport http interface transport
type HandlerTransport interface {
	Encode(response *fasthttp.Response, uuid string) (err error)
}

type handlerTransport struct{}

// Encode message to mq
func (t *handlerTransport) Encode(response *fasthttp.Response, uuid string) (err error) {
	req := HandlerResponse{
		UUID: uuid,
	}
	body, err := json.Marshal(req)
	if err != nil {
		return
	}
	response.SetBody(body)
	response.SetStatusCode(http.StatusOK)
	return
}

// NewHandlerTransport ...
func NewHandlerTransport() HandlerTransport {
	return &handlerTransport{}
}
