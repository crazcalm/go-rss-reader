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

	//Episode Content for one feed
	episode, err := feeds[0].GetEpisode(0)
	if err != nil {
		log.Panicln(err)
	}

	content, _, err := episode.Content()
	if err != nil {
		log.Panicln(err)
	}

	//Components
	header := gui.NewHeader("title", "Content goes here")
	footer := gui.NewFooter("footer", "Content goes here")
	pager := gui.NewPager("pager", content)

	//Display components
	g.SetManager(header, footer, pager)

	//keybindings
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, gui.Quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("pager", gocui.KeyArrowUp, gocui.ModNone, gui.PagerUp); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("pager", gocui.KeyArrowDown, gocui.ModNone, gui.PagerDown); err != nil {
		log.Panicln(err)
	}

	//Run code
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}
