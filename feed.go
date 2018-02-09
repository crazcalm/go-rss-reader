package rss

import (
	"fmt"
	"strings"

	"github.com/mmcdole/gofeed"
)

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

//Title -- returns the title of the feed. If there is no title,
//Then it returns the url of the feed.
func (f *Feed) Title() string {
	if !strings.EqualFold(f.Data.Title, "") {
		return f.Data.Title
	}
	return f.URL
}

//EpisodeTotal -- Returns the total num of episodes for the feed
func (f *Feed) EpisodeTotal() int {
	return len(f.Data.Items)
}
