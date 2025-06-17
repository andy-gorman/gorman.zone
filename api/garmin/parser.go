package garmin

import (
	"errors"
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

	return "", errors.New("no links found in document")
}
