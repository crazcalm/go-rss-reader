package rss

import (
	"log"

	"github.com/jroimartin/gocui"

	"github.com/crazcalm/go-rss-reader/database"
	"github.com/crazcalm/go-rss-reader/interface"
)

var (
	//CurrentEpisodeIndex -- Global container for the current episode index
	CurrentEpisodeIndex int
)

//EpisodesInit -- Episodes Init for the Gui
func EpisodesInit(g *gocui.Gui) error {
	episodeIDs, err := database.GetFeedEpisodeIDs(database.DB, CurrentFeedID)
	if err != nil {
		return err
	}

	//Create data need for the Epsisode screen
	var guiEpisodeData []gui.Episode
	for _, id := range episodeIDs {
		_, title, date, _, _, err := database.GetEpisode(database.DB, id)
		if err != nil {
			return err
		}

		guiEpisodeData = append(guiEpisodeData, gui.Episode{Date: date.Format("Jan 02"), Title: title})
	}

	//Components
	header := gui.NewHeader("title", "Content goes here!")
	footer := gui.NewFooter("footer", "Footer Content is here!")
	episodes := gui.NewEpisodes("episodes", guiEpisodeData)

	g.SetManager(header, footer, episodes)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, gui.Quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, gui.CursorUp); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, gui.CursorDown); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, SelectEpisode); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlB, gocui.ModNone, QuitEpisodes); err != nil {
		log.Panicln(err)
	}

	return nil
}

//EpisodeContentInit -- Initializes the Episode content for the Gui
func EpisodeContentInit(g *gocui.Gui) error {
	/*
		feed := FeedsData[CurrentFeedIndex]
		episode, err := feed.GetEpisode(CurrentEpisodeIndex)
		if err != nil {
			log.Fatal(err)
		}

		content, _, err := episode.Content()
		if err != nil {
			log.Panicln(err)
		}
	*/

	//Components
	header := gui.NewHeader("title", "Content goes here")
	footer := gui.NewFooter("footer", "Content goes here")
	pager := gui.NewPager("pager", "pending")

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

	if err := g.SetKeybinding("pager", gocui.KeyCtrlB, gocui.ModNone, QuitPager); err != nil {
		log.Panicln(err)
	}

	return nil
}
