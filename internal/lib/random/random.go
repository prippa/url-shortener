package random

import (
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewRandomString(length int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	charsetLen := len(charset)
	b := make([]byte, length)

	for i := range b {
		b[i] = charset[rnd.Intn(charsetLen)]
	}

	return string(b)
}
