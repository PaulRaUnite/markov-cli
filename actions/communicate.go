package actions

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

var Communicate = Page{
	"B",
	2,
	"Сообщающиеся ли i и j",
	"communicate",
	[]Element{
		{"i", "i", "0", true, 1, nil},
		{"j", "j", "0", true, 1, nil},
		{"calculate", "", "calculate", false, 1, errWrapper(communicateCalculate, "calculate")},
	},
	0,
	retView,
	retDelete,
}

func communicateCalculate(g *gocui.Gui, _ *gocui.View) error {
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
	ok, err := chain.Communicate(i, j)
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
