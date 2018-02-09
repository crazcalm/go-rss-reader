package rss

import (
	"testing"
)

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

