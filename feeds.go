package rss

import (
	"github.com/crazcalm/go_read_rss/interface"
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
func (f Feeds)GuiData()(data []gui.Feed){
	return
}