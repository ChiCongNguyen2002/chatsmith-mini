package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"chatsmith-mini/internal/chat"

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

type eventStream struct {
	mu   sync.Mutex
	w    gin.ResponseWriter
	done chan struct{}
	exit chan struct{}
}

func newEventStream(ctx *gin.Context) *eventStream {
	ctx.Header(headerContentType, contentTypeSSE)
	ctx.Header(headerCache, valueNoCache)
	ctx.Header(headerConnection, valueKeepAlive)
	ctx.Header(headerNoBuffering, valueNo)
	ctx.Status(http.StatusOK)

	stream := &eventStream{w: ctx.Writer, done: make(chan struct{}), exit: make(chan struct{})}
	go stream.heartbeat()

	return stream
}

func (s *eventStream) emit(event chat.Event) error {
	payload, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	return s.write(fmt.Sprintf(frameEvent, event.Name, payload))
}

func (s *eventStream) heartbeat() {
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

func (s *eventStream) stop() {
	close(s.done)
	<-s.exit
}

func (s *eventStream) write(frame string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := fmt.Fprint(s.w, frame); err != nil {
		return err
	}
	s.w.Flush()

	return nil
}
