package hex

import (
	"cmp"
	"fmt"
	"iter"
	"slices"
	"strings"

	geom "github.com/gravitton/geometry"
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

// ParseHex parses a hex in the form "(q,r)", the form String prints, each coordinate an integer
// [geom.Parse] accepts. A fractional coordinate is an error, not a rounded hex.
func ParseHex(s string) (Hex, error) {
	q, r, err := parseCoordinates[int](s, "hex")
	if err != nil {
		return Hex{}, err
	}

	return Hex{q, r}, nil
}

// RangeLen returns the number of hexes within radius n of a hex, the length of [Hex.Range] and
// [Hex.Spiral], for sizing the buffer their Append forms fill. It is zero for a negative radius.
func RangeLen(n int) int {
	if n < 0 {
		return 0
	}

	return 1 + 3*n*(n+1)
}

// RingLen returns the number of hexes at exactly distance radius from a hex, the length of
// [Hex.Ring], for sizing the buffer its Append form fills. It is one for a zero radius, the
// center alone, and zero for a negative one.
func RingLen(radius int) int {
	if radius < 0 {
		return 0
	}

	return max(6*radius, 1)
}

// parseCoordinates parses the "(q,r)" form String prints into its two coordinates, naming the
// type of the value in its errors.
func parseCoordinates[T geom.Number](s, name string) (T, T, error) {
	coordinates, opened := strings.CutPrefix(s, "(")
	coordinates, closed := strings.CutSuffix(coordinates, ")")
	before, after, separated := strings.Cut(coordinates, ",")
	if !opened || !closed || !separated {
		return 0, 0, fmt.Errorf("hex: invalid %s format %q", name, s)
	}

	q, err := geom.Parse[T](before)
	if err != nil {
		return 0, 0, fmt.Errorf("hex: invalid q value: %w", err)
	}
	r, err := geom.Parse[T](after)
	if err != nil {
		return 0, 0, fmt.Errorf("hex: invalid r value: %w", err)
	}

	return q, r, nil
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

// Subtract returns a new Hex that is the vector difference h - hex.
func (h Hex) Subtract(hex Hex) Hex {
	return Hex{h.Q - hex.Q, h.R - hex.R}
}

// Multiply returns a new Hex scaled by the given integer factor.
func (h Hex) Multiply(factor int) Hex {
	return Hex{h.Q * factor, h.R * factor}
}

// Lerp returns a new Hex at the interpolated position between h and hex, rounded to the
// nearest hex. It interpolates through [FractionalHex], so the result is a hex the straight
// line from h to hex passes through. It panics for a NaN or infinite t, as
// [FractionalHex.Round] does, and a finite t that carries the position beyond the range of int
// gives a platform-dependent hex.
func (h Hex) Lerp(hex Hex, t float64) Hex {
	return h.Float().Lerp(hex.Float(), t).Round()
}

// Turn returns the hex rotated by steps×60° around the origin, in the same sense as
// [Direction.Turn]: clockwise as drawn on a screen with Y pointing down.
// Negative steps rotate the other way.
func (h Hex) Turn(steps int) Hex {
	switch geom.Mod(steps, 6) {
	case 1:
		return Hex{-h.R, -h.S()}
	case 2:
		return Hex{h.S(), h.Q}
	case 3:
		return Hex{-h.Q, -h.R}
	case 4:
		return Hex{h.R, h.S()}
	case 5:
		return Hex{-h.S(), -h.Q}
	default:
		return h
	}
}

// TurnAround returns the hex rotated by steps×60° around center, in the same sense as
// [Direction.Turn]: clockwise as drawn on a screen with Y pointing down.
// Negative steps rotate the other way.
func (h Hex) TurnAround(center Hex, steps int) Hex {
	return center.Add(h.Subtract(center).Turn(steps))
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

// Neighbors returns the six neighboring hexes around h, ordered by increasing angle like
// [Directions]. It returns an array, so it allocates nothing.
func (h Hex) Neighbors() [6]Hex {
	var neighbors [6]Hex
	for i, direction := range Directions() {
		neighbors[i] = h.Neighbor(direction)
	}

	return neighbors
}

// DiagonalNeighbor returns the hex two steps from h that lies between two neighbors: the one
// reached by a step in the given direction and a step in the direction after it, a twelfth of a
// turn past the direction itself. It returns h itself for [DirectionNone].
func (h Hex) DiagonalNeighbor(direction Direction) Hex {
	return h.Neighbor(direction).Neighbor(direction.Turn(1))
}

// DiagonalNeighbors returns the six diagonal neighbors around h, ordered by increasing angle
// like [Directions], each a twelfth of a turn past the neighbor of the same index. It returns an
// array, so it allocates nothing.
func (h Hex) DiagonalNeighbors() [6]Hex {
	var neighbors [6]Hex
	for i, direction := range Directions() {
		neighbors[i] = h.DiagonalNeighbor(direction)
	}

	return neighbors
}

// Range returns the set of hexes within radius n around h, inclusive of h, ordered by q and
// then by r like [Hex.Compare]. It returns nil for a negative radius and h alone for zero.
func (h Hex) Range(n int) []Hex {
	if n < 0 {
		return nil
	}

	return h.AppendRange(make([]Hex, 0, RangeLen(n)), n)
}

// AppendRange appends the hexes Range returns to dst and returns the extended slice, so a
// caller reusing dst allocates nothing once it has room. A negative radius appends nothing.
func (h Hex) AppendRange(dst []Hex, n int) []Hex {
	for hex := range h.RangeSeq(n) {
		dst = append(dst, hex)
	}

	return dst
}

// RangeSeq returns an iterator over the hexes Range returns, in the same order, so a caller
// walks them without a buffer. A negative radius yields nothing.
func (h Hex) RangeSeq(n int) iter.Seq[Hex] {
	return func(yield func(Hex) bool) {
		bounds := h.rangeBounds(n)
		for q := bounds.minQ; q <= bounds.maxQ; q++ {
			first, last := bounds.row(q)
			for r := first; r <= last; r++ {
				if !yield(Hex{q, r}) {
					return
				}
			}
		}
	}
}

// RangeIntersection returns the hexes within radius n of h that are also within radius m of
// hex, ordered by q and then by r like [Hex.Compare]. It reads them off the cube bounds the
// two ranges share, without building either. It returns nil where the two do not meet, as for
// a negative radius.
func (h Hex) RangeIntersection(n int, hex Hex, m int) []Hex {
	size := h.rangeBounds(n).meet(hex.rangeBounds(m)).size()
	if size == 0 {
		return nil
	}

	return h.AppendRangeIntersection(make([]Hex, 0, size), n, hex, m)
}

// AppendRangeIntersection appends the hexes RangeIntersection returns to dst and returns the
// extended slice, so a caller reusing dst allocates nothing once it has room.
func (h Hex) AppendRangeIntersection(dst []Hex, n int, hex Hex, m int) []Hex {
	for shared := range h.RangeIntersectionSeq(n, hex, m) {
		dst = append(dst, shared)
	}

	return dst
}

// RangeIntersectionSeq returns an iterator over the hexes RangeIntersection returns, in the
// same order, so a caller walks them without a buffer.
func (h Hex) RangeIntersectionSeq(n int, hex Hex, m int) iter.Seq[Hex] {
	return func(yield func(Hex) bool) {
		bounds := h.rangeBounds(n).meet(hex.rangeBounds(m))
		for q := bounds.minQ; q <= bounds.maxQ; q++ {
			first, last := bounds.row(q)
			for r := first; r <= last; r++ {
				if !yield(Hex{q, r}) {
					return
				}
			}
		}
	}
}

// Ring returns the hexes at exactly distance radius from h, ordered by increasing angle from
// the [SMinus] corner: clockwise as drawn on a screen with Y pointing down, so the ring of
// radius one holds exactly [Hex.Neighbors]. It returns nil for a negative radius and h alone
// for zero.
func (h Hex) Ring(radius int) []Hex {
	if radius < 0 {
		return nil
	}

	return h.AppendRing(make([]Hex, 0, RingLen(radius)), radius)
}

// AppendRing appends the hexes Ring returns to dst and returns the extended slice, so a
// caller reusing dst allocates nothing once it has room. A negative radius appends nothing.
func (h Hex) AppendRing(dst []Hex, radius int) []Hex {
	for hex := range h.RingSeq(radius) {
		dst = append(dst, hex)
	}

	return dst
}

// RingSeq returns an iterator over the hexes Ring returns, in the same order, so a caller
// walks them without a buffer. A negative radius yields nothing.
func (h Hex) RingSeq(radius int) iter.Seq[Hex] {
	return func(yield func(Hex) bool) {
		if radius < 0 {
			return
		}
		if radius == 0 {
			yield(h)
			return
		}

		for _, direction := range Directions() {
			corner := h.Add(direction.Hex().Multiply(radius))
			edge := direction.Turn(2).Hex()

			for j := range radius {
				if !yield(corner.Add(edge.Multiply(j))) {
					return
				}
			}
		}
	}
}

// Spiral returns all hexes from h outward to radius, starting with h and expanding ring by
// ring. The result contains the same hexes as Range but in spiral order, useful for
// nearest-first traversal. It returns nil for a negative radius.
func (h Hex) Spiral(radius int) []Hex {
	if radius < 0 {
		return nil
	}

	return h.AppendSpiral(make([]Hex, 0, RangeLen(radius)), radius)
}

// AppendSpiral appends the hexes Spiral returns to dst and returns the extended slice, so a
// caller reusing dst allocates nothing once it has room. A negative radius appends nothing.
func (h Hex) AppendSpiral(dst []Hex, radius int) []Hex {
	for hex := range h.SpiralSeq(radius) {
		dst = append(dst, hex)
	}

	return dst
}

// SpiralSeq returns an iterator over the hexes Spiral returns, in the same order, so a caller
// walks them nearest first without a buffer and stops as soon as it has found what it looks
// for. A negative radius yields nothing.
func (h Hex) SpiralSeq(radius int) iter.Seq[Hex] {
	return func(yield func(Hex) bool) {
		for r := range radius + 1 {
			for hex := range h.RingSeq(r) {
				if !yield(hex) {
					return
				}
			}
		}
	}
}

// rangeBounds returns the bounds of the hexes within radius n of h.
func (h Hex) rangeBounds(n int) cubeBounds {
	return cubeBounds{h.Q - n, h.Q + n, h.R - n, h.R + n, h.S() - n, h.S() + n}
}

// DistanceTo returns the hex distance between h and the given hex.
func (h Hex) DistanceTo(hex Hex) int {
	return h.Subtract(hex).Length()
}

// DirectionTo returns the direction from h towards target: the direction of a neighbor, and for
// a hex further away the direction nearest to it. It compares the cube coordinates, so it is
// exact at any distance. A target exactly between two directions, as a diagonal neighbor is,
// takes the one of lower angle, the direction [Hex.DiagonalNeighbor] names it by. It returns
// [DirectionNone] for h itself.
func (h Hex) DirectionTo(target Hex) Direction {
	delta := target.Subtract(h)
	for _, direction := range Directions() {
		q, r, s := delta.Turn(-int(direction)).QRS()
		if q >= r && r > s {
			return direction
		}
	}

	return DirectionNone
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
	for hex := range h.LineSeq(target) {
		dst = append(dst, hex)
	}

	return dst
}

// LineSeq returns an iterator over the hexes Line returns, in the same order, so a caller walks
// from h towards target without a buffer and stops where it needs to.
func (h Hex) LineSeq(target Hex) iter.Seq[Hex] {
	return func(yield func(Hex) bool) {
		n := h.DistanceTo(target)
		for i := 0; i <= n; i++ {
			if !yield(h.lineAt(target, i, n)) {
				return
			}
		}
	}
}

// HasLineOfSight reports whether target is visible from h past the blocking hexes: whether the
// straight segment between the two centers is clear. A blocking hex the segment runs through
// blocks it, however small the corner it cuts. A hex the segment only touches does not: one
// met at a single corner never blocks, and where the segment runs along the edge between two
// hexes it is blocked only when both of them block. Neither h nor target counts as a blocker.
//
// The test is exact in integers and reads the same from either end, so two hexes always see
// each other or neither does. It is the visibility symmetric shadow casting computes, and it is
// stricter than [Hex.Line]: a hex off the line can clip the segment, and a hex of the line that
// the segment only touches does not block. It allocates nothing.
//
// It costs one pass over the blockers, and one more for each blocker the segment touches along
// an edge, to look for the hex facing it: at most two a step of the distance.
//
// The test multiplies the coordinates pairwise, so it is exact while the square of the distance
// from h to target and to every blocker fits an int a few times over, far beyond any map on a
// 64-bit platform. Further apart the products overflow and the answer means nothing.
func (h Hex) HasLineOfSight(target Hex, blocking []Hex) bool {
	line := h.sightLine(target)

	for _, blocker := range blocking {
		offset := blocker.Subtract(h)

		clearance := line.clearance(offset)
		if clearance < 0 {
			return false
		}
		if clearance == 0 {
			if facing, edge := line.facing(offset); edge && slices.Contains(blocking, h.Add(facing)) {
				return false
			}
		}
	}

	return true
}

// HasLineOfSightFunc reports whether target is visible from h past the hexes blocked reports,
// by the rule of [Hex.HasLineOfSight]: it holds exactly when HasLineOfSight does for the same
// blockers. It asks blocked about the hexes the segment meets, from h towards target, and no
// others, so it costs the distance however many hexes block, where HasLineOfSight costs its
// passes over the blockers: the form for a large or dense map kept in a map or a grid. A nil blocked
// blocks nothing. It allocates nothing.
func (h Hex) HasLineOfSightFunc(target Hex, blocked func(Hex) bool) bool {
	if blocked == nil {
		return true
	}

	line := h.sightLine(target)

	for offset := range line.crossed() {
		clearance := line.clearance(offset)
		if clearance > 0 || !blocked(h.Add(offset)) {
			continue
		}
		if clearance < 0 {
			return false
		}
		if facing, edge := line.facing(offset); edge && blocked(h.Add(facing)) {
			return false
		}
	}

	return true
}

// FieldOfView returns the candidates visible from h past the blocking hexes, in the order
// given: those [Hex.HasLineOfSight] sees, so the field is symmetric and a blocker casts the
// same shadow from every side. A candidate at distance one or less is always visible. It
// returns nil for no candidates.
func (h Hex) FieldOfView(candidates []Hex, blocking []Hex) []Hex {
	if len(candidates) == 0 {
		return nil
	}

	return h.AppendFieldOfView(make([]Hex, 0, len(candidates)), candidates, blocking)
}

// AppendFieldOfView appends the hexes FieldOfView returns to dst and returns the extended
// slice, so a caller reusing dst allocates nothing once it has room.
func (h Hex) AppendFieldOfView(dst []Hex, candidates []Hex, blocking []Hex) []Hex {
	for hex := range h.FieldOfViewSeq(candidates, blocking) {
		dst = append(dst, hex)
	}

	return dst
}

// FieldOfViewSeq returns an iterator over the hexes FieldOfView returns, in the same order, so
// a caller walks the visible candidates without a buffer and tests no candidate past the one
// it stops at.
func (h Hex) FieldOfViewSeq(candidates []Hex, blocking []Hex) iter.Seq[Hex] {
	return func(yield func(Hex) bool) {
		for _, candidate := range candidates {
			if h.HasLineOfSight(candidate, blocking) && !yield(candidate) {
				return
			}
		}
	}
}

// FieldOfViewFunc returns the candidates visible from h past the hexes blocked reports, in the
// order given: those [Hex.HasLineOfSightFunc] sees, the field [Hex.FieldOfView] returns for the
// same blockers. It returns nil for no candidates.
func (h Hex) FieldOfViewFunc(candidates []Hex, blocked func(Hex) bool) []Hex {
	if len(candidates) == 0 {
		return nil
	}

	return h.AppendFieldOfViewFunc(make([]Hex, 0, len(candidates)), candidates, blocked)
}

// AppendFieldOfViewFunc appends the hexes FieldOfViewFunc returns to dst and returns the
// extended slice, so a caller reusing dst allocates nothing once it has room.
func (h Hex) AppendFieldOfViewFunc(dst []Hex, candidates []Hex, blocked func(Hex) bool) []Hex {
	for hex := range h.FieldOfViewFuncSeq(candidates, blocked) {
		dst = append(dst, hex)
	}

	return dst
}

// FieldOfViewFuncSeq returns an iterator over the hexes FieldOfViewFunc returns, in the same
// order, so a caller walks the visible candidates without a buffer and tests no candidate past
// the one it stops at.
func (h Hex) FieldOfViewFuncSeq(candidates []Hex, blocked func(Hex) bool) iter.Seq[Hex] {
	return func(yield func(Hex) bool) {
		for _, candidate := range candidates {
			if h.HasLineOfSightFunc(candidate, blocked) && !yield(candidate) {
				return
			}
		}
	}
}

// sightLine returns the segment from the center of h to the center of target.
func (h Hex) sightLine(target Hex) sightLine {
	sight := target.Subtract(h)

	return sightLine{sight, sight.dot(sight), sight.width(), sight.edgeStep()}
}

// lineAt returns the hex at step i of the n steps of the line from h to target, both nudged
// by the same amount so the line never runs along a hex boundary.
func (h Hex) lineAt(target Hex, i, n int) Hex {
	const nudge = 1e-6

	start := FractionalHex{float64(h.Q) + nudge, float64(h.R) + 2*nudge}
	end := FractionalHex{float64(target.Q) + nudge, float64(target.R) + 2*nudge}

	return start.Lerp(end, float64(i)/float64(max(n, 1))).Round()
}

// width returns how far a hex reaches to either side of the line from the origin through h, in
// the unit of three times cross: a hex whose center lies nearer the line than that is run
// through, and one exactly that far is touched, at a corner or along an edge.
func (h Hex) width() int {
	q, r, s := h.QRS()

	return max(geom.Abs(q-r), geom.Abs(r-s), geom.Abs(s-q))
}

// edgeStep returns the step from a hex counterclockwise of the line from the origin through h
// to the hex facing it across the line, where the line runs along the edge the two share, as it
// does towards a diagonal neighbor. It is the zero hex for a line that runs along no edge.
func (h Hex) edgeStep() Hex {
	q, r, s := h.QRS()
	if h.IsZero() || q != r && r != s && s != q {
		return Hex{}
	}

	return Hex{(r - s) / h.width(), (s - q) / h.width()}
}

// cross returns the cross product of h and hex as vectors, up to a positive factor: positive
// where hex lies counterclockwise of h, zero where the two are parallel.
func (h Hex) cross(hex Hex) int {
	return h.Q*hex.R - h.R*hex.Q
}

// dot returns the dot product of h and hex as vectors, up to a positive factor: the sum of the
// products of their cube coordinates.
func (h Hex) dot(hex Hex) int {
	return h.Q*hex.Q + h.R*hex.R + (h.Q+h.R)*(hex.Q+hex.R)
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

// To converts the hex into the specified coordinate system, returning a geom.Point[int]. It
// panics for a system outside the seven, as [CoordinateSystem.To] does.
func (h Hex) To(system CoordinateSystem) geom.Point[int] {
	return system.To(h)
}

// Point returns (q,r) as a geom.Point[int].
func (h Hex) Point() geom.Point[int] {
	return geom.Pt(h.Q, h.R)
}

// Float converts the hex to a [FractionalHex], the counterpart of [FractionalHex.Round].
func (h Hex) Float() FractionalHex {
	return FractionalHex{float64(h.Q), float64(h.R)}
}

// String returns a compact representation of the hex as (q,r).
func (h Hex) String() string {
	return fmt.Sprintf("(%s,%s)", geom.String(h.Q), geom.String(h.R))
}
