package rss

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/crazcalm/go-rss-reader/database"
	"github.com/jroimartin/gocui"
)

//QuitPager -- Callback used to quit the pager view and return to the Episodes view
func QuitPager(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		err := EpisodesInit(g)
		if err != nil {
			return err
		}
	}

	return nil
}

//QuitEpisodes -- Callback used to quit the Episodes view and return to the feeds view
func QuitEpisodes(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		err := FeedsInit(g)
		if err != nil {
			return err
		}
	}
	return nil
}

//UpdateFeed -- Updates a single feed
func UpdateFeed(g *gocui.Gui, v *gocui.View) error {
	var title string

	if v != nil {
		_, cy := v.Cursor()

		line, err := v.Line(cy)
		if err != nil {
			return err
		}

		parts := strings.Split(strings.TrimSpace(line), " ")
		if len(parts) <= 0 {
			return fmt.Errorf("The selected line does not have enough parts to split: %s", line)
		}
		lastIndex := len(parts) - 1
		title = parts[lastIndex]
	}

	//Get feedID
	if strings.EqualFold(title, "") {
		return fmt.Errorf("The title obtained was an empty string")
	}
	feedID, err := database.GetFeedID(database.DB, title)
	if err != nil {
		return err
	}

	//UpdateFeed
	err = database.GetFeedInfo(database.DB, feedID)
	if err != nil {
		return err
	}

	//Refresh screen
	err = FeedsInit(g)
	if err != nil {
		return err
	}

	return err
}

//SelectEpisode -- Callback used to select an episode
func SelectEpisode(g *gocui.Gui, v *gocui.View) error {
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

		//Set Global Current Episode Index
		CurrentEpisodeIndex = feedIndex - 1 // Minus 1 to get the actual index

		err = EpisodeContentInit(g)
		if err != nil {
			return err
		}
	}
	return nil
}

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

		//Set CurrentFeedIndex
		CurrentFeedIndex = feedIndex - 1 // Minus 1 to get the actual index

		err = EpisodesInit(g)
		if err != nil {
			return err
		}

	}
	return nil
}
