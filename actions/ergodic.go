package actions

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

var Ergodic = Page{
	"B",
	1,
	"Существенно ли i",
	"ergodic",
	[]Element{
		{"i", "i", "0", true, 1, nil},
		{"calculate", "", "calculate", false, 1, errWrapper(ergodicCalculate, "calculate")},
	},
	0,
	retView,
	retDelete,
}

func ergodicCalculate(g *gocui.Gui, _ *gocui.View) error {
	if !chainReady {
		return errChainNotLoaded
	}
	i, err := getInt(g, "i")
	if err != nil {
		return err
	}
	ok, err := chain.Ergodic(i)
	if err != nil {
		return err
	}
	v, err := g.View("result")
	if err != nil {
		return err
	}
	v.Clear()
	fmt.Fprintln(v, ok)
	return nil
}
