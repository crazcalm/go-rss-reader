package gui

import (
	"strings"

	"github.com/jroimartin/gocui"
)

//Quit -- Callback used to quit application
func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

//CursorDown -- Callback used to scroll down on the feeds
func CursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()

		//Used to check if we have reached the last item
		nextLine, err := v.Line(cy + 1)
		if err != nil {
			return err
		}
		if strings.EqualFold(nextLine, "") {
			return nil
		}

		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

//PagerDown -- Callback used to scroll down on the Pager
func PagerDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy+1); err != nil {

			//Need to stop scrolling once I reach the end
			str, err := v.Line(1)
			if err != nil {
				return err
			}

			//Case: no more content
			if strings.EqualFold(str, "~") {
				return nil
			}

			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

//CursorUp -- Callback used to scroll up the Feeds
func CursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}

//PagerUp -- Callback used to scroll up on the Pager
func PagerUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()
		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}
