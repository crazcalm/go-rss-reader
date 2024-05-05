package gui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

// Feed -- Data structure used to hold the needed feed data
type Feed struct {
	Episodes string
	Title    string
}

// Feeds -- Gui component for the feeds
type Feeds struct {
	Name  string
	Feeds []Feed
}

// NewFeeds -- Creates a new Feed gui component
func NewFeeds(name string, feeds []Feed) *Feeds {
	return &Feeds{Name: name, Feeds: feeds}
}

func (f *Feeds) format(g *gocui.Gui) (result string) {
	maxX, _ := g.Size()

	for index, feed := range f.Feeds {
		if index > 0 {
			result += "\n"
		}
		index, err := leftPadExactLength(fmt.Sprintf("%d", index+1), " ", 4)
		if err != nil {
			log.Fatal(err)
		}

		episodes, err := leftPadExactLength(feed.Episodes, " ", 12)
		if err != nil {
			log.Fatal(err)
		}

		title := leftPad(feed.Title, " ", 1)

		line, err := rightPadExactLength(index+episodes+title, " ", maxX)
		if err != nil {
			log.Fatal(err)
		}

		result += line
	}
	//Added extra lines so that I can use them to as
	//a marker to stop scrolling down
	result += "\n\n\n\n\n"
	return
}

func (f *Feeds) location(g *gocui.Gui) (x, y, w, h int) {
	maxX, maxY := g.Size()
	x = -1
	y = 1
	w = maxX
	h = maxY - 4
	return
}

// Layout -- Tells gocui.Gui how to display this component
func (f *Feeds) Layout(g *gocui.Gui) error {
	x, y, w, h := f.location(g)
	v, err := g.SetView(f.Name, x, y, w, h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		//gocui.Gui settings
		g.Cursor = true

		//gocui.View settings
		v.Frame = false
		v.Highlight = true
		v.SelBgColor = gocui.ColorBlue
		v.SelFgColor = gocui.ColorYellow

		//Setting this view on top
		_, err = g.SetCurrentView(v.Name())
		if err != nil {
			return err
		}

		fmt.Fprintf(v, f.format(g))
	}
	return nil
}
