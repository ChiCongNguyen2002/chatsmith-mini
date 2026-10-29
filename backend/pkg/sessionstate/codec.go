package sessionstate

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"
	"time"
)

const (
	tokenSep = "."
	fieldSep = "|"
)

var errInvalidToken = errors.New("invalid token")

type Codec struct {
	key []byte
	now func() time.Time
}

func NewCodec(secret string) *Codec {
	return &Codec{key: []byte(secret), now: time.Now}
}

func (c *Codec) Seal(value string, ttl time.Duration) string {
	expires := c.now().Add(ttl).Unix()
	payload := base64.RawURLEncoding.EncodeToString([]byte(value + fieldSep + strconv.FormatInt(expires, 10)))

	return payload + tokenSep + c.sign(payload)
}

func (c *Codec) Open(token string) (string, error) {
	payload, signature, found := strings.Cut(token, tokenSep)
	if !found || !hmac.Equal([]byte(signature), []byte(c.sign(payload))) {
		return "", errInvalidToken
	}

	raw, err := base64.RawURLEncoding.DecodeString(payload)
	if err != nil {
		return "", errInvalidToken
	}
	idx := strings.LastIndex(string(raw), fieldSep)
	if idx < 0 {
		return "", errInvalidToken
	}
	expires, err := strconv.ParseInt(string(raw[idx+len(fieldSep):]), 10, 64)
	if err != nil || c.now().Unix() >= expires {
		return "", errInvalidToken
	}

	return string(raw[:idx]), nil
}

func (c *Codec) sign(payload string) string {
	mac := hmac.New(sha256.New, c.key)
	mac.Write([]byte(payload))

	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
