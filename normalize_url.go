package main 

import (
	"fmt"
	"net/url"
	"strings"
)

func normalizeURL(u string) (string, error) {
	if u[len(u)-1] == '/' {
		u = u[:len(u)-1]
	}

	parsedURL, err := url.Parse(u)
	if err != nil {
		return "", fmt.Errorf("url: %s, error: %w", u, err)
	}

	normURL := parsedURL.Host + parsedURL.Path
	normURL = strings.ToLower(normURL)

	return normURL, nil
}
