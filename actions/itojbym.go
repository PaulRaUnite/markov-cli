package actions

import (
	"fmt"
	"github.com/jroimartin/gocui"
)

var ITOJByM = Page{
	"A",
	2,
	"Вероятность перехода из i в j на m шаге",
	"itojbym",
	[]Element{
		{"i", "i", "0", true, 1, nil},
		{"j", "j", "1", true, 1, nil},
		{"m", "m", "2", true, 1, nil},
		{"calculate", "", "calculate", false, 1, errWrapper(itojbymCalculate, "calculate")},
	},
	0,
	retView,
	retDelete,
}

func itojbymCalculate(g *gocui.Gui, _ *gocui.View) error {
	if !chainReady {
		return errSpawn(g, "Марковская цепь не загружена", "itojbym")
	}
	i, err := getInt(g, "i")
	if err != nil {
		return err
	}
	j, err := getInt(g, "j")
	if err != nil {
		return err
	}
	m, err := getInt(g, "m")
	if err != nil {
		return err
	}
	result, err := chain.Probability(m, i, j)
	if err != nil {
		return err
	}
	v, err := g.View("result")
	if err != nil {
		return err
	}
	v.Clear()
	fmt.Fprintln(v, result)
	return nil
}
