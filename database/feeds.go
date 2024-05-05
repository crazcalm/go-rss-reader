package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/crazcalm/go-rss-reader/file"
	"github.com/crazcalm/go-rss-reader/interface"
)

//Feeds -- A slice of feeds
type Feeds []*Feed

//Len -- Returns the number of feeds in the slice.
//Needed for the sort interface.
func (f Feeds) Len() int { return len(f) }

//Less -- Does a string comparison on the Title.
//Needed for the sort interface
func (f Feeds) Less(i, j int) bool {
	return strings.Compare(f[i].Title, f[j].Title) == -1
}

//Swap -- Definining what it means to swap two feeds
//Needed for the sort interface
func (f Feeds) Swap(i, j int) { f[i], f[j] = f[j], f[i] }

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
		data = append(data, gui.Feed{Episodes: fmt.Sprintf("(%d/%d)", seen, total), Title: feed.Title})
	}
	return
}
