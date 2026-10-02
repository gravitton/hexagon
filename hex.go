package hex

import (
	"cmp"
	"fmt"
	"slices"

	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/ints"
)

// Hex represents a hexagon in axial (cube) coordinates using integer q and r.
// The third coordinate s is implied by s = -q - r.
type Hex struct {
	Q int `json:"q"`
	R int `json:"r"`
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

// Length returns the distance from the origin (0,0) in hex steps.
func (h Hex) Length() int {
	return (geom.Abs(h.Q) + geom.Abs(h.R) + geom.Abs(h.S())) / 2
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

// Lerp creates a new Hex at the interpolated position between h and hex, rounded to the
// nearest hex. It interpolates through [FractionalHex], so the result is a hex the straight
// line from h to hex passes through. It panics for a NaN or infinite t, as
// [FractionalHex.Round] does.
func (h Hex) Lerp(hex Hex, t float64) Hex {
	return h.Float().Lerp(hex.Float(), t).Round()
}

// Turn returns the hex rotated by steps×60° around the origin, in the same sense as
// [Direction.Turn]: clockwise as drawn on a screen with Y pointing down.
// Negative steps rotate the other way.
func (h Hex) Turn(steps int) Hex {
	return h.TurnAround(Hex{}, steps)
}

// TurnAround returns the hex rotated by steps×60° around center, in the same sense as
// [Direction.Turn]: clockwise as drawn on a screen with Y pointing down.
// Negative steps rotate the other way.
func (h Hex) TurnAround(center Hex, steps int) Hex {
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

// Neighbor returns the neighboring hex of h in the given direction, and h itself for
// [DirectionNone], which steps nowhere.
func (h Hex) Neighbor(direction Direction) Hex {
	return h.Add(direction.Hex())
}

// Neighbors returns the six neighboring hexes around h in axial coordinates, ordered by
// increasing angle like [Directions]. It returns an array, so it allocates nothing.
func (h Hex) Neighbors() [6]Hex {
	var neighbors [6]Hex
	for i, direction := range Directions() {
		neighbors[i] = h.Neighbor(direction)
	}

	return neighbors
}

// Range returns the set of hexes within radius n around h, inclusive of h, ordered by q and
// then by r like [Hex.Compare]. It returns nil for a negative radius and h alone for zero.
func (h Hex) Range(n int) []Hex {
	if n < 0 {
		return nil
	}

	return h.AppendRange(make([]Hex, 0, areaOf(n)), n)
}

// AppendRange appends the hexes Range returns to dst and returns the extended slice, so a
// caller reusing dst allocates nothing once it has room. A negative radius appends nothing.
func (h Hex) AppendRange(dst []Hex, n int) []Hex {
	for q := -n; q <= n; q++ {
		for r := max(-n, -q-n); r <= min(n, -q+n); r++ {
			dst = append(dst, Hex{h.Q + q, h.R + r})
		}
	}

	return dst
}

// Ring returns the hexes at exactly distance radius from h, ordered by increasing angle from
// the [SMinus] corner: clockwise as drawn on a screen with Y pointing down, so the ring of
// radius one holds exactly [Hex.Neighbors]. It returns nil for a negative radius and h alone
// for zero.
func (h Hex) Ring(radius int) []Hex {
	if radius < 0 {
		return nil
	}

	return h.AppendRing(make([]Hex, 0, max(6*radius, 1)), radius)
}

// AppendRing appends the hexes Ring returns to dst and returns the extended slice, so a
// caller reusing dst allocates nothing once it has room. A negative radius appends nothing.
func (h Hex) AppendRing(dst []Hex, radius int) []Hex {
	if radius < 0 {
		return dst
	}
	if radius == 0 {
		return append(dst, h)
	}

	for _, direction := range Directions() {
		corner := h.Add(direction.Hex().Multiply(radius))
		edge := direction.Turn(2).Hex()

		for j := range radius {
			dst = append(dst, corner.Add(edge.Multiply(j)))
		}
	}

	return dst
}

// Spiral returns all hexes from h outward to radius, starting with h and expanding ring by
// ring. The result contains the same hexes as Range but in spiral order, useful for
// nearest-first traversal. It returns nil for a negative radius.
func (h Hex) Spiral(radius int) []Hex {
	if radius < 0 {
		return nil
	}

	return h.AppendSpiral(make([]Hex, 0, areaOf(radius)), radius)
}

// AppendSpiral appends the hexes Spiral returns to dst and returns the extended slice, so a
// caller reusing dst allocates nothing once it has room. A negative radius appends nothing.
func (h Hex) AppendSpiral(dst []Hex, radius int) []Hex {
	for r := range radius + 1 {
		dst = h.AppendRing(dst, r)
	}

	return dst
}

// DistanceTo returns the hex distance between h and the given hex.
func (h Hex) DistanceTo(hex Hex) int {
	return h.Subtract(hex).Length()
}

// Line returns the sequence of hexes that connects h to target in a straight line, both ends
// included. The line is nudged off the hex boundaries, so it never lands on a tie, except far
// enough from the origin that float64 can no longer resolve the nudge.
func (h Hex) Line(target Hex) []Hex {
	return h.AppendLine(make([]Hex, 0, h.DistanceTo(target)+1), target)
}

// AppendLine appends the hexes Line returns to dst and returns the extended slice, so a caller
// reusing dst allocates nothing once it has room.
func (h Hex) AppendLine(dst []Hex, target Hex) []Hex {
	n := h.DistanceTo(target)
	for i := 0; i <= n; i++ {
		dst = append(dst, h.lineAt(target, i, n))
	}

	return dst
}

// HasLineOfSight reports whether target is visible from h past the blocking hexes. Neither h
// nor target counts as a blocker, only the hexes of the line strictly between them do. It walks
// the line without building it, so it allocates nothing.
func (h Hex) HasLineOfSight(target Hex, blocking []Hex) bool {
	if len(blocking) == 0 {
		return true
	}

	n := h.DistanceTo(target)
	for i := 1; i < n; i++ {
		if slices.Contains(blocking, h.lineAt(target, i, n)) {
			return false
		}
	}

	return true
}

// FieldOfView returns the candidates visible from h past the blocking hexes, in the order
// given. A candidate at distance one or less is always visible.
func (h Hex) FieldOfView(candidates []Hex, blocking []Hex) []Hex {
	return h.AppendFieldOfView(make([]Hex, 0, len(candidates)), candidates, blocking)
}

// AppendFieldOfView appends the hexes FieldOfView returns to dst and returns the extended
// slice, so a caller reusing dst allocates nothing once it has room.
func (h Hex) AppendFieldOfView(dst []Hex, candidates []Hex, blocking []Hex) []Hex {
	for _, candidate := range candidates {
		if h.HasLineOfSight(candidate, blocking) {
			dst = append(dst, candidate)
		}
	}

	return dst
}

// lineAt returns the hex at step i of the n steps of the line from h to target, both nudged
// by the same amount so the line never runs along a hex boundary.
func (h Hex) lineAt(target Hex, i, n int) Hex {
	const nudge = 1e-6

	start := FractionalHex{float64(h.Q) + nudge, float64(h.R) + 2*nudge}
	end := FractionalHex{float64(target.Q) + nudge, float64(target.R) + 2*nudge}

	return start.Lerp(end, float64(i)/float64(max(n, 1))).Round()
}

// Equal reports whether h and hex are the same hex. Axial coordinates are integers, so the
// comparison is exact and Hex is usable as a map key with == just as well.
func (h Hex) Equal(hex Hex) bool {
	return h == hex
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

// IsZero reports whether the hex is at the origin (0, 0).
func (h Hex) IsZero() bool {
	return h == Hex{}
}

// To converts the hex into the specified coordinate system, returning an ints.Point.
func (h Hex) To(system CoordinateSystem) ints.Point {
	return system.To(h)
}

// Point returns (q,r) as an [ints.Point].
func (h Hex) Point() ints.Point {
	return geom.Pt(h.Q, h.R)
}

// Float converts the hex to a [FractionalHex], the counterpart of [FractionalHex.Round].
func (h Hex) Float() FractionalHex {
	return FractionalHex{float64(h.Q), float64(h.R)}
}

// String returns a compact representation of the hex as (q,r).
func (h Hex) String() string {
	return fmt.Sprintf("(%d,%d)", h.Q, h.R)
}

// areaOf returns the number of hexes within radius n of a hex, the length of Range and Spiral.
func areaOf(n int) int {
	return 1 + 3*n*(n+1)
}
