package rss

import (
	"fmt"
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

		parts := strings.SplitN(strings.TrimSpace(line), ")", 2) //Splitting till I reach the title
		if len(parts) <= 0 {
			return fmt.Errorf("The selected line does not have enough parts to split: %s", line)
		}
		title = strings.TrimSpace(parts[len(parts)-1])
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
	return err
}


//UpdateFeeds -- Update all feeds being shown
func UpdateFeeds(g *gocui.Gui, v *gocui.View) error {
	for _, feed := range FeedsData {
		err := database.GetFeedInfo(database.DB, feed.ID)
		if err != nil {
			return err
		}
		err = FeedsInit(g)
		if err != nil {
			return err
		}
	}
	return nil
}


//SelectFeed -- Callback used to select an episode
func SelectFeed(g *gocui.Gui, v *gocui.View) (err error) {
	var title string
	var line string

	if v != nil {
		_, cy := v.Cursor()

		line, err = v.Line(cy)
		if err != nil {
			return err
		}

		parts := strings.SplitN(strings.TrimSpace(line), ")", 2) //Splitting till I reach the title
		if len(parts) <= 0 {
			return fmt.Errorf("The selected line does not have enough parts to split: %s", line)
		}
		title = strings.TrimSpace(parts[len(parts)-1])
	}

	//Get feedID
	if strings.EqualFold(title, "") {
		return fmt.Errorf("The title obtained was an empty string")
	}
	CurrentFeedID, err = database.GetFeedID(database.DB, title)
	if err != nil {
		return err
	}

	err = EpisodesInit(g)
	return err
}

//SelectEpisode -- Callback used to select a feed
func SelectEpisode(g *gocui.Gui, v *gocui.View) (err error) {
	var title string
	var line string
	titleStartIndex := 21

	if v != nil {
		_, cy := v.Cursor()

		line, err = v.Line(cy)
		if err != nil {
			return err
		}

		title = strings.TrimSpace(line[titleStartIndex:])
	}

	CurrentEpisodeID, err = database.GetEpisodeIDByFeedIDAndTitle(database.DB, CurrentFeedID, title)
	if err != nil {
		return err
	}

	err = EpisodeContentInit(g)
	return err
}
