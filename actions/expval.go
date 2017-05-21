package actions

import (
	"errors"
	"fmt"

	"github.com/jroimartin/gocui"
)

var ExpValUn = Page{
	"A",
	3,
	"Математическое ожидание(равномерное) на m шаге",
	"expcalum",
	[]Element{
		{"m", "m", "2", true, 1, nil},
		{"calculate", "", "calculate", false, 1, errWrapper(expvalumCalculate, "calculate")},
	},
	0,
	retView,
	retDelete,
}

func expvalumCalculate(g *gocui.Gui, _ *gocui.View) error {
	if !chainReady {
		return errChainNotLoaded
	}
	m, err := getInt(g, "m")
	if err != nil {
		return err
	}
	n := chain.Size()
	un := float64(1)/float64(n)
	slice := make([]float64, n)
	for i := 0; i < n; i++ {
		slice[i] = un
	}
	dist, err := chain.ExpectedValue(m, slice)
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

var ExpValCon = Page{
	"A",
	3,
	"Математическое ожидание(сконцентрированное в i) на m шаге",
	"expcalcon",
	[]Element{
		{"m", "m", "2", true, 1, nil},
		{"i", "i", "0", true, 1, nil},
		{"calculate", "", "calculate", false, 1, errWrapper(expvalconCalculate, "calculate")},
	},
	0,
	retView,
	retDelete,
}

var errInvalidI error = errors.New("invalid i value")
func expvalconCalculate(g *gocui.Gui, _ *gocui.View) error {
	if !chainReady {
		return errChainNotLoaded
	}
	m, err := getInt(g, "m")
	if err != nil {
		return err
	}
	i, err := getInt(g, "i")
	if err != nil {
		return err
	}
	n := chain.Size()
	slice := make([]float64, n)
	if i < 0 || i >= n {
		return errInvalidI
	}
	slice[i] = 1
	dist, err := chain.ExpectedValue(int(m), slice)
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
