package main

import (
	"log"
	"path/filepath"

	"github.com/crazcalm/go-rss-reader"
	"github.com/crazcalm/go-rss-reader/interface"
	"github.com/jroimartin/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	//Get info from file
	fileData := rss.ExtractFileContent(filepath.Join("test_data", "urls"))
	//log.Println(fileData)

	//Create feeds
	feeds := rss.NewFeeds(fileData)
	//log.Println(feeds)

	//Feed Gui info
	feedData := feeds.GuiData()
	//log.Println(feedData)

	/*
		//Testing
		feed1 := gui.Feed{"(12/15)", "title"}
		feed2 := gui.Feed{"(1/8)", "title"}
		feed3 := gui.Feed{"(110/159)", "title"}
		feed4 := gui.Feed{"(10/10)", "title"}
		feed5 := gui.Feed{"(0/0)", "title"}
		feed6 := gui.Feed{"(12/13)", "title"}
		feed7 := gui.Feed{"(0/10)", "title"}

		feedData := []gui.Feed{feed1, feed2, feed3, feed4, feed5, feed6, feed7}
	*/

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

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
