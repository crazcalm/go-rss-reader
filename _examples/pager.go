package main

import (
	"log"
	"path/filepath"

	"github.com/jroimartin/gocui"

	"github.com/crazcalm/go-rss-reader"
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
	rss.FeedsData = rss.NewFeeds(fileData)

	err = rss.EpisodeContentInit(g, 0, 0)
	if err != nil {
		log.Fatal(err)
	}

	//Run code
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

}
