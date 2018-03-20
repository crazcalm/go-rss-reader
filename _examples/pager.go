package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	//"github.com/jroimartin/gocui"

	"github.com/crazcalm/go-rss-reader"
	"github.com/crazcalm/html-to-text"
	//"github.com/crazcalm/go-rss-reader/interface"
)

func main() {
	/*
		g, err := gocui.NewGui(gocui.OutputNormal)
		if err != nil {
			log.Panicln(err)
		}
		defer g.Close()
	*/

	//Get info from file
	fileData := rss.ExtractFileContent(filepath.Join("test_data", "urls"))

	//Create feeds
	feeds := rss.NewFeeds(fileData)

	//Episode data for one feed
	episodeData := feeds[0].GuiItemsData()

	fmt.Println(episodeData)

	fmt.Printf("\nDescription: \n\n%s\n\nContent\n\n%s\n\n", feeds[0].Data.Items[0].Description, feeds[0].Data.Items[0].Content)

	text, _, err := htmltotext.Translate(strings.NewReader(feeds[0].Data.Items[0].Description))
	if err != nil {
		log.Panicln(err)
	}

	fmt.Println(strings.TrimSpace(text))

	//Components
	//header := gui.NewHeader("title", "Content goes here!")
	//footer := gui.NewFooter("footer", "Footer Content is here!")

}
