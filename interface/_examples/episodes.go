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

	ep1 := gui.NewEpisode("ep1", 1, "Jan 01", "title")
	ep2 := gui.NewEpisode("ep2", 2, "Feb 23", "title")
	ep3 := gui.NewEpisode("ep3", 3, "Mar 15", "title")
	ep4 := gui.NewEpisode("ep4", 4, "Sep 17", "title")

	g.SetManager(header, footer, ep1, ep2, ep3, ep4)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, gui.Quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, gui.EpisodeUp); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, gui.EpisodeDown); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
