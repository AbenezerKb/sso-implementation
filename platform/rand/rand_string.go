package rand

import (
	"crypto/rand"
	"io"
)

const (
	SmallLetters   = "abcdefghijklmnopqrstuvwxyz"
	CapitalLetters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits         = "0987654321"
)

func GenerateCustomRandomString(set string, length int) string {
	randString := make([]byte, length)
	_, _ = io.ReadAtLeast(rand.Reader, randString, length)
	for i := 0; i < len(randString); i++ {
		randString[i] = set[int(randString[i])%len(set)]
	}

	return string(randString)
}
