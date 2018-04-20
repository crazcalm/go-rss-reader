package database

import (
	//"strings"
	"testing"

	//"github.com/mmcdole/gofeed"

	"github.com/crazcalm/go-rss-reader/file"
)

/*
This function was changed so that it gets data from the database

func TestGuiData(t *testing.T) {
	var noData gofeed.Feed

	feed1 := &Feed{-1, URLs[snXML], Tags[snXML], getTestFeed(snXML)}
	feed2 := &Feed{-1, URLs[alertXML], Tags[alertXML], getTestFeed(alertXML)}
	feed3 := &Feed{-1, URLs[xkcdXML], Tags[xkcdXML], getTestFeed(xkcdXML)}
	feed4 := &Feed{-1, "NO URL", []string{}, &noData}

	tests := []struct {
		Feeds          Feeds
		ExpectedCount  int
		ExpectedEpData []string
	}{
		{Feeds{}, 0, []string{}},
		{Feeds{feed1}, 1, []string{"(10/10)"}},
		{Feeds{feed1, feed2, feed3}, 3, []string{"(10/10)", "(10/10)", "(4/4)"}},
		{Feeds{feed1, feed2, feed3, feed4}, 4, []string{"(10/10)", "(10/10)", "(4/4)", "(0/0)"}},
	}

	for _, test := range tests {
		data := test.Feeds.GuiData()

		if len(data) != test.ExpectedCount {
			t.Errorf("Expected %d feeds, but got %d", len(data), test.ExpectedCount)
		}

		for index, item := range data {
			if !strings.EqualFold(item.Episodes, test.ExpectedEpData[index]) {
				t.Errorf("Expected %s to give %s, but got %s", item.Title, test.ExpectedEpData[index], item.Episodes)
			}
		}
	}
}
*/

func TestNewFeeds(t *testing.T) {
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
