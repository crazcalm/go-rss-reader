package gui

import (
	"github.com/jroimartin/gocui"
)

//Quit -- Callback used to quit application
func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

//EpisodeDown -- Callback used to toggle down the Episode list
func EpisodeDown(g *gocui.Gui, v *gocui.View) error {
	if CurrentEpisode < len(EpisodeList)-1 {
		//Remove highlight from old view
		oldView, err := g.View(EpisodeList[CurrentEpisode])
		if err != nil {
			return err
		}
		oldView.Highlight = false

		//Add highlight to new view
		CurrentEpisode++
		newView, err := g.SetViewOnTop(EpisodeList[CurrentEpisode])
		if err != nil {
			return err
		}
		newView.Highlight = true

	}
	return nil
}

//FeedDown -- Callback used to toggle down the feed list
func FeedDown(g *gocui.Gui, v *gocui.View) error {
	if CurrentFeed < len(FeedList)-1 {
		//Remove highlight from old view
		oldView, err := g.View(FeedList[CurrentFeed])
		if err != nil {
			return err
		}
		oldView.Highlight = false

		//Add highlight to new view
		CurrentFeed++
		newView, err := g.SetViewOnTop(FeedList[CurrentFeed])
		if err != nil {
			return err
		}
		newView.Highlight = true

	}
	return nil
}

//EpisodeUp -- Callback used to toggle up the Episode list
func EpisodeUp(g *gocui.Gui, v *gocui.View) error {
	if CurrentEpisode > 0 {
		//Remove highlight from old view
		oldView, err := g.View(EpisodeList[CurrentEpisode])
		if err != nil {
			return err
		}
		oldView.Highlight = false

		//Add highlight to new view
		CurrentEpisode--
		newView, err := g.SetViewOnTop(EpisodeList[CurrentEpisode])
		if err != nil {
			return err
		}
		newView.Highlight = true

	}
	return nil
}

//FeedUp -- Callback used to toggle up the feed list
func FeedUp(g *gocui.Gui, v *gocui.View) error {
	if CurrentFeed > 0 {
		//Remove highlight from old view
		oldView, err := g.View(FeedList[CurrentFeed])
		if err != nil {
			return err
		}
		oldView.Highlight = false

		//Add highlight to new view
		CurrentFeed--
		newView, err := g.SetViewOnTop(FeedList[CurrentFeed])
		if err != nil {
			return err
		}
		newView.Highlight = true

	}
	return nil
}
