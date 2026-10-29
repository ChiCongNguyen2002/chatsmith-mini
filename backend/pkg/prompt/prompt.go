package prompt

import (
	"strings"
	"unicode/utf8"
)

const (
	RoleUser      = "user"
	charsPerToken = 4
	tokenOverhead = 4
	titleRunes    = 40
)

type Message struct {
	Role    string
	Content string
}

func EstimateTokens(text string) int64 {
	return int64((utf8.RuneCountInString(text) + charsPerToken - 1) / charsPerToken)
}

func Tokens(messages []Message) int64 {
	var total int64
	for _, message := range messages {
		total += EstimateTokens(message.Content) + tokenOverhead
	}

	return total
}

func Build(history []Message, budget int) []Message {
	start := len(history)
	var total int64
	for i := len(history) - 1; i >= 0; i-- {
		tokens := EstimateTokens(history[i].Content) + tokenOverhead
		if total+tokens > int64(budget) && start < len(history) {
			break
		}
		total += tokens
		start = i
	}

	kept := history[start:]
	for len(kept) > 1 && kept[0].Role != RoleUser {
		kept = kept[1:]
	}

	messages := make([]Message, 0, len(kept))
	for _, item := range kept {
		last := len(messages) - 1
		if last >= 0 && messages[last].Role == item.Role {
			messages[last].Content += "\n\n" + item.Content
			continue
		}
		messages = append(messages, Message{Role: item.Role, Content: item.Content})
	}

	return messages
}

func Title(content string) string {
	flat := strings.Join(strings.Fields(content), " ")
	runes := []rune(flat)
	if len(runes) > titleRunes {
		return string(runes[:titleRunes])
	}

	return flat
}
