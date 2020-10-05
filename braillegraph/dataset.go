package graph

import (
	"errors"
)

type GraphDataSets []GraphDataSet

func (gs GraphDataSets) GetMinMaxXY() (minX, minY, maxX, maxY float64) {
	if len(gs) == 0 {
		return
	}
	minX, minY, maxX, maxY = gs[0].GetMinMaxXY()
	for i := 1; i < len(gs); i++ {
		g := gs[i]
		localMinX, localMinY, localMaxX, localMaxY := g.GetMinMaxXY()
		if localMinX < minX {
			minX = localMinX
		}
		if localMinY < minY {
			minY = localMinY
		}
		if localMaxX > maxX {
			maxX = localMaxX
		}
		if localMaxY > maxY {
			maxY = localMaxY
		}
	}
	return minX, minY, maxX, maxY
}

type GraphDataSet struct {
	Name string
	Fg   uint8 // point colour
	X    []float64
	Y    []float64
}

// NewGraphDataSet creates a new GraphDataSet object
func NewGraphDataSet(name string, fg uint8, x []float64, y []float64) (GraphDataSet, error) {
	if len(x) != len(y) {
		return GraphDataSet{}, errors.New("uneven data points")
	}
	if len(x) == 0 || len(y) == 0 {
		return GraphDataSet{}, errors.New("no x or y data")
	}

	return GraphDataSet{
		Name: name,
		Fg:   fg,
		X:    x,
		Y:    y,
	}, nil
}

func (g GraphDataSet) Iterate(fn func(x, y float64)) {
	minLen := len(g.Y)
	if len(g.X) < len(g.Y) { // this can happen as datasets are in the middle of being changed
		minLen = len(g.X)
	}
	for i := 0; i < minLen; i++ {
		fn(g.X[i], g.Y[i])
	}
}

// Get the MinX of datapoints
func (g GraphDataSet) GetMinMaxXY() (minX, minY, maxX, maxY float64) {
	minX, minY = g.X[0], g.Y[0]
	maxX, maxY = g.X[0], g.Y[0]
	for _, y := range g.Y {
		for _, x := range g.X {
			if x < minX {
				minX = x
			}
			if y < minY {
				minY = y
			}
			if x > maxX {
				maxX = x
			}
			if y > maxY {
				maxY = y
			}
		}
	}
	return minX, minY, maxX, maxY
}
