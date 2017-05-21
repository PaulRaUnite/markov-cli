package actions

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jroimartin/gocui"
)

var ReturnTo string

type Element struct {
	Name     string
	Title    string
	Body     string
	Editable bool
	Part     int
	Action   func(g *gocui.Gui, v *gocui.View) error
}

func Spawn(p *Page, g *gocui.Gui) error {
	if p.Elements == nil || len(p.Elements) == 0 {
		return nil
	}
	maxX, _ := g.Size()
	parts := 0
	for _, value := range p.Elements {
		if value.Part < 0 {
			parts += 1
			continue
		}
		parts += value.Part
	}
	cell := (maxX - Shift) / parts
	prev := Shift + 1
	for i, e := range p.Elements {
		next := prev + cell*e.Part - 2
		if i == len(p.Elements)-1 {
			next = maxX - 1
		}
		if v, err := g.SetView(e.Name, prev, 0, next, 2); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}
			if e.Body != "" {
				fmt.Fprintln(v, e.Body)
			}
			if e.Title != "" {
				v.Title = e.Title
			}
			if e.Editable {
				v.Editable = true
			}
		}
		prev = next + 1
		if err := g.SetKeybinding(e.Name, gocui.KeyArrowUp, gocui.ModNone, p.cursorLeft); err != nil {
			return err
		}
		if err := g.SetKeybinding(e.Name, gocui.KeyArrowDown, gocui.ModNone, p.cursorRight); err != nil {
			return err
		}
		if e.Action != nil {
			if err := g.SetKeybinding(e.Name, gocui.KeyEnter, gocui.ModNone, e.Action); err != nil {
				return err
			}
		}
		if i == p.CurEl {
			_, err := setCurrentViewOnTop(g, e.Name)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func UnSpawn(es []Element, g *gocui.Gui) {
	for _, value := range es {
		g.DeleteView(value.Name)
		g.DeleteKeybindings(value.Name)
	}
}

type Page struct {
	Title    string
	BodySize int
	Body     string
	ViewName string
	Elements []Element
	CurEl    int
	View     func(g *gocui.Gui, p Page) error
	Delete   func(g *gocui.Gui, p Page) error
}

func (p *Page) Spawn(g *gocui.Gui) error {
	err := Spawn(p, g)
	if err != nil {
		return err
	}
	if p.View == nil {
		return nil
	}
	return p.View(g, *p)
}
func (p *Page) UnSpawn(g *gocui.Gui) error {
	UnSpawn(p.Elements, g)
	if p.Delete == nil {
		return nil
	}
	return p.Delete(g, *p)
}

func (p *Page) cursorLeft(g *gocui.Gui, v *gocui.View) error {
	if p.CurEl > 0 {
		p.CurEl--
	} else {
		_, err := setCurrentViewOnTop(g, ReturnTo)
		return err
	}
	_, err := setCurrentViewOnTop(g, p.Elements[p.CurEl].Name)
	return err
}

func (p *Page) cursorRight(g *gocui.Gui, v *gocui.View) error {
	if p.CurEl < len(p.Elements)-1 {
		p.CurEl++
	}
	_, err := setCurrentViewOnTop(g, p.Elements[p.CurEl].Name)
	return err
}

var Shift int = 21

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func getInt(g *gocui.Gui, view string) (int, error) {
	raw, err := g.View(view)
	if err != nil {
		return 0, err
	}
	res, err := strconv.ParseInt(strings.TrimSpace(raw.Buffer()), 10, 64)
	if err != nil {
		return 0, err
	}
	return int(res), nil
}

func getFloat(g *gocui.Gui, view string) (float64, error) {
	raw, err := g.View(view)
	if err != nil {
		return 0, err
	}
	res, err := strconv.ParseFloat(strings.TrimSpace(raw.Buffer()), 64)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func getFloatSlice(g *gocui.Gui, view string) ([]float64, error) {
	raw, err := g.View(view)
	if err != nil {
		return nil, err
	}
	parts := strings.Fields(raw.Buffer())
	var nums []float64
	for _, part := range parts {
		num, err := strconv.ParseFloat(part, 64)
		if err != nil {
			return nil, err
		}
		nums = append(nums, num)
	}
	return nums, nil
}

func retView(g *gocui.Gui, p Page) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("result", Shift+1, 3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Result"
		v.Wrap = true
	}
	return nil
}
func retDelete(g *gocui.Gui, p Page) error {
	g.DeleteView("result")
	return nil
}
