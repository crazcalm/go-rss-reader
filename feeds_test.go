package rss

import (
	"strings"
	"testing"

	"github.com/mmcdole/gofeed"
)

func TestGuiData(t *testing.T) {
	var noData gofeed.Feed

	feed1 := &Feed{URLs[snXML], Tags[snXML], getTestFeed(snXML)}
	feed2 := &Feed{URLs[alertXML], Tags[alertXML], getTestFeed(alertXML)}
	feed3 := &Feed{URLs[xkcdXML], Tags[xkcdXML], getTestFeed(xkcdXML)}
	feed4 := &Feed{"NO URL", []string{}, &noData}

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
