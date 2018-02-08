package rss

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

//Feeds -- A slice of feeds
type Feeds []*Feed

//Feed -- Data structure used to hold a feed
type Feed struct {
	URL  string
	Tags []string
	Data *gofeed.Feed
}

//NewFeed -- Used to create a new Feed
func NewFeed(fileData FileData) (*Feed, error) {
	fp := gofeed.NewParser()
	data, err := fp.ParseURL(fileData.URL)
	if err != nil {
		return &Feed{URL: fileData.URL, Tags: fileData.Tags, Data: data}, fmt.Errorf("Error occured while trying to parse feed")
	}
	return &Feed{URL: fileData.URL, Tags: fileData.Tags, Data: data}, nil
}

//NewFeeds -- Used to create a slice of Feeds
func NewFeeds(fileData []FileData) (feeds Feeds) {
	for _, d := range fileData {
		feed, err := NewFeed(d)
		if err != nil {
			continue
		}
		feeds = append(feeds, feed)
	}
	return feeds
}
