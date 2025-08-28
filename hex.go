package hex

import (
	"fmt"
	"math"

	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/floats"
	"github.com/gravitton/geometry/types/ints"
)

// Axial (and Cube) coordinates
type Hex struct {
	Q, R int
}

// H is shorthand for Hex{q, r}.
func H(q, r int) Hex {
	return Hex{q, r}
}

func (h Hex) S() int {
	return -h.Q - h.R
}

func (h Hex) Add(hex Hex) Hex {
	return Hex{h.Q + hex.Q, h.R + hex.R}
}

func (h Hex) AddQR(deltaQ, deltaR int) Hex {
	return Hex{h.Q + deltaQ, h.R + deltaR}
}

func (h Hex) Subtract(hex Hex) Hex {
	return Hex{h.Q - hex.Q, h.R - hex.R}
}

func (h Hex) SubtractQR(deltaQ, deltaR int) Hex {
	return Hex{h.Q - deltaQ, h.R - deltaR}
}

func (h Hex) Multiply(factor int) Hex {
	return Hex{h.Q * factor, h.R * factor}
}

func (h Hex) Length() int {
	return int((math.Abs(float64(h.Q)) + math.Abs(float64(h.R)) + math.Abs(float64(h.S()))) / 2)
}

func (h Hex) DistanceTo(hex Hex) int {
	return h.Subtract(hex).Length()
}

func (h Hex) ToPoint() ints.Point {
	return geom.P(h.Q, h.R)
}

func (h Hex) To(coordType CoordinateType) ints.Point {
	return HexTo(h, coordType)
}

// Range returns the set of hexagons around the hexagon for a given radius
func (h Hex) Range(n int) []Hex {
	if n <= 0 {
		return nil
	}

	results := make([]Hex, 0, 1+n*(n+1)*3)
	for q := -n; q <= n; q++ {
		for r := max(-n, -q-n); r <= min(n, -q+n); r++ {
			results = append(results, h.AddQR(q, r))
		}
	}

	return results
}

func (h Hex) String() string {
	return fmt.Sprintf("(%+d,%+d)", h.Q, h.R)
}

type FractionalHex struct {
	Q, R float64
}

// F is shorthand for FractionalHex{q, r}.
func F(q, r float64) FractionalHex {
	return FractionalHex{q, r}
}

func (h FractionalHex) S() float64 {
	return -h.Q - h.R
}

func (h FractionalHex) ToPoint() floats.Point {
	return geom.P(h.Q, h.R)
}

func (h FractionalHex) Round() Hex {
	q := math.Round(h.Q)
	r := math.Round(h.R)
	s := math.Round(h.S())

	qDiff := math.Abs(q - h.Q)
	rDiff := math.Abs(r - h.R)
	sDiff := math.Abs(s - h.S())

	if qDiff > rDiff && qDiff > sDiff {
		q = -r - s
	} else if rDiff > sDiff {
		r = -q - s
	}

	return Hex{int(q), int(r)}
}

func (h FractionalHex) String() string {
	return fmt.Sprintf("(%+.2f,%+.2f)", h.Q, h.R)
}
