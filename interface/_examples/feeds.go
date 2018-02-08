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

	feed1 := gui.Feed{"(12/15)", "title"}
	feed2 := gui.Feed{"(1/8)", "title"}
	feed3 := gui.Feed{"(110/159)", "title"}
	feed4 := gui.Feed{"(10/10)", "title"}
	feed5 := gui.Feed{"(0/0)", "title"}
	feed6 := gui.Feed{"(12/13)", "title"}
	feed7 := gui.Feed{"(0/10)", "title"}

	feedData := []gui.Feed{feed1, feed2, feed3, feed4, feed5, feed6, feed7}

	//Components
	header := gui.NewHeader("title", "Content goes here!")
	footer := gui.NewFooter("footer", "Footer Content is here!")
	feeds := gui.NewFeeds("feeds", feedData)

	g.SetManager(header, footer, feeds)

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
