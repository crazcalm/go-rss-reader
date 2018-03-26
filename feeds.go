package rss

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/jroimartin/gocui"

	"github.com/crazcalm/go-rss-reader/interface"
)

var (
	//FeedsData -- Global container for the Feeds
	FeedsData Feeds
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

//GuiData -- The data needed to create the gui interface for feeds
func (f Feeds) GuiData() (data []gui.Feed) {
	for _, item := range f {
		data = append(data, gui.Feed{fmt.Sprintf("(%d/%d)", item.EpisodeTotal(), item.EpisodeTotal()), item.Title()})
	}
	return
}

//FeedsInit -- Feeds Init for the Gui
func FeedsInit(g *gocui.Gui) error {
	//Get info from file
	fileData := ExtractFileContent(filepath.Join("test_data", "urls"))

	//Create feeds
	FeedsData = NewFeeds(fileData)

	//Feed Gui info
	feedData := FeedsData.GuiData()

	//Components
	headerGui := gui.NewHeader("title", "Content goes here!")
	footerGui := gui.NewFooter("footer", "Footer Content is here!")
	feedsGui := gui.NewFeeds("feeds", feedData)

	g.SetManager(headerGui, footerGui, feedsGui)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, gui.Quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, gui.CursorUp); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, gui.CursorDown); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, SelectFeed); err != nil {
		log.Panic(err)
	}

	return nil
}
