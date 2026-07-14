package utils

import (
	"errors"
	"net/url"
)

const (
	HTTP  = "http"
	HTTPS = "https"
)

func ValidateURL(rawURL string) error {
	u, err := url.Parse(rawURL)
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
