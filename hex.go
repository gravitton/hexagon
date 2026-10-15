package hex

import (
	"cmp"
	"fmt"
	"iter"
	"slices"
	"strings"

	geom "github.com/gravitton/geometry"
)

// Hex is a cell of a hexagonal grid in axial coordinates. The third cube coordinate, s, is -q-r,
// derived and never stored.
type Hex struct {
	Q int `json:"q"`
	R int `json:"r"`
}

// Pt is shorthand for Hex{q, r}.
func Pt(q, r int) Hex {
	return Hex{q, r}
}

// ParseHex returns the hex in the form String prints, and an error for any other string. Each
// coordinate is an integer geom.Parse accepts; a fractional one is an error, not rounded.
func ParseHex(s string) (Hex, error) {
	q, r, err := parseCoordinates[int](s, "hex")
	if err != nil {
		return Hex{}, err
	}

	return Hex{q, r}, nil
}

// RangeLen returns the length of Range and Spiral for the radius (1+3r(r+1)), and zero for a
// negative one.
func RangeLen(radius int) int {
	if radius < 0 {
		return 0
	}

	return 1 + 3*radius*(radius+1)
}

// RingLen returns the length of Ring for the radius (6r), one for a zero radius and zero for a
// negative one.
func RingLen(radius int) int {
	if radius < 0 {
		return 0
	}

	return max(6*radius, 1)
}

// parseCoordinates returns the two coordinates of the form String prints, naming the type of the
// value in its errors.
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

// S returns the third cube coordinate (-q-r).
func (h Hex) S() int {
	return -h.Q - h.R
}

// QR returns Q and R.
func (h Hex) QR() (int, int) {
	return h.Q, h.R
}

// QRS returns Q, R and S.
func (h Hex) QRS() (int, int, int) {
	return h.Q, h.R, h.S()
}

// Length returns the distance from the origin in steps.
func (h Hex) Length() int {
	return (geom.Abs(h.Q) + geom.Abs(h.R) + geom.Abs(h.S())) / 2
}

// Add returns the sum of the hexes.
func (h Hex) Add(hex Hex) Hex {
	return Hex{h.Q + hex.Q, h.R + hex.R}
}

// Subtract returns the hex less the given one.
func (h Hex) Subtract(hex Hex) Hex {
	return Hex{h.Q - hex.Q, h.R - hex.R}
}

// Multiply returns the hex with both coordinates multiplied by the factor.
func (h Hex) Multiply(factor int) Hex {
	return Hex{h.Q * factor, h.R * factor}
}

// Lerp returns the hex the position interpolated toward the given one lies in, extrapolated
// outside [0, 1]. Line is nudged off the boundaries between hexes and Lerp is not, so on a tie the
// two can pick different hexes. It panics for a NaN or infinite t.
func (h Hex) Lerp(hex Hex, t float64) Hex {
	return h.Float().Lerp(hex.Float(), t).Round()
}

// Turn returns the hex turned about the origin by steps sixth-turns of increasing angle; negative
// steps go the other way.
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

// TurnAround returns the hex turned about the pivot by steps sixth-turns of increasing angle;
// negative steps go the other way.
func (h Hex) TurnAround(pivot Hex, steps int) Hex {
	return pivot.Add(h.Subtract(pivot).Turn(steps))
}

// ReflectQ returns the hex reflected across the q-axis, swapping r and s.
func (h Hex) ReflectQ() Hex {
	return Hex{h.Q, h.S()}
}

// ReflectR returns the hex reflected across the r-axis, swapping q and s.
func (h Hex) ReflectR() Hex {
	return Hex{h.S(), h.R}
}

// ReflectS returns the hex reflected across the s-axis, swapping q and r.
func (h Hex) ReflectS() Hex {
	return Hex{h.R, h.Q}
}

// Neighbor returns the neighbor in the direction, and the hex itself for DirectionNone.
func (h Hex) Neighbor(direction Direction) Hex {
	return h.Add(direction.Hex())
}

// Neighbors returns the six neighbors in the order of Directions.
func (h Hex) Neighbors() [6]Hex {
	var neighbors [6]Hex
	for i, direction := range Directions() {
		neighbors[i] = h.Neighbor(direction)
	}

	return neighbors
}

// DiagonalNeighbor returns the hex two steps away between the neighbor in the direction and the
// one after it, a twelfth of a turn past the direction, and the hex itself for DirectionNone.
func (h Hex) DiagonalNeighbor(direction Direction) Hex {
	return h.Neighbor(direction).Neighbor(direction.Turn(1))
}

// DiagonalNeighbors returns the six diagonal neighbors in the order of Directions, each a twelfth
// of a turn past the neighbor of the same index.
func (h Hex) DiagonalNeighbors() [6]Hex {
	var neighbors [6]Hex
	for i, direction := range Directions() {
		neighbors[i] = h.DiagonalNeighbor(direction)
	}

	return neighbors
}

// Range returns the hexes within radius of the hex, itself included, in the order of Compare. It
// returns nil for a negative radius.
func (h Hex) Range(radius int) []Hex {
	if radius < 0 {
		return nil
	}

	return slices.AppendSeq(make([]Hex, 0, RangeLen(radius)), h.RangeSeq(radius))
}

// RangeSeq returns an iterator over the hexes of Range, in the same order.
func (h Hex) RangeSeq(radius int) iter.Seq[Hex] {
	return func(yield func(Hex) bool) {
		for q := -radius; q <= radius; q++ {
			for r := max(-radius, -q-radius); r <= min(radius, radius-q); r++ {
				if !yield(Hex{h.Q + q, h.R + r}) {
					return
				}
			}
		}
	}
}

// Ring returns the hexes at exactly distance radius from the hex, by increasing angle from the
// SMinus corner, so the ring of radius one is Neighbors. It returns nil for a negative radius and
// the hex alone for zero.
func (h Hex) Ring(radius int) []Hex {
	if radius < 0 {
		return nil
	}

	return slices.AppendSeq(make([]Hex, 0, RingLen(radius)), h.RingSeq(radius))
}

// RingSeq returns an iterator over the hexes of Ring, in the same order.
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

// Spiral returns the hexes of Range nearest first: the hex, then ring by ring outward in the order
// of Ring. It returns nil for a negative radius.
func (h Hex) Spiral(radius int) []Hex {
	if radius < 0 {
		return nil
	}

	return slices.AppendSeq(make([]Hex, 0, RangeLen(radius)), h.SpiralSeq(radius))
}

// SpiralSeq returns an iterator over the hexes of Spiral, in the same order.
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

// DistanceTo returns the distance to the given hex in steps.
func (h Hex) DistanceTo(hex Hex) int {
	return h.Subtract(hex).Length()
}

// DirectionTo returns the direction nearest to the angle toward the given hex, and DirectionNone
// for the hex itself. A hex exactly between two directions takes the one the other is a sixth-turn
// past, the direction DiagonalNeighbor reaches it by.
func (h Hex) DirectionTo(hex Hex) Direction {
	delta := hex.Subtract(h)
	for _, direction := range Directions() {
		q, r, s := delta.Turn(-int(direction)).QRS()
		if q >= r && r > s {
			return direction
		}
	}

	return DirectionNone
}

// Line returns the hexes of the straight line from the hex to target, one a step, both ends
// included. The line is nudged off the boundaries between hexes, so it never lands on a tie, except
// far enough from the origin that float64 cannot resolve the nudge.
func (h Hex) Line(target Hex) []Hex {
	return slices.AppendSeq(make([]Hex, 0, h.DistanceTo(target)+1), h.LineSeq(target))
}

// LineSeq returns an iterator over the hexes of Line, in the same order.
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

// HasLineOfSight reports whether the segment between the centers of the hex and target is clear of
// the blocking hexes. Neither end counts as a blocker. It is stricter than Line: a hex off the line
// can block the segment, and a hex of the line the segment only touches does not.
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

// HasLineOfSightFunc reports whether the segment between the centers of the hex and target is
// clear of the hexes blocked reports, exactly when HasLineOfSight is for the same blockers. It asks
// blocked only about the hexes that can block the segment, each at most once, from the hex toward
// target, never about either end. A nil blocked blocks nothing.
func (h Hex) HasLineOfSightFunc(target Hex, blocked func(Hex) bool) bool {
	if blocked == nil {
		return true
	}

	line := h.sightLine(target)

	var previous Hex
	for offset := range line.crossed() {
		switch clearance := line.clearance(offset); {
		case clearance < 0:
			if blocked(h.Add(offset)) {
				return false
			}
		case clearance == 0:
			if facing, edge := line.facing(offset); edge && facing != previous && blocked(h.Add(offset)) && blocked(h.Add(facing)) {
				return false
			}
		}

		previous = offset
	}

	return true
}

// FieldOfView returns the candidates HasLineOfSight sees from the hex past the blocking hexes, in
// the order given, and nil for no candidates.
func (h Hex) FieldOfView(candidates []Hex, blocking []Hex) []Hex {
	if len(candidates) == 0 {
		return nil
	}

	return slices.AppendSeq(make([]Hex, 0, len(candidates)), h.FieldOfViewSeq(candidates, blocking))
}

// FieldOfViewSeq returns an iterator over the hexes of FieldOfView, in the same order.
func (h Hex) FieldOfViewSeq(candidates []Hex, blocking []Hex) iter.Seq[Hex] {
	return func(yield func(Hex) bool) {
		for _, candidate := range candidates {
			if h.HasLineOfSight(candidate, blocking) && !yield(candidate) {
				return
			}
		}
	}
}

// FieldOfViewFunc returns the candidates HasLineOfSightFunc sees from the hex past the hexes
// blocked reports, in the order given, and nil for no candidates.
func (h Hex) FieldOfViewFunc(candidates []Hex, blocked func(Hex) bool) []Hex {
	if len(candidates) == 0 {
		return nil
	}

	return slices.AppendSeq(make([]Hex, 0, len(candidates)), h.FieldOfViewFuncSeq(candidates, blocked))
}

// FieldOfViewFuncSeq returns an iterator over the hexes of FieldOfViewFunc, in the same order.
func (h Hex) FieldOfViewFuncSeq(candidates []Hex, blocked func(Hex) bool) iter.Seq[Hex] {
	return func(yield func(Hex) bool) {
		for _, candidate := range candidates {
			if h.HasLineOfSightFunc(candidate, blocked) && !yield(candidate) {
				return
			}
		}
	}
}

// FieldOfViewRange returns the hexes of Spiral(radius) that FieldOfViewFunc sees past the hexes
// blocked reports, in that order. It asks blocked about each hex at most once, and returns nil for
// a negative radius.
func (h Hex) FieldOfViewRange(radius int, blocked func(Hex) bool) []Hex {
	if radius < 0 {
		return nil
	}

	return slices.AppendSeq(make([]Hex, 0, RangeLen(radius)), h.FieldOfViewRangeSeq(radius, blocked))
}

// FieldOfViewRangeSeq returns an iterator over the hexes of FieldOfViewRange, in the same order.
func (h Hex) FieldOfViewRangeSeq(radius int, blocked func(Hex) bool) iter.Seq[Hex] {
	return func(yield func(Hex) bool) {
		if radius < 0 || !yield(h) {
			return
		}

		var sweep shadowCast
		if blocked != nil && radius > 1 {
			sweep = newShadowCast(radius)
		}

		for distance := 1; distance < radius; distance++ {
			if !sweep.ring(h, distance, blocked, yield) {
				return
			}
		}

		if radius > 0 {
			sweep.ring(h, radius, nil, yield)
		}
	}
}

// sightLine returns the segment from the center of the hex to the center of target.
func (h Hex) sightLine(target Hex) sightLine {
	sight := target.Subtract(h)

	return sightLine{sight, sight.dot(sight), sight.width(), sight.edgeStep()}
}

// lineAt returns the hex at step i of the n steps of the line from the hex to target, both nudged
// by the same amount so the line never runs along a hex boundary.
func (h Hex) lineAt(target Hex, i, n int) Hex {
	const nudge = 1e-6

	start := FractionalHex{float64(h.Q) + nudge, float64(h.R) + 2*nudge}
	end := FractionalHex{float64(target.Q) + nudge, float64(target.R) + 2*nudge}

	return start.Lerp(end, float64(i)/float64(max(n, 1))).Round()
}

// width returns how far a hex reaches to either side of the line from the origin through this
// one, in the unit of three times cross: a hex whose center lies nearer the line than that is run
// through, and one exactly that far is touched, at a corner or along an edge.
func (h Hex) width() int {
	q, r, s := h.QRS()

	return max(geom.Abs(q-r), geom.Abs(r-s), geom.Abs(s-q))
}

// edgeStep returns the step from a hex counterclockwise of the line from the origin through this
// one to the hex facing it across the line, where the line runs along the edge the two share, as
// it does toward a diagonal neighbor. It is the zero hex for a line that runs along no edge.
func (h Hex) edgeStep() Hex {
	q, r, s := h.QRS()
	if h.IsZero() || q != r && r != s && s != q {
		return Hex{}
	}

	width := h.width()

	return Hex{(r - s) / width, (s - q) / width}
}

// cross returns the cross product of the hexes as vectors, up to a positive factor: positive
// where the given hex lies counterclockwise of this one, zero where the two are parallel.
func (h Hex) cross(hex Hex) int {
	return h.Q*hex.R - h.R*hex.Q
}

// dot returns the dot product of the hexes as vectors, up to a positive factor: the sum of the
// products of their cube coordinates.
func (h Hex) dot(hex Hex) int {
	return h.Q*hex.Q + h.R*hex.R + (h.Q+h.R)*(hex.Q+hex.R)
}

// Equal reports whether the hexes have equal coordinates.
func (h Hex) Equal(hex Hex) bool {
	return h == hex
}

// Compare returns -1, 0 or +1 as the hex sorts before, with or after the given one, by Q and then
// by R, as cmp.Compare does.
func (h Hex) Compare(hex Hex) int {
	if c := cmp.Compare(h.Q, hex.Q); c != 0 {
		return c
	}

	return cmp.Compare(h.R, hex.R)
}

// IsZero reports whether the hex is the origin.
func (h Hex) IsZero() bool {
	return h == Hex{}
}

// To returns the coordinates of the hex in the system. It panics for a system outside the seven.
func (h Hex) To(system CoordinateSystem) geom.Point[int] {
	return system.To(h)
}

// Point returns Q and R as a point.
func (h Hex) Point() geom.Point[int] {
	return geom.Pt(h.Q, h.R)
}

// Vector returns Q and R as a vector.
func (h Hex) Vector() geom.Vector[int] {
	return geom.Vec(h.Q, h.R)
}

// Float converts the hex to a FractionalHex.
func (h Hex) Float() FractionalHex {
	return FractionalHex{float64(h.Q), float64(h.R)}
}

// String returns the hex as its coordinates: (q,r).
func (h Hex) String() string {
	return h.Point().String()
}
