package actions

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/PaulRaUnite/mark-lib"
	"github.com/jroimartin/gocui"
	"github.com/sqweek/dialog"
)

var File Page = Page{
	"",
	1,
	"Загрузить файл",
	"loadfile",
	[]Element{{"load", "", "Загрузить", false, 1, errWrapper(loadFile, "load")}},
	0,
	retView,
	retDelete,
}

var (
	filename   string
	Buf        []byte
	chain      mark_lib.Chain
	chainReady = false
	re         = regexp.MustCompile(`\r?\n`)
)

func loadFile(g *gocui.Gui, _ *gocui.View) error {
	filename, _ = dialog.File().Load()
	if filename == "" {
		return nil
	}
	var err error
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	Buf, err = ioutil.ReadAll(file)
	v, err := g.View("result")
	if err != nil {
		return err
	}
	v.Clear()
	fmt.Fprintln(v, filename)
	fmt.Fprintln(v, re.ReplaceAllString(string(Buf), "\n"))
	matrix, err := ProcessFile(filename)
	if err != nil {
		return err
	}
	chain, err = mark_lib.NewChain(matrix)
	if err != nil {
		return err
	}

	fmt.Fprintln(v, "")
	fmt.Fprintln(v, chain)
	chainReady = true
	return nil
}

var errInvalidQuantity = errors.New("wrong quantity of transitions in file")

func ProcessFile(filename string) ([][]float64, error) {
	//open file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	//get lines and split by " "
	var strMatrix [][]string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		strMatrix = append(strMatrix, strings.Fields(scanner.Text()))
	}
	//if something goes wrong
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	//Buf = append([]byte(fmt.Sprintln(strMatrix)))
	//convert string into float64
	var matrix [][]float64
	for _, subarr := range strMatrix {
		temp := make([]float64, len(strMatrix))
		if len(subarr)%2 != 0 {
			return nil, errInvalidQuantity
		}
		//main matrix and products
		for i := 0; i < len(subarr)/2; i++ {
			to, err := strconv.ParseInt(subarr[i*2], 10, 64)
			if err != nil {
				return nil, err
			}
			val, err := strconv.ParseFloat(strings.Replace(subarr[i*2+1], ",", ".", 1), 64)
			temp[to] = val
		}
		matrix = append(matrix, temp)
	}
	return matrix, nil
}
