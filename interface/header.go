package gui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

//Header -- Gui Component used for the title bar and info bar
type Header struct {
	name    string
	content string
}

//NewHeader -- Creates a new Bar gui component
func NewHeader(name, content string) *Header {
	return &Header{name: name, content: content}
}

func (header *Header) location(g *gocui.Gui) (x, y, w, h int) {
	maxX, _ := g.Size()
	x = -1
	y = -1
	w = maxX
	h = y + 2
	return
}

//Layout -- Tells gocui.Gui how to display this component
func (header *Header) Layout(g *gocui.Gui) error {
	x, y, w, h := header.location(g)
	v, err := g.SetView(header.name, x, y, w, h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		//Set colors
		v.FgColor = gocui.ColorYellow
		v.BgColor = gocui.ColorBlue

		//testing
		v.Frame = false

		_, err = fmt.Fprintf(v, header.content)
	}
	return err
}
