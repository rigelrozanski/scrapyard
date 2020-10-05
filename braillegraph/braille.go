package graph

import "errors"

// Braille point layout
// (0,0) (1,0)
// (0,1) (1,1)
// (0,2) (1,2)
// (0,3) (1,3)
type Braille [2][4]bool

// Rune maps each point in braille to a dot identifier and
// calculates the corresponding unicode symbol.
func (b Braille) Rune() rune {
	braillePts := [8]bool{
		b[0][0], b[0][1], b[0][2],
		b[1][0], b[1][1], b[1][2],
		b[0][3], b[1][3]}
	a := 0
	for i, pt := range braillePts {
		if pt {
			a += 1 << uint(i)
		}
	}
	return rune(a) + '\u2800'
}

func ConvertToBrailleRune(input [][]bool) (out [][]rune, err error) {

	if len(input) == 0 {
		return out, errors.New("no braille input")
	}

	// initialize the 2D array of all brailles
	lenXs := len(input)
	lenYs := len(input[0])
	brWidth := lenXs / 2
	if lenXs%2 != 0 {
		brWidth++
	}
	brHeight := lenYs / 4
	if lenYs%4 != 0 {
		brHeight++
	}
	br := make([][]Braille, brWidth)
	for i := range br {
		br[i] = make([]Braille, brHeight)
	}

	// Flip the Ys of the input
	inputFlipped := make([][]bool, lenXs)
	for i := range inputFlipped {
		inputFlipped[i] = make([]bool, lenYs)
	}
	for x, ys := range input {
		for y, pt := range ys {
			inputFlipped[x][lenYs-1-y] = pt
		}
	}

	// add all points to the brailles
	for x, ys := range inputFlipped {
		for y, pt := range ys {
			br[x/2][y/4][x%2][y%4] = pt
		}
	}

	// convert all the brailles to runes
	out = make([][]rune, brWidth)
	for i := range out {
		out[i] = make([]rune, brHeight)
	}

	// add all points to the brailles
	for x, ys := range br {
		for y, brpt := range ys {
			out[x][y] = brpt.Rune()
		}
	}

	return out, nil
}
