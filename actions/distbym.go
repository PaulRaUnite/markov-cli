package actions

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/kataras/go-errors"
)

var DistByM = Page{
	"A",
	2,
	"Распределение вероятности на m шаге",
	"distbym",
	[]Element{
		{"dist", "dist", "", true, 2, nil},
		{"m", "m", "2", true, 1, nil},
		{"calculate", "", "calculate", false, 1, errWrapper(distbymCalculate, "calculate")},
	},
	0,
	retView,
	retDelete,
}

var (
	errChainNotLoaded error = errors.New("chain is not loaded")
	errInvalidDistLen error = errors.New("distribution lenght is not valid")
	errInvalidDistSum error = errors.New("distribution sum doesn't equal to 1")
)

func distbymCalculate(g *gocui.Gui, _ *gocui.View) error {
	if !chainReady {
		return errChainNotLoaded
	}
	m, err := getInt(g, "m")
	if err != nil {
		return err
	}
	slice, err := getFloatSlice(g, "dist")
	if err != nil {
		return err
	}
	if chain.Size() != len(slice) {
		return errInvalidDistLen
	}
	sum := float64(0)
	for _, value := range slice {
		sum += value
	}
	if sum != 1 {
		return errInvalidDistSum
	}
	dist, err := chain.Distribution(m, slice)
	if err != nil {
		return err
	}
	v, err := g.View("result")
	if err != nil {
		return err
	}
	v.Clear()
	fmt.Fprintln(v, dist)
	return nil
}
