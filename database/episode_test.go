package database

import (
	"log"
	"strings"
	"testing"
)

//getTestAnswer -- Used to fetch locally stored expected results
func getTestAnswer(path []string) string {
	file := getTestFile(path)
	return string(file)
}

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

func TestHeader(t *testing.T) {
	tests := []struct {
		Feed       *Feed
		EpisodeNum int
		AnswerPath []string
	}{
		{&Feed{1, URLs[snXML], "Security Now (MP3)", Tags[snXML], getTestFeed(snXML)}, 0, []string{"test_data", "episodes", "header", "sn0.txt"}},
		{&Feed{2, URLs[alertXML], "US-CERT Alerts", Tags[alertXML], getTestFeed(alertXML)}, 0, []string{"test_data", "episodes", "header", "alerts0.txt"}},
		{&Feed{3, URLs[xkcdXML], "xkcd.com", Tags[xkcdXML], getTestFeed(xkcdXML)}, 0, []string{"test_data", "episodes", "header", "xkcd0.txt"}},
	}

	for _, test := range tests {
		episode, err := test.Feed.GetEpisode(test.EpisodeNum)
		if err != nil {
			log.Fatalf("Human error that should not happen: %s", err.Error())
		}
		header := episode.Header()
		expected := string(getTestFile(test.AnswerPath))

		if !strings.EqualFold(header, expected) {
			t.Errorf("Expected: \n%s\n\nGot:\n%s\n", expected, header)
		}
	}

}

func TestContent(t *testing.T) {
	tests := []struct {
		Feed        *Feed
		EpisodeNum  int
		AnswerPath  []string
		ExpectedErr bool
	}{
		{&Feed{1, URLs[snXML], "Security Now (MP3)", Tags[snXML], getTestFeed(snXML)}, 0, []string{"test_data", "episodes", "content", "sn0.txt"}, false},
		{&Feed{2, URLs[xkcdXML], "xkcd.com", Tags[xkcdXML], getTestFeed(xkcdXML)}, 0, []string{"test_data", "episodes", "content", "xkcd0.txt"}, false},
	}

	for _, test := range tests {
		episode, err := test.Feed.GetEpisode(test.EpisodeNum)
		if err != nil {
			log.Fatalf("Human error that should not happen: %s", err.Error())
		}
		content, _, err := episode.Content()
		expected := string(getTestFile(test.AnswerPath))

		if test.ExpectedErr && err == nil {
			t.Errorf("Expected an error, but none was received")
		}

		if !test.ExpectedErr && err != nil {
			t.Errorf("Received an unexpected error: %s", err.Error())
		}

		if !strings.EqualFold(content, expected) {
			t.Errorf("Expected: \n%s\n\nGot:\n%s\n", expected, content)
		}
	}
}
