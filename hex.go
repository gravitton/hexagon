package hex

import (
	"fmt"
	"math"
	"slices"

	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/floats"
	"github.com/gravitton/geometry/types/ints"
)

// Hex represents a hexagon in axial (cube) coordinates using integer q and r.
// The third coordinate s is implied by s = -q - r.
type Hex struct {
	Q, R int
}

// H is shorthand for Hex{q, r}.
func H(q, r int) Hex {
	return Hex{q, r}
}

// S returns the implied s coordinate (-q - r).
func (h Hex) S() int {
	return -h.Q - h.R
}

// QR returns the (q, r) coordinates.
func (h Hex) QR() (int, int) {
	return h.Q, h.R
}

// QRS returns the (q, r, s) coordinates where s is implied.
func (h Hex) QRS() (int, int, int) {
	return h.Q, h.R, h.S()
}

// Add returns a new Hex that is the vector sum of h and hex.
func (h Hex) Add(hex Hex) Hex {
	return Hex{h.Q + hex.Q, h.R + hex.R}
}

// Subtract creates a new Hex that is the vector difference h - hex.
func (h Hex) Subtract(hex Hex) Hex {
	return Hex{h.Q - hex.Q, h.R - hex.R}
}

// Multiply creates a new Hex scaled by the given integer factor.
func (h Hex) Multiply(factor int) Hex {
	return Hex{h.Q * factor, h.R * factor}
}

// Length creates the distance from the origin (0,0) in hex steps.
func (h Hex) Length() int {
	return int((math.Abs(float64(h.Q)) + math.Abs(float64(h.R)) + math.Abs(float64(h.S()))) / 2)
}

// DistanceTo returns the hex distance between h and the given hex.
func (h Hex) DistanceTo(hex Hex) int {
	return h.Subtract(hex).Length()
}

// To converts the hex into the specified coordinate system, returning an ints.Point.
func (h Hex) To(system CoordinateSystem) ints.Point {
	return To(h, system)
}

// ToPoint returns (q,r) as an ints.Point.
func (h Hex) ToPoint() ints.Point {
	return geom.P(h.Q, h.R)
}

// Neighbors returns the six neighboring hexes around h in axial coordinates.
func (h Hex) Neighbors() []Hex {
	neighbors := make([]Hex, len(Directions))
	for i, v := range Directions {
		neighbors[i] = Hex{h.Q + v.X, h.R + v.Y}
	}

	return neighbors
}

// Neighbor returns the neighboring hex of h in the given direction.
func (h Hex) Neighbor(direction Direction) Hex {
	v := Directions[direction]

	return Hex{h.Q + v.X, h.R + v.Y}
}

// Range returns the set of hexes within radius n around h, inclusive of h.
// When n < 0, it returns nil, when n == 0, it returns itself.
func (h Hex) Range(n int) []Hex {
	if n < 0 {
		return nil
	}

	results := make([]Hex, 0, 1+n*(n+1)*3)
	for q := -n; q <= n; q++ {
		for r := max(-n, -q-n); r <= min(n, -q+n); r++ {
			results = append(results, Hex{h.Q + q, h.R + r})
		}
	}

	return results
}

// Line returns a list of hexes that connect hex to the target hex in a straight line.
func (h Hex) Line(target Hex) []Hex {
	n := h.DistanceTo(target)
	step := 1.0 / math.Max(float64(n), 1.0)

	// nudge the line in one direction to avoid landing on hex side boundaries.
	e := 1e-6
	start := FractionalHex{float64(h.Q) + e, float64(h.R) + 2*e}
	end := FractionalHex{float64(target.Q) + e, float64(target.R) + 2*e}

	results := make([]Hex, 0, n)
	for i := 0; i <= n; i++ {
		results = append(results, start.Lerp(end, step*float64(i)).Round())
	}

	return results
}

// HasLineOfSight checks if the target hex is visible from the hex, taking into consideration a set of blocking hexagons
func (h Hex) HasLineOfSight(target Hex, blocking []Hex) bool {
	line := h.Line(target)
	for _, lineHex := range line[:len(line)-1] {
		if slices.Contains(blocking, lineHex) {
			return false
		}
	}

	return true
}

// FieldOfView returns a list of hexes that are visible from the hex, taking into consideration a set of blocking hexagons
func (h Hex) FieldOfView(candidates []Hex, blocking []Hex) []Hex {
	results := make([]Hex, 0, len(candidates))
	for _, candidate := range candidates {
		if len(blocking) == 0 || h.DistanceTo(candidate) <= 1 || h.HasLineOfSight(candidate, blocking) {
			results = append(results, candidate)
		}
	}

	return results
}

// String returns a compact representation of the hex as (q,r).
func (h Hex) String() string {
	return fmt.Sprintf("(%d,%d)", h.Q, h.R)
}

// FractionalHex represents a hex with floating-point axial coordinates.
// Useful for interpolation and conversions from pixel space before rounding.
type FractionalHex struct {
	Q, R float64
}

// F is shorthand for FractionalHex{q, r}.
func F(q, r float64) FractionalHex {
	return FractionalHex{q, r}
}

// S returns the implied s coordinate (-q - r).
func (h FractionalHex) S() float64 {
	return -h.Q - h.R
}

// QR returns the (q, r) coordinates.
func (h FractionalHex) QR() (float64, float64) {
	return h.Q, h.R
}

// QRS returns the (q, r, s) coordinates where s is implied.
func (h FractionalHex) QRS() (float64, float64, float64) {
	return h.Q, h.R, h.S()
}

// ToPoint returns the axial (q,r) as a floats.Point.
func (h FractionalHex) ToPoint() floats.Point {
	return geom.P(h.Q, h.R)
}

// Lerp creates a new FractionalHex in linear interpolation towards given hex.
func (h FractionalHex) Lerp(hex FractionalHex, t float64) FractionalHex {
	return FractionalHex{geom.Lerp(h.Q, hex.Q, t), geom.Lerp(h.R, hex.R, t)}
}

// Round converts a FractionalHex to the nearest Hex while preserving q+r+s=0.
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

// String returns a compact representation of the fractional hex as (q,r) with 2 decimals.
func (h FractionalHex) String() string {
	return fmt.Sprintf("(%.2f,%.2f)", h.Q, h.R)
}
