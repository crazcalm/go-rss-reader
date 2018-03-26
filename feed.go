package rss

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jroimartin/gocui"
	"github.com/mmcdole/gofeed"

	"github.com/crazcalm/go-rss-reader/interface"
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

	return &Episode{f.Title(), f.Data.Author, f.Data.Items[num]}, nil
}

//GetEpisodes -- Gets all episodes
func (f Feed) GetEpisodes() []*Episode {
	var episodes []*Episode

	for _, item := range f.Data.Items {
		episodes = append(episodes, &Episode{f.Title(), f.Data.Author, item})
	}

	return episodes
}

//EpisodesInit -- Episodes Init for the Gui
func EpisodesInit(g *gocui.Gui, feedIndex int) error {
	feed := FeedsData[feedIndex]

	//Episode data for one feed
	episodeData := feed.GuiItemsData()

	//Components
	header := gui.NewHeader("title", "Content goes here!")
	footer := gui.NewFooter("footer", "Footer Content is here!")
	episodes := gui.NewEpisodes("episodes", episodeData)

	g.SetManager(header, footer, episodes)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, gui.Quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, gui.CursorUp); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, gui.CursorDown); err != nil {
		log.Panicln(err)
	}

	return nil
}
