package utils

import (
	"net/url"
	"strconv"
)

func FormQueryFromURL(url url.URL) (desc bool, since string, limit int, flat string) {
	queryURL := url.Query()
	descString := queryURL.Get("desc")
	if descString == "true" {
		desc = true
	}
	since = queryURL.Get("since")
	limit, err := strconv.Atoi(queryURL.Get("limit"))
	if err != nil || limit < 1 {
		limit = 100
	}
	flat = queryURL.Get("sort")
	return desc, since, limit, flat
}
