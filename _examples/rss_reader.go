package main

import (
	"fmt"
	"log"
	"log/slog"
	
	"github.com/jroimartin/gocui"

	"github.com/crazcalm/go-rss-reader/database"
	"github.com/crazcalm/go-rss-reader"
	"github.com/crazcalm/go-rss-reader/cli"
)

func main() {
	slog.Debug("Starting the program")
	
	cli.GlobalConfig = cli.NewConfig()
	err := cli.GlobalConfig.CliParse()
	if err != nil {
		log.Fatalf("Failed to parse cli args: %s", err.Error())
	}

	slog.Debug("Successfully Parsed the Cli Args")
	
	if cli.GlobalConfig.DBExist() {
		_, err = database.Init(fmt.Sprintf("file:%s?_foreign_keys=1", cli.GlobalConfig.GetDBPath()), false)
		if err != nil {
			log.Fatalf("Unable to connect to DB at %s: %s", cli.GlobalConfig.GetDBPath(), err.Error())
		}
	} else {
		_, err = database.Create(cli.GlobalConfig.GetDBPath())
		if err != nil {
			log.Fatal("Unable to create DB at %s: %s", cli.GlobalConfig.GetDBPath(), err)
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
