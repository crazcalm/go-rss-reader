package rss

import (
	"fmt"
	"log"
	"strings"

	"github.com/jroimartin/gocui"

	"github.com/crazcalm/go-rss-reader/database"
	"github.com/crazcalm/go-rss-reader/interface"
	"github.com/crazcalm/html-to-text"
)

var (
	//CurrentEpisodeID -- Global reference to the episode ID of the currently viewed
	//or last viewed episode
	CurrentEpisodeID int64
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
		_, title, date, seen, _, err := database.GetEpisode(database.DB, id)
		if err != nil {
			return err
		}

		if seen == 1 {
			guiEpisodeData = append(guiEpisodeData, gui.Episode{Date: date.Format("Jan 02"), Title: title, Seen: true})

		} else {
			guiEpisodeData = append(guiEpisodeData, gui.Episode{Date: date.Format("Jan 02"), Title: title, Seen: false})

		}
	}

	//Components
	header := gui.NewHeader("title", gui.HeaderText)
	footer := gui.NewFooter("footer", gui.EpisodesFooterText)
	episodes := gui.NewEpisodes("episodes", guiEpisodeData)

	g.SetManager(header, footer, episodes)

	//Quit
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, gui.Quit); err != nil {
		log.Panicln(err)
	}

	//Scroll up
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, gui.CursorUp); err != nil {
		log.Panicln(err)
	}

	//Scroll down
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, gui.CursorDown); err != nil {
		log.Panicln(err)
	}

	//Select Episode
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, SelectEpisode); err != nil {
		log.Panicln(err)
	}

	//Back
	if err := g.SetKeybinding("", gocui.KeyCtrlB, gocui.ModNone, QuitEpisodes); err != nil {
		log.Panicln(err)
	}

	return nil
}

//EpisodeContentInit -- Initializes the Episode content for the Gui
func EpisodeContentInit(g *gocui.Gui) error {
	feedTitle, episodeTitle, author, episodeLink, date, mediaContent, err := database.GetEpisodeHeaderData(database.DB, CurrentFeedID, CurrentEpisodeID)
	if err != nil {
		return err
	}

	episodeHeader := database.FormatEpisodeHeader(feedTitle, episodeTitle, author, episodeLink, date, mediaContent)

	_, _, _, _, rawData, err := database.GetEpisode(database.DB, CurrentEpisodeID)
	if err != nil {
		return err
	}

	content, links, err := htmltotext.Translate(strings.NewReader(rawData))
	if err != nil {
		return fmt.Errorf("Error occurred when parsing raw data: (%s). Returning raw data", err.Error())
	}

	content = strings.TrimSpace(content)
	//If the rawData was not html, content will have spaces, but no data
	if strings.EqualFold(content, "") {
		content = rawData
	}

	if len(links) != 0 {
		//Add the links to the content
		content += "\n\n\nLinks:\n"
		for index, link := range links {
			content += fmt.Sprintf("%d: %s\n", index, link)
		}
	}

	body := fmt.Sprintf("%s\n%s", episodeHeader, strings.TrimSpace(content))

	//Mark episode as seen
	err = database.MarkEpisodeAsSeen(database.DB, CurrentEpisodeID)
	if err != nil {
		return err
	}

	//Components
	header := gui.NewHeader("title", gui.HeaderText)
	footer := gui.NewFooter("footer", gui.PagerFooterText)
	pager := gui.NewPager("pager", body)

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
