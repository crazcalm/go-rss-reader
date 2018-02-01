package gui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

var (
	//CurrentFeed -- A marker for the current feed being viewed
	CurrentFeed = 0
	//FeedList -- List of all feeds
	FeedList = []string{}
)

//Feed -- Gui component for the feed as presented in the list
type Feed struct {
	name 		string
	index 		int
	episodes 	string
	title 		string
}

//NewFeed -- Creates a new Feed gui component
func NewFeed(name string, index int, episodes, title string) *Feed {
	FeedList = append(FeedList, name)
	return &Feed{name:name, index:index, episodes:episodes, title:title}
}

func (f *Feed) format() string {
	index, err := leftPadExactLength(fmt.Sprintf("%d", f.index), " ", 4)
	if err != nil {
		log.Fatal(err)
	}
	episodes, err := leftPadExactLength(f.episodes, " ", 12)
	if err != nil {
		log.Fatal(err) //Deal with properly later
	}
	title := leftPad(f.title, " ", 1)
	return index + episodes + title
}

func (f *Feed) location (g *gocui.Gui) (x, y, w, h int) {
	maxX, _ := g.Size()
	x = -1
	y = 2 * f.index - 1
	w = maxX
	h = y + 2
	return
}

//Layout -- Tells gocui.Gui how to display this component
func (f *Feed) Layout (g *gocui.Gui) (error) {
	maxX, _ := g.Size()
	x, y, w, h := f.location(g)
	v, err := g.SetView(f.name, x, y, w, h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		//gocui.View settings
		v.SelBgColor = gocui.ColorBlue
		v.SelFgColor = gocui.ColorYellow

		if v.Name() == "feed1" {
			g.SetViewOnTop(v.Name())
			v.Highlight = true
		}
		content, err := rightPadExactLength(f.format()," ", maxX)
		if err != nil {
			return err
		}

		fmt.Fprintf(v, content)
	}
	return nil
}
