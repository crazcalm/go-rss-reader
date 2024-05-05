package database

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/mmcdole/gofeed"

	"github.com/crazcalm/go-rss-reader/file"
)

//Feed -- Data structure used to hold a feed
type Feed struct {
	ID    int64
	URL   string
	Title string
	Tags  []string
	Data  *gofeed.Feed
}

//GetFeedDataFromSite -- gets the feed data from the feed url and returns it
func GetFeedDataFromSite(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Error trying to get the raw feed data from %s: %s", url, err.Error())
	}
	defer func() {
		if err = resp.Body.Close(); err != nil {
			err = fmt.Errorf("Errorr occurred while closing the response body: %s", err.Error())
		}
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), err
}

//NewFeed -- Used to create a new Feed. Id the id is equal to -1, then
//all of the database interactions will not happen
func NewFeed(id int64, fileData file.Data) (*Feed, error) {
	var db *sql.DB
	var err error
	if id != -1 {
		//Initialize database
		db, err = Init(TestDB, false)
		if err != nil {
			log.Fatal(err)
		}
	}

	//Parse the feed
	fp := gofeed.NewParser()
	data, err := fp.ParseURL(fileData.URL)
	if err != nil {
		return &Feed{ID: id, URL: fileData.URL, Tags: fileData.Tags, Data: data}, fmt.Errorf("Error occured while trying to parse feed")
	}

	//Initalize a feed
	feed := &Feed{ID: id, URL: fileData.URL, Tags: fileData.Tags, Data: data}

	if id != -1 {

		//get raw data for feed
		rawData, err := GetFeedDataFromSite(feed.URL)
		if err != nil {
			log.Fatal(err)
		}

		//Add the feed data to the database
		err = UpdateFeedRawData(db, feed.ID, rawData)
		if err != nil {
			log.Fatal(err)
		}

		err = UpdateFeedTitle(db, feed.ID, feed.Title)
		if err != nil {
			log.Fatal(err)
		}
	}

	return feed, nil
}
