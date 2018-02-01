package main

import (
	"log"

	"github.com/crazcalm/go_read_rss/interface"
	"github.com/jroimartin/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	//Components
	header := gui.NewHeader("title", "Content goes here!")
	footer := gui.NewFooter("footer", "Footer Content is here!")
	feed1  := gui.NewFeed("feed1", 1, "(12/15)", "title")
	feed2  := gui.NewFeed("feed2", 2, "(1/8)", "title")
	feed3  := gui.NewFeed("feed3", 3,"(110/159)", "title")
	feed4  := gui.NewFeed("feed4", 4, "(10/10)", "title")
	feed5  := gui.NewFeed("feed5", 5, "(0/0)", "title")
	feed6  := gui.NewFeed("feed6", 6, "(12/13)", "title")
	feed7  := gui.NewFeed("feed7", 7, "(0/10)", "title")

	g.SetManager(header, footer, feed1, feed2, feed3, feed4, feed5, feed6, feed7)

    if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, gui.Quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, gui.FeedUp); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, gui.FeedDown); err != nil {
		log.Panicln(err)
	}
	
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	 }	
}