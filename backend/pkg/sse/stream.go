package sse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	heartbeatInterval = 15 * time.Second
	headerCache       = "Cache-Control"
	headerConnection  = "Connection"
	headerNoBuffering = "X-Accel-Buffering"
	headerContentType = "Content-Type"
	contentTypeSSE    = "text/event-stream"
	valueNoCache      = "no-cache"
	valueKeepAlive    = "keep-alive"
	valueNo           = "no"
	frameEvent        = "event: %s\ndata: %s\n\n"
	frameHeartbeat    = ": ping\n\n"
)

type Stream struct {
	mu   sync.Mutex
	w    gin.ResponseWriter
	done chan struct{}
	exit chan struct{}
}

func Open(ctx *gin.Context) *Stream {
	ctx.Header(headerContentType, contentTypeSSE)
	ctx.Header(headerCache, valueNoCache)
	ctx.Header(headerConnection, valueKeepAlive)
	ctx.Header(headerNoBuffering, valueNo)
	ctx.Status(http.StatusOK)

	stream := &Stream{w: ctx.Writer, done: make(chan struct{}), exit: make(chan struct{})}
	go stream.heartbeat()

	return stream
}

func (s *Stream) Send(name string, payload any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return s.write(fmt.Sprintf(frameEvent, name, data))
}

func (s *Stream) Close() {
	close(s.done)
	<-s.exit
}

func (s *Stream) heartbeat() {
	defer close(s.exit)
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-s.done:
			return
		case <-ticker.C:
			if s.write(frameHeartbeat) != nil {
				return
			}
		}
	}
}

func (s *Stream) write(frame string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := fmt.Fprint(s.w, frame); err != nil {
		return err
	}
	s.w.Flush()

	return nil
}
