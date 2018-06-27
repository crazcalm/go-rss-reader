package database

import (
	"strings"
	"testing"
)

func TestFormatEpisodeHeader(t *testing.T) {
	feedName := "feedName"
	epTitle := "Episode Title"
	author := "Jane Doe"
	date := "Jan 01, 2045"
	mediaContent := "link to media"
	epLink := "Episode link"

	tests := []struct {
		FeedName     string
		EpTitle      string
		Author       string
		EpLink       string
		Date         string
		MediaContent string
		Expect       string
	}{
		{"", "", "", "", "", "", ""},
		{feedName, epTitle, author, epLink, date, mediaContent,
			"Feed: feedName\nTitle: Episode Title\nAuthor: Jane Doe\nLink: Episode link\nDate: Jan 01, 2045\nMedia Link: link to media\n"},
	}

	for _, test := range tests {
		result := FormatEpisodeHeader(test.FeedName, test.EpTitle, test.Author, test.EpLink, test.Date, test.MediaContent)
		if !strings.EqualFold(result, test.Expect) {
			t.Errorf("Expected %s, but got %s", test.Expect, result)
		}
	}
}
