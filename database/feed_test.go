package database

import (
	"strings"
	"testing"
)

func TestGetFeedDataFromSite(t *testing.T) {
	tests := []string{
		"http://www.leoville.tv/podcasts/sn.xml",
		"https://www.xkcd.com/rss.xml",
		"https://latenightlinux.com/feed/mp3",
		"http://feeds.feedburner.com/InternetHistoryPodcast",
		"https://feeds.mozilla-podcasts.org/irl",
	}

	for _, test := range tests {
		body, err := GetFeedDataFromSite(test)
		if err != nil {
			t.Errorf("Error for url (%s): %s", test, err.Error())
		}

		if strings.EqualFold(body, "") {
			t.Errorf("The body was empty for url %s", test)
		}
	}
}
