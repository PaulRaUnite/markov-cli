package actions

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

var AttSet = Page{
	"B",
	2,
	"Множество достижимости для i",
	"attset",
	[]Element{
		{"i", "i", "0", true, 1, nil},
		{"calculate", "", "calculate", false, 1, errWrapper(attsetCalculate, "calculate")},
	},
	0,
	retView,
	retDelete,
}

func attsetCalculate(g *gocui.Gui, _ *gocui.View) error {
	if !chainReady {
		return errChainNotLoaded
	}
	i, err := getInt(g, "i")
	if err != nil {
		return err
	}
	set, err := chain.AttainabilitySet(i)
	if err != nil {
		return err
	}
	v, err := g.View("result")
	if err != nil {
		return err
	}
	v.Clear()
	var setInt []int
	for i := range set {
		setInt = append(setInt, i)
	}
	fmt.Fprintln(v, setInt)
	return nil
}
