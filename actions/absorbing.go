package actions

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

var Absorbing = Page{
	"B",
	2,
	"Все поглощающие состояния",
	"absorbing",
	[]Element{
		{"calculate", "", "calculate", false, 1, errWrapper(absorbingCalculate, "calculate")},
	},
	0,
	retView,
	retDelete,
}

func absorbingCalculate(g *gocui.Gui, _ *gocui.View) error {
	if !chainReady {
		return errChainNotLoaded
	}
	abs, err := chain.AbsorbingClasses()
	if err != nil {
		return err
	}
	v, err := g.View("result")
	if err != nil {
		return err
	}
	v.Clear()
	for _, value := range abs {
		fmt.Fprintln(v, []int{value})
	}
	return nil
}
