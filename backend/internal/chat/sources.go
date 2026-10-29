package chat

import (
	"encoding/json"

	"chatsmith-mini/internal/provider"
)

type sourceLog []searchSourceData

func (l *sourceLog) handle(emit func(Event) error, se provider.SearchEvent) error {
	switch se.Type {
	case "started":
		return emit(Event{Name: EventSearchStarted})
	case "source":
		item := searchSourceData{URL: se.URL, Title: se.Title}
		*l = append(*l, item)
		return emit(Event{Name: EventSearchSource, Data: item})
	case "completed":
		return emit(Event{Name: EventSearchCompleted})
	}

	return nil
}

func marshalSources(items []searchSourceData) (any, error) {
	if len(items) == 0 {
		return nil, nil
	}
	raw, err := json.Marshal(items)
	if err != nil {
		return nil, err
	}

	return string(raw), nil
}
