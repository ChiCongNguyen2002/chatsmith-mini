package llm

import (
	"bufio"
	"io"
	"strings"
)

const (
	maxSSELine    = 1 << 20
	sseBufferSize = 64 << 10
	sseFieldEvent = "event"
	sseFieldData  = "data"
	sseComment    = ":"
	sseSeparator  = ":"
)

type sseHandler func(event, data string) error

type sseFrame struct {
	event string
	data  strings.Builder
}

func (f *sseFrame) setField(field, value string) {
	switch field {
	case sseFieldEvent:
		f.event = value
	case sseFieldData:
		if f.data.Len() > 0 {
			f.data.WriteByte('\n')
		}
		f.data.WriteString(value)
	}
}

func (f *sseFrame) flush(handle sseHandler) error {
	if f.data.Len() == 0 {
		f.event = ""
		return nil
	}

	err := handle(f.event, f.data.String())
	f.event = ""
	f.data.Reset()

	return err
}

func readSSE(r io.Reader, handle sseHandler) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, sseBufferSize), maxSSELine)

	var frame sseFrame
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if err := frame.flush(handle); err != nil {
				return err
			}
			continue
		}
		if strings.HasPrefix(line, sseComment) {
			continue
		}

		field, value, _ := strings.Cut(line, sseSeparator)
		frame.setField(field, strings.TrimPrefix(value, " "))
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return frame.flush(handle)
}
