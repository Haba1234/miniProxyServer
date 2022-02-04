package proxyserver

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

const defaultInterval = 3 * time.Second

type RequestResponse struct {
	Request
	QueryResult
}

type Service struct {
	Stats    map[uuid.UUID]RequestResponse
	Interval time.Duration
}

func NewService() *Service {
	return &Service{
		Stats:    make(map[uuid.UUID]RequestResponse),
		Interval: defaultInterval,
	}
}

func (s *Service) toSave(id uuid.UUID, req Request, res QueryResult) {
	s.Stats[id] = RequestResponse{
		Request:     req,
		QueryResult: res,
	}
}
