package hex

import (
	"cmp"
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

// Pt is shorthand for Hex{q, r}.
func Pt(q, r int) Hex {
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

// Length returns the distance from the origin (0,0) in hex steps.
func (h Hex) Length() int {
	return (geom.Abs(h.Q) + geom.Abs(h.R) + geom.Abs(h.S())) / 2
}

// DistanceTo returns the hex distance between h and the given hex.
func (h Hex) DistanceTo(hex Hex) int {
	return h.Subtract(hex).Length()
}

// Rotate returns the hex rotated by steps×60° around the origin, in the same sense as
// [Direction.Rotate]: clockwise as drawn on a screen with Y pointing down.
// Negative steps rotate the other way.
func (h Hex) Rotate(steps int) Hex {
	return h.RotateAround(Hex{}, steps)
}

// RotateAround returns the hex rotated by steps×60° around center, in the same sense as
// [Direction.Rotate]: clockwise as drawn on a screen with Y pointing down.
// Negative steps rotate the other way.
func (h Hex) RotateAround(center Hex, steps int) Hex {
	relative := h.Subtract(center)
	steps = geom.Mod(steps, 6)
	for i := 0; i < steps; i++ {
		relative = Hex{-relative.R, relative.Q + relative.R}
	}

	return center.Add(relative)
}

// ReflectQ returns the hex reflected across the q-axis (q unchanged, r and s swapped).
func (h Hex) ReflectQ() Hex {
	return Hex{h.Q, -h.Q - h.R}
}

// ReflectR returns the hex reflected across the r-axis (r unchanged, q and s swapped).
func (h Hex) ReflectR() Hex {
	return Hex{-h.Q - h.R, h.R}
}

// ReflectS returns the hex reflected across the s-axis (s unchanged, q and r swapped).
func (h Hex) ReflectS() Hex {
	return Hex{h.R, h.Q}
}

// Neighbors returns the six neighboring hexes around h in axial coordinates.
func (h Hex) Neighbors() []Hex {
	neighbors := make([]Hex, len(Directions))
	for i, direction := range Directions {
		neighbors[i] = h.Neighbor(direction)
	}

	return neighbors
}

// Neighbor returns the neighboring hex of h in the given direction.
func (h Hex) Neighbor(direction Direction) Hex {
	v := direction.Offset()

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

// Ring returns the hexes at exactly distance radius from h, ordered by increasing angle from
// the [SMinus] corner: clockwise as drawn on a screen with Y pointing down. Ring(1) is exactly
// [Hex.Neighbors].
// Returns nil for radius < 0 and a slice containing only h for radius == 0.
func (h Hex) Ring(radius int) []Hex {
	if radius < 0 {
		return nil
	}
	if radius == 0 {
		return []Hex{h}
	}

	results := make([]Hex, 0, 6*radius)
	for _, direction := range Directions {
		corner := direction.Hex().Multiply(radius)
		edge := direction.Rotate(2).Hex()

		for j := range radius {
			results = append(results, h.Add(corner).Add(edge.Multiply(j)))
		}
	}

	return results
}

// Spiral returns all hexes from h outward to radius, starting with h and expanding
// ring by ring. The result contains the same hexes as Range but in spiral order,
// useful for nearest-first traversal.
func (h Hex) Spiral(radius int) []Hex {
	if radius < 0 {
		return nil
	}

	results := make([]Hex, 0, 1+3*radius*(radius+1))
	results = append(results, h)
	for r := 1; r <= radius; r++ {
		results = append(results, h.Ring(r)...)
	}

	return results
}

// Line returns the sequence of hexes that connects this hex to target in a straight line.
func (h Hex) Line(target Hex) []Hex {
	n := h.DistanceTo(target)
	step := 1.0 / math.Max(float64(n), 1.0)

	// nudge the line in one direction to avoid landing on hex side boundaries.
	e := 1e-6
	start := FractionalHex{float64(h.Q) + e, float64(h.R) + 2*e}
	end := FractionalHex{float64(target.Q) + e, float64(target.R) + 2*e}

	results := make([]Hex, 0, n+1)
	for i := 0; i <= n; i++ {
		results = append(results, start.Lerp(end, step*float64(i)).Round())
	}

	return results
}

// HasLineOfSight checks if the target hex is visible from this hex, taking into account a set of blocking hexagons.
// Neither h nor target counts as a blocker — only the hexes strictly between them do.
func (h Hex) HasLineOfSight(target Hex, blocking []Hex) bool {
	line := h.Line(target)
	if len(line) < 3 {
		return true
	}

	for _, lineHex := range line[1 : len(line)-1] {
		if slices.Contains(blocking, lineHex) {
			return false
		}
	}

	return true
}

// FieldOfView returns the subset of candidate hexes visible from this hex, taking into account a set of blocking hexagons.
func (h Hex) FieldOfView(candidates []Hex, blocking []Hex) []Hex {
	results := make([]Hex, 0, len(candidates))
	for _, candidate := range candidates {
		if len(blocking) == 0 || h.DistanceTo(candidate) <= 1 || h.HasLineOfSight(candidate, blocking) {
			results = append(results, candidate)
		}
	}

	return results
}

// IsZero reports whether the hex is at the origin (0, 0).
func (h Hex) IsZero() bool {
	return h.Q == 0 && h.R == 0
}

// Compare returns -1, 0, or +1 as h sorts before, with, or after hex, ordering by
// Q and then by R, the axial counterpart of [geom.Point.Compare]. It follows the
// [cmp.Compare] convention.
func (h Hex) Compare(hex Hex) int {
	if c := cmp.Compare(h.Q, hex.Q); c != 0 {
		return c
	}

	return cmp.Compare(h.R, hex.R)
}

// To converts the hex into the specified coordinate system, returning an ints.Point.
func (h Hex) To(system CoordinateSystem) ints.Point {
	return system.To(h)
}

// Point returns (q,r) as a [ints.Point].
func (h Hex) Point() ints.Point {
	return geom.Pt(h.Q, h.R)
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

// FracPt is shorthand for FractionalHex{q, r}.
func FracPt(q, r float64) FractionalHex {
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

// Point returns the axial (q,r) as a floats.Point.
func (h FractionalHex) Point() floats.Point {
	return geom.Pt(h.Q, h.R)
}

// String returns a compact representation of the fractional hex as (q,r) with 2 decimals.
func (h FractionalHex) String() string {
	return fmt.Sprintf("(%.2f,%.2f)", h.Q, h.R)
}
