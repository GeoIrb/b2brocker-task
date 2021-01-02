package service

import (
	"context"
	"time"
)

// Service load
type Service struct {
	sendBack func(ctx context.Context, requestUUID string) error
}

// Handler to requests
func (s *Service) Handler(ctx context.Context, requestUUID string) (err error) {
	var timeout time.Duration
	for _, s := range requestUUID {
		timeout += time.Duration(s)
	}
	time.Sleep(timeout)

	return s.sendBack(ctx, requestUUID)
}

// NewService ...
func NewService(
	sendBack func(ctx context.Context, requestUUID string) error,
) *Service {
	return &Service{
		sendBack: sendBack,
	}
}
