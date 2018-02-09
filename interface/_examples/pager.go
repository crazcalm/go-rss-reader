package main

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/crazcalm/go-rss-reader/interface"
	"github.com/jroimartin/gocui"
)

func main() {
	//get test data
	b, err := ioutil.ReadFile(filepath.Join("test_data", "pager_content.txt"))
	if err != nil {
		log.Panicln(err)
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	//components
	header := gui.NewHeader("title", "Content goes here")
	footer := gui.NewFooter("footer", "Content goes here")
	pager := gui.NewPager("pager", string(b))

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
