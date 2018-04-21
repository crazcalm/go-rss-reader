package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/crazcalm/go-rss-reader/file"
	"github.com/crazcalm/go-rss-reader/interface"
)

//Feeds -- A slice of feeds
type Feeds []*Feed

//LoadFeeds -- Used to load the feed info from the database
func LoadFeeds(db *sql.DB, fileData map[int64]file.Data) (feeds Feeds) {
	for id := range fileData {
		feed, err := LoadFeed(db, id)
		if err != nil {
			log.Fatal(err)
		}
		feeds = append(feeds, feed)
	}
	return feeds
}

//NewFeeds -- Used to create a slice of Feeds
func NewFeeds(fileData map[int64]file.Data) (feeds Feeds) {
	for id, data := range fileData {
		feed, err := NewFeed(id, data)
		if err != nil {
			continue
		}
		feeds = append(feeds, feed)
	}
	return feeds
}

//GuiData -- The data needed to create the gui interface for feeds
func (f Feeds) GuiData(db *sql.DB) (data []gui.Feed) {
	for _, feed := range f {
		seen, total, err := GetFeedEpisodeSeenRatio(db, feed.ID)
		if err != nil {
			log.Fatal(err)
		}
		//TODO: fix feed.Title() -- I should be able to call it here!!!
		data = append(data, gui.Feed{fmt.Sprintf("(%d/%d)", seen, total), feed.Title})
	}
	return
}
