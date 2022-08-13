package utils

import (
	"math/rand"
	"net/url"
	"strings"
	"time"
)

func payloadEncode(payload string) string {
	return url.QueryEscape(payload)
}

func generateRandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	var chars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
