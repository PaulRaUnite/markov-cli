package actions

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

var errReturn string

func errSpawn(g *gocui.Gui, text string, ret string) error {
	errReturn = ret
	maxX, maxY := g.Size()
	xShift, yShift := 20, 5
	if v, err := g.SetView("error_box", maxX/2-xShift, maxY/2-yShift, maxX/2+xShift, maxY/2+yShift); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Error"
	}
	if v, err := g.SetView("error_text", maxX/2-xShift+2, maxY/2-yShift+1, maxX/2+xShift-2, maxY/2+yShift-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		fmt.Fprintln(v, text)
	}
	if v, err := g.SetView("error_button", maxX/2-xShift+5, maxY/2+yShift-4, maxX/2+xShift-5, maxY/2+yShift-2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Нажмите Enter")
	}
	_, err := setCurrentViewOnTop(g, "error_button")
	if err != nil {
		return err
	}
	if err := g.SetKeybinding("error_button", gocui.KeyEnter, gocui.ModNone, errUnSpawn); err != nil {
		return err
	}
	return nil
}

func errUnSpawn(g *gocui.Gui, v *gocui.View) error {
	g.DeleteView("error_box")
	g.DeleteView("error_text")
	g.DeleteView("error_button")

	_, err := setCurrentViewOnTop(g, errReturn)
	if err != nil {
		return err
	}
	return nil
}

func errWrapper(f func(*gocui.Gui, *gocui.View) error, retview string) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		err := f(g, v)
		if err != nil {
			return errSpawn(g, err.Error(), retview)
		}
		return nil
	}
}
