package database

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/mmcdole/gofeed"

	"github.com/crazcalm/go-rss-reader/file"
)

var (
	testFeed1 = Feed{ID: 1, URL: "1.com", Title: "1", Tags: []string{}, Data: new(gofeed.Feed)}
	testFeed2 = Feed{ID: 2, URL: "2.com", Title: "2", Tags: []string{}, Data: new(gofeed.Feed)}
	testFeed3 = Feed{ID: 3, URL: "3.com", Title: "3", Tags: []string{}, Data: new(gofeed.Feed)}
	testFeed4 = Feed{ID: 4, URL: "4.com", Title: "4", Tags: []string{}, Data: new(gofeed.Feed)}
	testFeed5 = Feed{ID: 5, URL: "5.com", Title: "5", Tags: []string{}, Data: new(gofeed.Feed)}
	testFeeds = Feeds{&testFeed5, &testFeed4, &testFeed3, &testFeed2, &testFeed1}
)

func TestFeedsLen(t *testing.T) {
	total := 5
	if testFeeds.Len() != total {
		t.Errorf("Expected Len to be %d, but got %d", total, testFeeds.Len())
	}
}

func TestFeedsLess(t *testing.T) {
	if testFeeds.Less(1, 2) {
		t.Errorf("%s should not be less than %s", testFeeds[1].Title, testFeeds[2].Title)
	}
	if testFeeds.Less(1, 1) {
		t.Errorf("%s is not less than %s", testFeeds[1].Title, testFeeds[1].Title)
	}
	if !testFeeds.Less(2, 1) {
		t.Errorf("%s should be less than %s", testFeeds[2].Title, testFeeds[1].Title)
	}

}

func TestFeedsSwap(t *testing.T) {
	tempt1 := testFeeds[1]
	tempt2 := testFeeds[2]

	testFeeds.Swap(1, 2)

	if tempt1 != testFeeds[2] && tempt2 != testFeeds[1] {
		t.Error("Index 1 and 2 were not swapped")
	}
}

func TestLoadFeeds(t *testing.T) {
	//Make sure that the test database exist
	dbPath := filepath.Join(".", "test_data", "database", "feeds.db")
	dbPath, err := filepath.Abs(dbPath)
	if err != nil {
		t.Errorf("Failed to get the absolute path for database: %s", err.Error())
	}

	if !Exist(dbPath) {
		t.Errorf("test database for path (%s) does not exist", dbPath)
	}

	//Make sure that the urls test file exists
	urls := filepath.Join(".", "test_data", "urls", "urls")
	urls, err = filepath.Abs(urls)
	if err != nil {
		t.Errorf("Failed to get the absolute path for urls file: %s", err.Error())
	}

	//Setting the TestDB path and initializing the database
	TestDB = fmt.Sprintf("file:%s?_foreign_keys=1", dbPath)
	db, err := Init(TestDB, false)
	if err != nil {
		t.Errorf("Failed to initial database: %s", err.Error())
	}

	//Preparing needed data for LoadFeeds
	urlFileData := file.ExtractFileContent(urls)
	urlFileDataMap, err := AddFeedFileData(db, urlFileData)
	if err != nil {
		t.Errorf("failed to add file data to database: %s", err.Error())
	}

	feeds := LoadFeeds(db, urlFileDataMap)

	if len(feeds) <= 0 {
		t.Errorf("No feeds were loaded...")
	}

}

func TestNewFeeds(t *testing.T) {
	//TODO: Rethink how I do tests for the feeds.go file
	//file.Data variables
	goodData := file.Data{URL: "http://www.leoville.tv/podcasts/sn.xml", Tags: []string{"security", "favorite", "audio"}}
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
