package pagination

import (
	"encoding/base64"
	"strconv"
	"strings"

	"chatsmith-mini/pkg/failure"
	"chatsmith-mini/pkg/uid"
)

const separator = "_"

type Keyset struct {
	Micro int64
	ID    string
}

func Encode(micro int64, id string) string {
	raw := strconv.FormatInt(micro, 10) + separator + id

	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

func Decode(value string) (Keyset, error) {
	raw, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return Keyset{}, failure.New(failure.CodeBadRequest)
	}

	micro, id, found := strings.Cut(string(raw), separator)
	if !found || !uid.Valid(id) {
		return Keyset{}, failure.New(failure.CodeBadRequest)
	}

	unixMicro, err := strconv.ParseInt(micro, 10, 64)
	if err != nil {
		return Keyset{}, failure.New(failure.CodeBadRequest)
	}

	return Keyset{Micro: unixMicro, ID: id}, nil
}
