package utils

import (
	"fmt"
	"job_posting_retreiver/constant"
	"net/url"
)

func CleanURL(inURL string) string {
	return StripQueryParam(inURL, constant.REDUNDANT_PARAMS)
}

func StripQueryParam(inURL string, keys []string) string {
	u, err := url.Parse(inURL)
	if err != nil {
		// TODO: log or handle error, in the meanwhile just return the original
		fmt.Println("Error URL:", err)
		return inURL
	}
	q := u.Query()
	for _, stripKey := range keys {
		q.Del(stripKey)
	}
	u.RawQuery = q.Encode()
	return u.String()
}
