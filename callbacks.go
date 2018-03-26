package rss

import (
	"strconv"
	"strings"

	"github.com/jroimartin/gocui"
)

//SelectFeed -- Callback used to select a feed
func SelectFeed(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, cy := v.Cursor()

		line, err := v.Line(cy)
		if err != nil {
			return err
		}

		//Parsing for the index number on the line
		parts := strings.Split(strings.TrimSpace(line), " ")
		feedIndex, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			return err
		}

		err = EpisodesInit(g, feedIndex-1) // Minus 1 to get the actual index
		if err != nil {
			return err
		}

	}
	return nil
}
