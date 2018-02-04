package gui

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)

//Pager -- Gui component used for viewing the content of a feed episode
type Pager struct {
	name    string
	content string
}

//NewPager -- Creates a new gui Pager component
func NewPager(name, content string) *Pager {
	return &Pager{name: name, content: content}
}

func (p *Pager) format() string {
	tail := strings.Repeat("\n~", 100)
	return p.content + tail
}

func (p *Pager) location(g *gocui.Gui) (x, y, w, h int) {
	maxX, maxY := g.Size()
	x = -1
	y = 1 //Needs to be right below the header
	w = maxX
	h = maxY - 4 //Needs to be right above the footer

	return
}

//Layout -- Tells gocui.Gui how to display this component
func (p *Pager) Layout(g *gocui.Gui) error {
	x, y, w, h := p.location(g)
	v, err := g.SetView(p.name, x, y, w, h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		//Gui settings
		g.Cursor = true

		//view settings
		v.Frame = false

		//Setting this view on top
		_, err = g.SetCurrentView(v.Name())
		if err != nil {
			return err
		}

		fmt.Fprint(v, p.format())
	}
	return nil
}
