package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/crazcalm/go-rss-reader/file"
	"github.com/crazcalm/go-rss-reader/interface"
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

	body, err := ioutil.ReadAll(resp.Body)
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

//EpisodeTotal -- Returns the total num of episodes for the feed
func (f *Feed) EpisodeTotal() int {
	return len(f.Data.Items)
}

//GuiItemsData -- The data needed to create the episodes interface
func (f *Feed) GuiItemsData() (data []gui.Episode) {
	monthsMap := map[time.Month]string{
		1:  "Jan",
		2:  "Feb",
		3:  "Mar",
		4:  "Apr",
		5:  "May",
		6:  "Jun",
		7:  "Jul",
		8:  "Aug",
		9:  "Sep",
		10: "Oct",
		11: "Nov",
		12: "Dec",
	}

	for _, episode := range f.Data.Items {
		//Getting the information needed to format the date
		_, month, day := episode.PublishedParsed.Date()

		//If day is a single digit, add a leading zero
		if day < 10 {
			data = append(data, gui.Episode{Date: fmt.Sprintf("%s 0%d", monthsMap[month], day), Title: episode.Title})
			continue
		}

		data = append(data, gui.Episode{Date: fmt.Sprintf("%s %d", monthsMap[month], day), Title: episode.Title})
	}

	return
}

//GetEpisode -- Gets one episode [Counting starts from 0]
func (f Feed) GetEpisode(num int) (*Episode, error) {
	episode := &Episode{"", &gofeed.Person{}, &gofeed.Item{}}
	err := fmt.Errorf("No episode found")

	if num < 0 {
		return episode, err
	} else if num >= len(f.Data.Items) {
		return episode, err
	}

	return &Episode{f.Title, f.Data.Author, f.Data.Items[num]}, nil
}

//GetEpisodes -- Gets all episodes
func (f Feed) GetEpisodes() []*Episode {
	var episodes []*Episode

	for _, item := range f.Data.Items {
		episodes = append(episodes, &Episode{f.Title, f.Data.Author, item})
	}

	return episodes
}
