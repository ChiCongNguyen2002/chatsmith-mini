package llm

import "encoding/json"

const (
	anthropicServerToolUse   = "server_tool_use"
	anthropicSearchResult    = "web_search_tool_result"
	anthropicSearchResultRow = "web_search_result"

	searchToolType = "web_search_20250305"
	searchToolName = "web_search"
	maxSearchUses  = 3
)

type searchResultItem struct {
	Type  string `json:"type"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

func (a *Anthropic) handleBlockStart(event anthropicEvent, onSearch func(SearchEvent) error) error {
	if onSearch == nil {
		return nil
	}

	switch event.ContentBlock.Type {
	case anthropicServerToolUse:
		if event.ContentBlock.Name == searchToolName {
			return onSearch(SearchEvent{Type: "started"})
		}
	case anthropicSearchResult:
		for _, src := range parseSearchResults(event.ContentBlock.Content) {
			if err := onSearch(SearchEvent{Type: "source", URL: src.URL, Title: src.Title}); err != nil {
				return err
			}
		}

		return onSearch(SearchEvent{Type: "completed"})
	}

	return nil
}

func parseSearchResults(raw json.RawMessage) []searchResultItem {
	if raw == nil {
		return nil
	}
	var items []searchResultItem
	_ = json.Unmarshal(raw, &items)
	var results []searchResultItem
	for _, item := range items {
		if item.Type == anthropicSearchResultRow && item.URL != "" {
			results = append(results, item)
		}
	}

	return results
}
