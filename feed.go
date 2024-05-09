package rss

import (
	"log"
	"path/filepath"
	"sort"

	"github.com/jroimartin/gocui"

	"github.com/crazcalm/go-rss-reader/cli"
	"github.com/crazcalm/go-rss-reader/database"
	"github.com/crazcalm/go-rss-reader/file"
	"github.com/crazcalm/go-rss-reader/interface"
)

var (
	//FeedsData -- Global container for the Feeds
	FeedsData database.Feeds
	//CurrentFeedID -- Global container for the current Feed ID
	CurrentFeedID int64
	//URLFile -- is the path to the url file
	URLFile = filepath.Join(".go-rss-reader", "urls")
)

// FeedsInitWithConfig -- Feeds Init for the Gui
func FeedsInitWithConfig(g *gocui.Gui) error {
	var err error

	//Get info from file
	fileData := file.ExtractFileContent(cli.GlobalConfig.GetUrlFile())

	// TODO: refactor this properly
	db := database.DB
	
	//Add file data to the database
	feedIDToFileDataMap, err := database.AddFeedFileData2(fileData, db)
	if err != nil {
		log.Fatal(err)
	}

	//Load the Feeds
	FeedsData = database.LoadFeeds(db, feedIDToFileDataMap)

	//Sort FeedsData by Title
	sort.Sort(FeedsData)

	//Feed Gui info
	feedData := FeedsData.GuiData(db)

	//Components
	headerGui := gui.NewHeader("title", gui.HeaderText)
	footerGui := gui.NewFooter("footer", gui.FeedsFooterText)
	feedsGui := gui.NewFeeds("feeds", feedData)

	g.SetManager(headerGui, footerGui, feedsGui)

	//Quit
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, gui.Quit); err != nil {
		log.Panicln(err)
	}

	//Scroll up
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, gui.CursorUp); err != nil {
		log.Panicln(err)
	}

	//Scroll down
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, gui.CursorDown); err != nil {
		log.Panicln(err)
	}

	//Select Feed
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, SelectFeed); err != nil {
		log.Panic(err)
	}

	//Refresh a feed
	if err := g.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, UpdateFeed); err != nil {
		log.Panic(err)
	}

	//Refresh all the feeds
	if err := g.SetKeybinding("", gocui.KeyCtrlA, gocui.ModNone, UpdateFeeds); err != nil {
		log.Panic(err)
	}

	return nil
}
