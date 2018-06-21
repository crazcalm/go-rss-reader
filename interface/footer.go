package gui

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

//Footer -- Gui Component used for the title bar and info bar
type Footer struct {
	name    string
	content string
}

//NewFooter -- Creates a new Bar gui component
func NewFooter(name, content string) *Footer {
	return &Footer{name: name, content: content}
}

func (f *Footer) location(g *gocui.Gui) (x, y, w, h int) {
	maxX, maxY := g.Size()
	x = -1
	y = maxY - 4
	w = maxX
	h = y + 2
	return
}

//Layout -- Tells gocui.Gui how to display this component
func (f *Footer) Layout(g *gocui.Gui) error {
	x, y, w, h := f.location(g)
	v, err := g.SetView(f.name, x, y, w, h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		//Set colors
		v.FgColor = gocui.ColorYellow
		v.BgColor = gocui.ColorBlue

		//gocui.View Settings
		v.Frame = false

		_, err = fmt.Fprintf(v, f.content)
	}
	return err
}
