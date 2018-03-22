package rss

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

func TestHeader(t *testing.T) {
	tests := []struct {
		Feed       *Feed
		EpisodeNum int
		AnswerPath []string
	}{
		{&Feed{URLs[snXML], Tags[snXML], getTestFeed(snXML)}, 0, []string{"test_data", "episodes", "header", "sn0.txt"}},
		{&Feed{URLs[alertXML], Tags[alertXML], getTestFeed(alertXML)}, 0, []string{"test_data", "episodes", "header", "alerts0.txt"}},
		{&Feed{URLs[xkcdXML], Tags[xkcdXML], getTestFeed(xkcdXML)}, 0, []string{"test_data", "episodes", "header", "xkcd0.txt"}},
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

func TestContent(t *testing.T) {}

func TestLinks(t *testing.T) {}
