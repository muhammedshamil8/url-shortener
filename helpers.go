package main

import (
	"errors"
	"math/rand/v2"
	"net/url"
)

const shortCodeLength = 6
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const (
	HTTP  = "http"
	HTTPS = "https"
)

func generateShortCode() (string, error) {
	code := make([]byte, shortCodeLength)

	for i := 0; i < shortCodeLength; i++ {
		code[i] = charset[rand.IntN(len(charset))]
	}

	return string(code), nil
}

func validateURL(rawURL string) error {
	u,err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	if u.Scheme != HTTP && u.Scheme != HTTPS {
		return errors.New("invalid url scheme")
	}
	if u.Host == "" {
		return errors.New("invalid url host")
	}
	return nil
}
