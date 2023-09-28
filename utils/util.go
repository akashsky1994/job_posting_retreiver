package utils

import (
	"fmt"
	"job_posting_retreiver/constant"
	"net/url"
	"regexp"
)

func CleanURL(inURL string) string {
	if matched, _ := regexp.MatchString(`^https://jobs.lever.co/|https://boards.greenhouse.io/|https://jobs.ashbyhq.com/`, inURL); matched {
		return StripQueryParam(inURL, []string{})
	}
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
	if len(keys) == 0 {
		u.RawQuery = ""
	} else {
		u.RawQuery = q.Encode()
	}

	return u.String()
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
