package garmin

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func ParseLivetrackLinkFromEmail(htmlEmailBody string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlEmailBody))
	if err != nil {
		return "", err
	}

	for n := range doc.Descendants() {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" && strings.Contains(a.Val, "https://livetrack.garmin.com") {
					return a.Val, nil
				}
			}
		}
	}

	return "", errors.New("No links found in document")
}

func IsLivetrackLinkActive(url string) bool {
	res, err := http.Get(url)
	if err != nil {
		slog.Error("Unable to fetch livetrack url", "error", err.Error())
		return false
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return false
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		slog.Error("Unable to parse response", "error", err.Error())
	}

	for n := range doc.Descendants() {
		if n.Type == html.TextNode {
			if n.Data == "The LiveTrack session you are looking for has ended." {
				return false
			}
		}
	}
	return true
}
