package gui

import (
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

//Episode -- Gui component for the episodes of a feed
type Episode struct {
	name	string
	index	int
	date	string
	title	string
}

//NewEpisode -- Creates a new Episode gui component
func NewEpisode(name string, index int, date, title string) *Episode{
	return &Episode{name:name, index:index, date:date, title:title}
}

func (e *Episode) format() string {
	index, err := leftPadExactLength(fmt.Sprintf("%d", e.index), " ", 4)
	if err != nil {
		log.Fatal(err)
	}
	date, err := leftPadExactLength(e.date, " ", 8)
	if err != nil {
		log.Fatal(err) //Deal with properly later
	}
	title := leftPad(e.title, " ", 2)
	return index + date + title
}

func (e *Episode) location (g *gocui.Gui) (x, y, w, h int) {
	maxX, _ := g.Size()
	x = -1
	y = 2 * e.index - 1
	w = maxX
	h = y + 2
	return
}

//Layout -- Tells gocui.Gui how to display this component
func (e *Episode) Layout (g *gocui.Gui) (error) {
	x, y, w, h := e.location(g)
	v, err := g.SetView(e.name, x, y, w, h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		fmt.Fprintf(v, e.format())
	}
	return nil
}