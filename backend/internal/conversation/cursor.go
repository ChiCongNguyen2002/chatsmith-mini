package conversation

import (
	"encoding/base64"
	"strconv"
	"strings"
	"time"

	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/uid"
)

type cursor struct {
	updatedAt time.Time
	id        string
}

func encodeCursor(item Conversation) string {
	raw := strconv.FormatInt(item.UpdatedAt.UnixMicro(), 10) + cursorSeparator + item.ID

	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

func decodeCursor(value string) (*cursor, error) {
	raw, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return nil, failure.New(failure.CodeBadRequest)
	}

	micro, id, found := strings.Cut(string(raw), cursorSeparator)
	if !found || !uid.Valid(id) {
		return nil, failure.New(failure.CodeBadRequest)
	}

	unixMicro, err := strconv.ParseInt(micro, 10, 64)
	if err != nil {
		return nil, failure.New(failure.CodeBadRequest)
	}

	return &cursor{updatedAt: time.UnixMicro(unixMicro).UTC(), id: id}, nil
}
