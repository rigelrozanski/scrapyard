package graph

import (
	"math"

	"github.com/gdamore/tcell"
)

///////
// title line
// braille lines
// bottom axis line
// bottom axis values
// bottom axis labels
///////
// left axis label
// left axis line
// left axis braille
// legend ?

type Graph struct {
	Title          string
	XAxis          GraphAxis // positions: 0=center,1=right // TODO make consts
	YAxis          GraphAxis // positions: 0=center-vertical,1=top,2=center-horizontal // TODO make consts
	HasLegend      bool      // legend only works if first data points are text
	LegendLocation uint8     // 0=bottom,1=right,2=top
	Width          uint16    // graph width in runes
	Height         uint16    // graph height in runes

	// colours
	AxisAndIntersectionFg tcell.Color // intersection datapoints will also be this color
	TextFg                tcell.Color
	PlotBg                tcell.Color
	NonPlotBg             tcell.Color
	AxisBg                tcell.Color
	TitleBg               tcell.Color

	// Data
	Datasets GraphDataSets
	Buffer   [][]rune
	BufferFg [][]uint16 // TODO
	BufferBg [][]uint16
}

type GraphAxis struct {
	showAxis    bool
	label       string
	labelPos    uint8 // 0=center,1=right // TODO make consts
	Prec        uint8 // precision digits
	HasMinValue bool  // TODO utilize
	HasMaxValue bool
	MinValue    float64
	MaxValue    float64
}

// Graph
func (g *Graph) DrawBuffer() {

	// calculate the min and max for y and x based on all data points
	minX, minY, maxX, maxY := g.Datasets.GetMinMaxXY()
	// TODO give a 5% leaway on all sides?

	// calculate plot area

	ptsWidth := g.Width * 2
	ptsHeight := g.Height * 4
	xAxisUnit := (maxX - minX) / float64(ptsWidth-2) // the -2 is a correction factor for possible rounding over the max
	yAxisUnit := (maxY - minY) / float64(ptsHeight-2)
	if xAxisUnit == 0 || yAxisUnit == 0 {
		return // TODO return some error in the formula output
	}

	// TODO always round left, always round up on braille edge

	for _, ds := range g.Datasets {

		pts := make([][]bool, ptsWidth)
		for i := range pts {
			pts[i] = make([]bool, ptsHeight)
		}

		fn := func(x, y float64) {
			xInt := int(math.Round((x - minX) / xAxisUnit))
			yInt := int(math.Round((y - minY) / yAxisUnit))
			if xInt >= len(pts) {
				return // TODO Error in formula output
				//panic(fmt.Sprintf("debug xInt: %v, len(pts):%v\n", xInt, len(pts)))
			}
			if yInt >= len(pts[0]) {
				return // TODO Error in formula output
				//panic(fmt.Sprintf("debug yInt: %v, len(pts):%v\n", yInt, len(pts[0])))
			}
			//panic(fmt.Sprintf("debug xInt, yInt, y, minY, yAxisUnit: %v\n%v\n%v\n%v\n%v\n", xInt, yInt, y, minY, yAxisUnit))
			pts[xInt][yInt] = true
		}
		ds.Iterate(fn)

		g.Buffer, _ = ConvertToBrailleRune(pts) // XXX only considers final dataset
	}
}

func (g Graph) IterateRunes(fn func(x, y int, r rune, style tcell.Style)) {

	for x, ys := range g.Buffer {
		for y, r := range ys {

			// TODO replace proper colours
			style := tcell.StyleDefault.Foreground(g.AxisAndIntersectionFg).Background(g.PlotBg)
			fn(x, y, r, style)
		}
	}
}
