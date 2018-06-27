package database

import (
	//"strings"
	"testing"

	//"github.com/mmcdole/gofeed"

	"github.com/crazcalm/go-rss-reader/file"
)

func TestNewFeeds(t *testing.T) {
	//TODO: Rethink how I do tests for the feeds.go file
	//file.Data variables
	goodData := file.Data{"http://www.leoville.tv/podcasts/sn.xml", []string{"security", "favorite", "audio"}}
	//badData := file.Data{"http://www.linux-magazine.com/rs/feed/lmi_full", []string{"error"}}
	//noData := file.Data{"", []string{}}

	tests := []struct {
		Data        map[int64]file.Data
		ExpectedNum int
	}{
		//{[]file.Data{}, 0},
		//{[]file.Data{goodData, goodData, goodData}, 3},
		//{[]file.Data{goodData, badData, goodData}, 2},
		//{[]file.Data{goodData, badData, noData}, 1},
		{map[int64]file.Data{-1: goodData}, 1},
	}

	for _, test := range tests {
		feeds := NewFeeds(test.Data)

		if len(feeds) != test.ExpectedNum {
			t.Errorf("Expected %d feed(s), but got %d feed(s)", test.ExpectedNum, len(feeds))
		}
	}
}
