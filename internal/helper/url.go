package helper

import (
	"net/url"
)

func IsValidURL(testURL string) bool {
	parsedURL, err := url.Parse(testURL)
	if err != nil {
		return false
	}

	// Проверка на наличие схемы и хоста
	return parsedURL.Scheme != "" && parsedURL.Host != ""
}
