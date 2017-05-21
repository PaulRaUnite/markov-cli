package main

import (
	"flag"
	"fmt"
	"log"
	"math"

	"github.com/PaulRaUnite/mark-lib"
	"github.com/PaulRaUnite/markov-cli/actions"
	"github.com/jroimartin/gocui"
)

var (
	acts []actions.Page = []actions.Page{
		actions.File,
		actions.ITOJByM,
		actions.DistByM,
		actions.ExpValUn,
		actions.ExpValCon,
		actions.AttFromIToJ,
		actions.AttSet,
		actions.Ergodic,
		actions.Communicate,
		actions.EqClasses,
		actions.Absorbing,
	}
	previous int = 0
	now      int = 0
)

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if now > 0 {
		now--
	}
	_, err := setCurrentViewOnTop(g, acts[now].ViewName)
	return err
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if now < len(acts)-1 {
		now++
	}
	_, err := setCurrentViewOnTop(g, acts[now].ViewName)
	return err
}

func selectAction(g *gocui.Gui, v *gocui.View) error {
	prev := acts[previous]
	prev.UnSpawn(g)
	previous = now
	actions.ReturnTo = acts[now].ViewName
	return acts[now].Spawn(g)
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	for _, value := range acts {
		if err := g.SetKeybinding(value.ViewName, gocui.KeyEnter, gocui.ModNone, selectAction); err != nil {
			return err
		}
		if err := g.SetKeybinding(value.ViewName, gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
			return err
		}
		if err := g.SetKeybinding(value.ViewName, gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
			return err
		}
	}
	return nil
}
func layout(g *gocui.Gui) error {
	prev := 0
	for _, value := range acts {
		next := prev + value.BodySize + 1
		if v, err := g.SetView(value.ViewName, 0, prev, actions.Shift, next); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			v.Title = value.Title
			v.Wrap = true
			if value.ViewName == "loadfile" {
				if _, err = setCurrentViewOnTop(g, "loadfile"); err != nil {
					return err
				}
			}
			fmt.Fprintln(v, value.Body)
		}
		prev = next
	}
	/*if v, err := g.SetView("list", 0, 0, actions.Shift, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "List"
	}*/
	return nil
}

func main() {
	tolerance := flag.Int("tolerance", 8, "digits after comma to compare floats")
	flag.Parse()
	mark_lib.TOLERANCE = math.Pow(10, -float64(*tolerance))
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	g.SetManagerFunc(layout)
	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	log.Println("ha?!", actions.Buf)
}

/*
func layoutList(g *gocui.Gui) error {
	_, maxY := g.Size()
	if v, err := g.SetView("list", 0, 0, 8, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "List"
		v.Highlight = true
		v.Editable  = true
		v.Wrap = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "вероятность перехода")
		fmt.Fprintln(v, "B")
		fmt.Fprintln(v, "нужно добавить новое")
		fmt.Fprintln(v, "Load")
	}
	switch tab {
	case "A":
		layoutA(g)
	case "B":
		layoutB(g)
	case "C":
		layoutC(g)
	case "Load":
		layoutLoad(g)
	}
	return nil
}
func chooseView(g *gocui.Gui, v *gocui.View) error {
	_, err := setCurrentViewOnTop(g, v.Name())
	return err
}

func layoutA(g *gocui.Gui) error {
	maxX, _ := g.Size()
	if v, err := g.SetView("a1", 9, 0, maxX-1, 4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "1"
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}
	shift := 10
	cellSize := (maxX-shift)/4
	if v, err := g.SetView("a1i", shift, 1, shift + cellSize-1, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "i")
		v.Title = "i"
		v.Highlight = true
		v.Editable  = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}
	if err := g.SetKeybinding("a1i", gocui.MouseLeft, gocui.ModNone, chooseView); err != nil {
			return err
	}
	if v, err := g.SetView("a1j", shift + cellSize, 1, shift + cellSize*2-1, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "j"
		v.Highlight = true
		v.Editable  = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}
	if err := g.SetKeybinding("a1j", gocui.MouseLeft, gocui.ModNone, chooseView); err != nil {
		return err
	}
	if v, err := g.SetView("a1m", shift + cellSize * 2, 1, shift + cellSize*3-1, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "m"
		v.Highlight = true
		v.Editable  = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}
	if err := g.SetKeybinding("a1m", gocui.MouseLeft, gocui.ModNone, chooseView); err != nil {
		return err
	}
	if v, err := g.SetView("a1calculate", shift + cellSize *3, 1, maxX - 2, 3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v,"calculate")
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}
	//////
	if v, err := g.SetView("a2", shift-1, 6, maxX-1, 10); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "2"
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}
	cellSize = (maxX-shift)*3/4
	if v, err := g.SetView("a2distribution", shift, 7, shift + cellSize - 1, 9); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "distribution"
		v.Highlight = true
		v.Editable  = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}
	if err := g.SetKeybinding("a2distribution", gocui.MouseLeft, gocui.ModNone, chooseView); err != nil {
		return err
	}
	if v, err := g.SetView("a2calculate", shift + cellSize, 7, maxX - 2, 9); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "calculate")
		v.Highlight = true
		v.Editable  = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}
	//////
	if v, err := g.SetView("a3", shift-1, 12, maxX-1, 17); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "3"
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}
	cellSize = (maxX-shift)/4
	if v, err := g.SetView("a3m", shift, 13, shift + cellSize - 1, 16); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "m"
		v.Highlight = true
		v.Editable = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}
	if err := g.SetKeybinding("a3m", gocui.MouseLeft, gocui.ModNone, chooseView); err != nil {
		return err
	}
	if v, err := g.SetView("a3list", shift + cellSize, 13, shift + cellSize*2 - 1, 16); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "list"
		v.Highlight = true
		v.Editable = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		fmt.Fprintln(v, "uniform")
		fmt.Fprintln(v, "concentrated")
	}
	if v, err := g.SetView("a3in", shift + cellSize*2, 13, shift+ cellSize*3 - 1, 16); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "element"
		v.Highlight = true
		v.Editable  = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}
	if err := g.SetKeybinding("a3in", gocui.MouseLeft, gocui.ModNone, chooseView); err != nil {
		return err
	}
	if v, err := g.SetView("a3calculate", shift + cellSize*3, 13, maxX - 2, 16); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "calculate")
		v.Highlight = true
		v.Editable  = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
	}
	return nil
}
func layoutB(g *gocui.Gui) error {
	return nil
}

func layoutC(g *gocui.Gui) error {
	return nil
}

func layoutLoad(g *gocui.Gui) error {
	return nil
}


func showMsg(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}

	maxX, maxY := g.Size()
	if v, err := g.SetView("msg", maxX/2-10, maxY/2, maxX/2+10, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, l)
	}
	return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}
	return nil
}
*/
