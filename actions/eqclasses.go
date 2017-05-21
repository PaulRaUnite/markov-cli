package actions

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

var EqClasses = Page{
	"B",
	2,
	"Сообщающиеся классы эквивалентности",
	"comclasses",
	[]Element{
		{"calculate", "", "calculate", false, 1, errWrapper(comclassesCalculate, "calculate")},
	},
	0,
	retView,
	retDelete,
}

func comclassesCalculate(g *gocui.Gui, _ *gocui.View) error {
	if !chainReady {
		return errChainNotLoaded
	}
	classes, err := chain.EqualityClasses()
	if err != nil {
		return err
	}
	v, err := g.View("result")
	if err != nil {
		return err
	}
	v.Clear()
	for _, class := range classes {
		ok, err := chain.CommunicatingClass(class)
		if err != nil {
			return err
		}
		if ok {
			fmt.Fprintln(v, class)
		}
	}
	return nil
}
