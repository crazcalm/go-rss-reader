package rss

import (
	"database/sql"
	"log"
	"path/filepath"
	"sort"

	"github.com/jroimartin/gocui"

	"github.com/crazcalm/go-rss-reader/database"
	"github.com/crazcalm/go-rss-reader/file"
	"github.com/crazcalm/go-rss-reader/interface"
)

var (
	//FeedsData -- Global container for the Feeds
	FeedsData database.Feeds
	//CurrentFeedID -- Global container for the current Feed ID
	CurrentFeedID int64
)

//FeedsInit -- Feeds Init for the Gui
func FeedsInit(g *gocui.Gui) error {
	var err error
	var db *sql.DB

	//Get info from file
	fileData := file.ExtractFileContent(filepath.Join("test_data", "urls"))

	//Establish the database
	if database.Exist(database.DBPath) {
		db, err = database.Init(database.TestDB, false)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		db, err = database.Create(database.DBPath)
		if err != nil {
			log.Fatal(err)
		}
		_, err = database.Init(database.TestDB, true)
		if err != nil {
			log.Fatal(err)
		}
	}

	//Add file data to the database
	feedIDToFileDataMap, err := database.AddFeedFileData(fileData)
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
	headerGui := gui.NewHeader("title", Header)
	footerGui := gui.NewFooter("footer", FeedFooter)
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
