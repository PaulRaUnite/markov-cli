package actions

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

var AttFromIToJ = Page{
	"B",
	2,
	"Достижимо ли состояние j из i",
	"attfromtoj",
	[]Element{
		{"i", "i", "0", true, 1, nil},
		{"j", "j", "1", true, 1, nil},
		{"calculate", "", "calculate", false, 1, errWrapper(attfromitojCalculate, "calculate")},
	},
	0,
	retView,
	retDelete,
}

func attfromitojCalculate(g *gocui.Gui, _ *gocui.View) error {
	if !chainReady {
		return errChainNotLoaded
	}
	i, err := getInt(g, "i")
	if err != nil {
		return err
	}
	j, err := getInt(g, "j")
	if err != nil {
		return err
	}
	ok, err := chain.Attainability(i, j)
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
