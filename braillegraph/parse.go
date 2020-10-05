package graph

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell"
)

//=GRAPH(X1:X1,Y1:Y1,X2:X2,Y2:Y2...)
//.WithSize(width,height)
//.WithColours(X1Col,X2Col...)
//.WithTitle("your title")
//.WithXLabel("x-axis label)
//.WithYLabel("y-axis label)              <- this one will be funny because only up and down letters
//.WithMinX(X)
//.WithMaxX(X)
//.WithMinY(Y)
//.WithMaxY(Y)
func NewGraphFromFormula(f string, baseCoor cell.Coor, s zstore.ZStore, axes cell.Axes) (g Graph, err error) {
	formulaAttr := strings.Split(f, ".With")

	// get the data
	ds, err := getDatasetFromFormula(formulaAttr[0], baseCoor, s, axes)
	if err != nil {
		return g, err
	}

	// defaults
	var width, height uint16 = 50, 15

	for i := 1; i < len(formulaAttr); i++ {
		fa := formulaAttr[i]
		switch {
		case strings.HasPrefix(fa, "Size"):
			fa = trimPrefixAndBrackets(fa, "Size")
			width, height, err = getSize(fa)
			if err != nil {
				return g, err
			}

		case strings.HasPrefix(fa, "Colours"):
			fa = trimPrefixAndBrackets(fa, "Colours")
		case strings.HasPrefix(fa, "Title"):
			fa = trimPrefixAndBrackets(fa, "Title")
		case strings.HasPrefix(fa, "XLabel"):
			fa = trimPrefixAndBrackets(fa, "XLabel")
		case strings.HasPrefix(fa, "YLabel"):
			fa = trimPrefixAndBrackets(fa, "YLabel")
		case strings.HasPrefix(fa, "MinX"):
			fa = trimPrefixAndBrackets(fa, "MinX")
		case strings.HasPrefix(fa, "MinY"):
			fa = trimPrefixAndBrackets(fa, "MinY")
		case strings.HasPrefix(fa, "MaxX"):
			fa = trimPrefixAndBrackets(fa, "MaxX")
		case strings.HasPrefix(fa, "MaxY"):
			fa = trimPrefixAndBrackets(fa, "MaxY")
		default:
			return g, fmt.Errorf("unknown graph attribute: %v", fa)
		}
	}

	g = Graph{
		Width:                 width,
		Height:                height,
		AxisAndIntersectionFg: tcell.NewRGBColor(255, 255, 255),
		PlotBg:                tcell.NewRGBColor(0, 0, 0),
		Datasets:              ds,
	}
	g.DrawBuffer()

	return g, nil
}

func trimPrefixAndBrackets(str, prefix string) string {
	str = strings.TrimPrefix(str, prefix+"(")
	str = strings.TrimSuffix(str, ")")
	return str
}

func getDatasetFromFormula(formula string, baseCoor cell.Coor, s zstore.ZStore, axes cell.Axes) (ds GraphDataSets, err error) {
	formula = strings.TrimPrefix(formula, "=GRAPH(")
	formula = strings.TrimSuffix(formula, ")")
	split := strings.Split(formula, ",")
	if len(split) != 2 {
		return ds, errors.New("at least two data ranges (Xs, and Ys) must be provided")
	}
	Xs, err := zspatial.GetRelativeRange(baseCoor, split[0], axes)
	if err != nil {
		return ds, err
	}
	Ys, err := zspatial.GetRelativeRange(baseCoor, split[1], axes)
	if err != nil {
		return ds, err
	}
	ds = GraphDataSets{ // XXX incomplete
		GraphDataSet{
			X: Xs.GetCalculated(s),
			Y: Ys.GetCalculated(s),
		},
	}
	return ds, nil
}

func getSize(contents string) (width, height uint16, err error) {
	split := strings.Split(contents, ",")
	if len(split) != 2 {
		return 0, 0, errors.New("not 2 integers in size attribute")
	}
	w, err := strconv.Atoi(split[0])
	if err != nil {
		return 0, 0, err
	}
	h, err := strconv.Atoi(split[1])
	if err != nil {
		return 0, 0, err
	}
	return uint16(w), uint16(h), nil
}
