package chat

import (
	"strings"
	"unicode/utf8"

	"chatsmith-mini/internal/conversation"
	"chatsmith-mini/internal/provider"
)

func estimateTokens(text string) int64 {
	return int64((utf8.RuneCountInString(text) + charsPerToken - 1) / charsPerToken)
}

func promptTokens(messages []provider.Message) int64 {
	var total int64
	for _, message := range messages {
		total += estimateTokens(message.Content) + tokenOverhead
	}

	return total
}

func buildPrompt(history []conversation.Message, budget int) []provider.Message {
	start := len(history)
	var total int64
	for i := len(history) - 1; i >= 0; i-- {
		tokens := estimateTokens(history[i].Content) + tokenOverhead
		if total+tokens > int64(budget) && start < len(history) {
			break
		}
		total += tokens
		start = i
	}

	kept := history[start:]
	for len(kept) > 1 && kept[0].Role != conversation.RoleUser {
		kept = kept[1:]
	}

	messages := make([]provider.Message, 0, len(kept))
	for _, item := range kept {
		last := len(messages) - 1
		if last >= 0 && messages[last].Role == string(item.Role) {
			messages[last].Content += "\n\n" + item.Content
			continue
		}
		messages = append(messages, provider.Message{Role: string(item.Role), Content: item.Content})
	}

	return messages
}

func makeTitle(content string) string {
	flat := strings.Join(strings.Fields(content), " ")
	runes := []rune(flat)
	if len(runes) > titleRunes {
		return string(runes[:titleRunes])
	}

	return flat
}

func cost(usage usageData, inputMicro, outputMicro int64) int64 {
	total := usage.InputTokens*inputMicro + usage.OutputTokens*outputMicro

	return (total + microPerMTok - 1) / microPerMTok
}
