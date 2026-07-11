package main

import (
	"math/rand/v2"
)

const shortCodeLength = 6
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateShortCode() (string, error) {
	code := make([]byte, shortCodeLength)

	for i := 0; i < shortCodeLength; i++ {
		code[i] = charset[rand.IntN(len(charset))]
	}

	return string(code), nil
}
