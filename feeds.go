package rss

import (
	"fmt"

	"github.com/crazcalm/go-rss-reader/interface"
)

//Feeds -- A slice of feeds
type Feeds []*Feed

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

//GuiData -- Data needs to create the gui interface for feeds
func (f Feeds) GuiData() (data []gui.Feed) {
	for _, item := range f {
		data = append(data, gui.Feed{fmt.Sprintf("(%d/%d)", item.EpisodeTotal(), item.EpisodeTotal()), item.Title()})
	}
	return
}
