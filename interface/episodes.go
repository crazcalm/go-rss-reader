package gui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

//Episode -- Data structure used to hold the needed episode data
type Episode struct {
	Date  string
	Title string
	Seen  bool
}

//Episodes -- Gui component for the Episodes
type Episodes struct {
	Name     string
	Episodes []Episode
}

//NewEpisodes -- Creates a new Episode gui component
func NewEpisodes(name string, episodes []Episode) *Episodes {
	return &Episodes{Name: name, Episodes: episodes}
}

func (e *Episodes) format(g *gocui.Gui) (result string) {
	maxX, _ := g.Size()

	for index, ep := range e.Episodes {
		if index > 0 {
			result += "\n"
		}

		index, err := leftPadExactLength(fmt.Sprintf("%d", index+1), " ", 4)
		if err != nil {
			log.Fatal(err)
		}

		var seen string
		if ep.Seen {
			seen, err = leftPadExactLength("seen", " ", 8)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			seen, err = leftPadExactLength("unseen", " ", 8)
			if err != nil {
				log.Fatal(err)
			}
		}

		date, err := leftPadExactLength(ep.Date, " ", 8)
		if err != nil {
			log.Fatal(err)
		}

		title := leftPad(ep.Title, " ", 2)

		line, err := rightPadExactLength(index+seen+date+title, " ", maxX)
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

func (e *Episodes) location(g *gocui.Gui) (x, y, w, h int) {
	maxX, maxY := g.Size()
	x = -1
	y = 1
	w = maxX
	h = maxY - 4
	return
}

//Layout -- Tells gocui.Gui how to display this component
func (e *Episodes) Layout(g *gocui.Gui) error {
	x, y, w, h := e.location(g)
	v, err := g.SetView(e.Name, x, y, w, h)
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

		fmt.Fprint(v, e.format(g))
	}
	return nil
}
