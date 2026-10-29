package sessionstate

import (
	"crypto/rand"
	"encoding/base64"
)

const randomBytes = 32

type Tokens struct{}

func NewTokens() *Tokens {
	return &Tokens{}
}

func (t *Tokens) NewToken() (string, error) {
	buf := make([]byte, randomBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}
