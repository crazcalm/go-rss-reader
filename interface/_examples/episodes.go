package main

import (
	"log"

	"github.com/crazcalm/go-rss-reader/interface"
	"github.com/jroimartin/gocui"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	//Episodes
	ep1 := gui.Episode{"Jan 01", "title", true}
	ep2 := gui.Episode{"Feb 23", "title", false}
	ep3 := gui.Episode{"Mar 15", "title", false}
	ep4 := gui.Episode{"Sep 17", "title", true}

	//Components
	header := gui.NewHeader("title", gui.HeaderText)
	footer := gui.NewFooter("footer", gui.EpisodesFooterText)
	episodes := gui.NewEpisodes("episodes", []gui.Episode{ep1, ep2, ep3, ep4})

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
