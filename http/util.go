package http

import (
	"strings"
)

const (
	defaultHTTPScheme  = "http://"
	defaultHTTPSScheme = "https://"
)

func PrepareURL(url string) string {
	if strings.HasPrefix(url, defaultHTTPScheme) || strings.HasPrefix(url, defaultHTTPSScheme) {
		return url
	}

	return defaultHTTPScheme + url
}
