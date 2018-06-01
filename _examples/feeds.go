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

func main() {
	flag.StringVar(&dbPath, "db", "test.db", "Path the database")
	flag.Parse()
	if len(strings.TrimSpace(dbPath)) == 0 {
		log.Fatal(fmt.Errorf("Did not pass in a database"))
	}
	database.DBPath = dbPath

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
