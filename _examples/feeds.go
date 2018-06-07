package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"

	"github.com/crazcalm/go-rss-reader"
	"github.com/crazcalm/go-rss-reader/database"
)

var dbPath string
var urlPath string

func main() {
	//Note: You may set the default paths to your database and urls file here
	flag.StringVar(&dbPath, "db", "test.db", "Path to the database")
	flag.StringVar(&urlPath, "url", "test_data/urls", "Path to the url file ")
	flag.Parse()
	if len(strings.TrimSpace(dbPath)) == 0 {
		log.Fatal(fmt.Errorf("Did not pass in a database"))
	}
	database.DBPath = dbPath

	if len(strings.TrimSpace(urlPath)) == 0 {
		log.Fatal(fmt.Errorf("Did not pass in a url file"))
	}
	rss.URLFile = urlPath

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
