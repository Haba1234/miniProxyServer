package proxyserver

import (
	"context"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

var requestCnt int64

type key string

const reqID key = "request_id"

type Logger struct {
	handler http.Handler
}

func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := atomic.AddInt64(&requestCnt, 1)
	ctx := context.WithValue(r.Context(), reqID, requestID)
	rCtx := r.Clone(ctx)
	l.handler.ServeHTTP(w, rCtx)
	log.Printf("%s %s %v request_id: %d", r.Method, r.URL.Path, time.Since(start), requestID)
}

func NewLogger(handlerToWrap http.Handler) *Logger {
	return &Logger{handlerToWrap}
}
