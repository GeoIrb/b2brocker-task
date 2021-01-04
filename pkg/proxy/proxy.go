package proxy

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/google/uuid"
)

// Proxy service
type Proxy struct {
	mutex   sync.Mutex
	request map[string]chan struct{}

	sendNext func(ctx context.Context, data []byte) error
}

// Handler to requests
func (p *Proxy) Handler(ctx context.Context) (requestUUID string, err error) {
	requestUUID = uuid.New().String()
	//todo
	message := struct {
		UUID string `json:"uuid"`
	}{requestUUID}

	data, err := json.Marshal(message)
	if err != nil {
		return
	}

	ch := make(chan struct{})
	p.mutex.Lock()
	p.request[requestUUID] = ch
	p.mutex.Unlock()

	if err = p.sendNext(ctx, data); err != nil {
		p.mutex.Lock()
		delete(p.request, requestUUID)
		p.mutex.Unlock()
		return
	}
	<-ch
	return
}

// Compliter complites requests
func (p *Proxy) Compliter(ctx context.Context, requestUUID string) {
	p.mutex.Lock()
	if ch, isExist := p.request[requestUUID]; isExist {
		close(ch)
		delete(p.request, requestUUID)
	}
	p.mutex.Unlock()
}

// New proxy service
func New(
	sendNext func(ctx context.Context, data []byte) error,
) *Proxy {
	return &Proxy{
		request:  make(map[string]chan struct{}),
		sendNext: sendNext,
	}
}
