package rss

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"github.com/jroimartin/gocui"

	"github.com/crazcalm/go-rss-reader/database"
	"github.com/crazcalm/go-rss-reader/file"
	"github.com/crazcalm/go-rss-reader/interface"
)

var (
	//FeedsData -- Global container for the Feeds
	FeedsData database.Feeds
	//CurrentFeedIndex -- Global container for the current Feed index
	CurrentFeedIndex int
)

//FeedsInit -- Feeds Init for the Gui
func FeedsInit(g *gocui.Gui) error {
	var err error
	var db *sql.DB

	//fmt.Println("Init: Started")

	//Get info from file
	fileData := file.ExtractFileContent(filepath.Join("test_data", "urls"))

	//fmt.Print("FileData: ")
	//fmt.Println(fileData)

	//fmt.Println("Init: Extracted file content.")

	//Establish the database
	if database.Exist(database.TestDBPath) {
		db, err = database.Init(database.TestDB, false)
		if err != nil {
			log.Fatal(err)
		}

	} else {
		db, err = database.Create(database.TestDBPath)
		if err != nil {
			log.Fatal(err)
		}
		_, err = database.Init(database.TestDB, true)
		if err != nil {
			log.Fatal(err)
		}
	}

	//fmt.Println("Init: Established database")

	//Add file data to the database
	feedIDToFileDataMap, err := database.AddFeedFileData(fileData)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(feedIDToFileDataMap)

	//fmt.Println("Init: Added feed file data")

	//Create feeds
	//Note: This needs to be moved the load/refresh feed functionality
	//that does not exist yet
	//FeedsData = NewFeeds(feedIDToFileDataMap)

	//Load the Feeds
	FeedsData = database.LoadFeeds(db, feedIDToFileDataMap)

	//fmt.Println(FeedsData)

	//Feed Gui info
	feedData := FeedsData.GuiData(db)

	//fmt.Print("Init: feedData: ")
	//fmt.Println(feedData)

	//Components
	headerGui := gui.NewHeader("title", "Content goes here!")
	footerGui := gui.NewFooter("footer", "Footer Content is here!")
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

	//Refesh a feed
	if err := g.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone, UpdateFeed); err != nil {
		log.Panic(err)
	}

	return nil
}
