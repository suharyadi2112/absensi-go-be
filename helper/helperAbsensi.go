package helper

import (
	"math/rand"
	"strings"
	"time"
)

func GenerateRandomString(length int) string {
	const charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var sb strings.Builder
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < length; i++ {
		randomIndex := rand.Intn(len(charset))
		sb.WriteByte(charset[randomIndex])
	}

	return sb.String()
}
