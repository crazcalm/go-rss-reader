package main

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"

	"github.com/crazcalm/go-rss-reader/database"
	"github.com/crazcalm/go-rss-reader"
	"github.com/crazcalm/go-rss-reader/cli"
)

func main() {
	cli.GlobalConfig = cli.NewConfig()
	config := cli.GlobalConfig
	err := config.CliParse()
	if err != nil {
		log.Panicln(err)
	}
	

	if config.DBExist() {
		_, err = database.Init(fmt.Sprintf("file:%s?_foreign_keys=1", config.GetDBPath()), false)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		_, err = database.Create(config.GetDBPath())
		if err != nil {
			log.Fatal(err)
		}
		
	}

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	err = rss.FeedsInit(g)
	if err != nil {
		log.Fatal(err)
	}

	if err = g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
