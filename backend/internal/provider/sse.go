package provider

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

func readSSE(r io.Reader, handle sseHandler) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, sseBufferSize), maxSSELine)

	var (
		event string
		data  strings.Builder
	)
	flush := func() error {
		if data.Len() == 0 {
			event = ""
			return nil
		}

		err := handle(event, data.String())
		event = ""
		data.Reset()

		return err
	}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			if err := flush(); err != nil {
				return err
			}
			continue
		}
		if strings.HasPrefix(line, sseComment) {
			continue
		}

		field, value, _ := strings.Cut(line, sseSeparator)
		value = strings.TrimPrefix(value, " ")
		switch field {
		case sseFieldEvent:
			event = value
		case sseFieldData:
			if data.Len() > 0 {
				data.WriteByte('\n')
			}
			data.WriteString(value)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return flush()
}
