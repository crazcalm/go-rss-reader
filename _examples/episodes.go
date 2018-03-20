package main

import (
	"log"
	"path/filepath"

	"github.com/jroimartin/gocui"

	"github.com/crazcalm/go-rss-reader"
	"github.com/crazcalm/go-rss-reader/interface"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	//Get info from file
	fileData := rss.ExtractFileContent(filepath.Join("test_data", "urls"))

	//Create feeds
	feeds := rss.NewFeeds(fileData)

	//Episode data for one feed
	episodeData := feeds[0].GuiItemsData()

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

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
