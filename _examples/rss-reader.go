package main

import (
	"log"

	"github.com/jroimartin/gocui"

	"github.com/crazcalm/go-rss-reader"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	err = rss.FeedsInit(g)
	if err != nil {
		log.Fatal(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
