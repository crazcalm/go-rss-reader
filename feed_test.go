package rss

import (
	"strings"
	"testing"
)

func checkForInternet() bool {
	/*
		Use: https://stackoverflow.com/questions/10040954/alternative-to-google-finance-api

		That post has a link to a stock market api call that returns json. Given
		that it is stock market data it will, hopefully, always be online and reachable
		from all countries.

		I can check the response from that call and decided from that whether or not
		I have internet.
	*/
	//TODO: implement later
	return true
}

func TestNewFeeds(t *testing.T) {
	//FileData variables
	goodData := FileData{"http://www.leoville.tv/podcasts/sn.xml", []string{"security", "favorite", "audio"}}
	badData := FileData{"http://www.linux-magazine.com/rs/feed/lmi_full", []string{"error"}}
	noData := FileData{"", []string{}}

	tests := []struct {
		Data        []FileData
		ExpectedNum int
	}{
		{[]FileData{}, 0},
		{[]FileData{goodData, goodData, goodData}, 3},
		{[]FileData{goodData, badData, goodData}, 2},
		{[]FileData{goodData, badData, noData}, 1},
	}

	for _, test := range tests {
		feeds := NewFeeds(test.Data)

		if len(feeds) != test.ExpectedNum {
			t.Errorf("Expected %d feed(s), but got %d feed(s)", test.ExpectedNum, len(feeds))
		}
	}
}

func TestNewFeed(t *testing.T) {
	//FileData variables
	goodData := FileData{"http://www.leoville.tv/podcasts/sn.xml", []string{"security", "favorite", "audio"}}
	badData := FileData{"http://www.linux-magazine.com/rs/feed/lmi_full", []string{"error"}}
	noData := FileData{"", []string{}}

	//Error
	errorStr := "Error occured while trying to parse feed"

	tests := []struct {
		Data      FileData
		ExpectErr bool
		Err       string
		Title     string
	}{
		{goodData, false, "", "Security Now (MP3)"},
		{badData, true, errorStr, ""},
		{noData, true, errorStr, ""},
	}

	for _, test := range tests {
		feed, err := NewFeed(test.Data)

		//Case: Received and error while expecting an error
		if err != nil && test.ExpectErr {
			if !strings.EqualFold(err.Error(), test.Err) {
				t.Errorf("Expected error: %s, received err -- %s", test.Err, err.Error())
			}
			continue
		}

		//Case: Received and error while not expecting an error
		if err != nil && !test.ExpectErr {
			t.Errorf("Unexpected error: %s", err.Error())
		}

		//Case: Expected an error, but did not receive an error
		if err == nil && test.ExpectErr {
			t.Errorf("Expected error: %s, but no error was received", err.Error())
		}

		//Case: Did not expect an error, and I did not receive an error
		if err == nil && !test.ExpectErr {
			if !strings.EqualFold(feed.Data.Title, test.Title) {
				t.Errorf("Expected title: %s, but got %s", test.Title, feed.Data.Title)
			}
			continue
		}
	}
}

//func TestTotalEpisodes(t *testing.T){}
