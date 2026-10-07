package hex

import (
	"iter"

	geom "github.com/gravitton/geometry"
)

// sightLine is the segment between the centers of two hexes as the sight test reads it, every
// hex given as its offset from the one the segment starts at: the step to the target, the
// length of the segment and the reach of a hex to either side of it in the units dot and cross
// return, and the step across the edge the segment runs along, zero where it runs along none.
type sightLine struct {
	sight  Hex
	reach  int
	width  int
	across Hex
}

// clearance returns how far the hex at offset stands clear of the segment: negative where the
// segment runs through it, zero where it only touches it, at a corner or along an edge, and
// positive where it misses it or does not reach it strictly between its ends.
func (l sightLine) clearance(offset Hex) int {
	if along := l.sight.dot(offset); along <= 0 || along >= l.reach {
		return 1
	}

	aside := 3 * l.sight.cross(offset)

	return max(aside, -aside) - l.width
}

// facing returns the hex across the edge the segment runs along from the hex at offset, one the
// segment touches, and whether there is such an edge: a hex touched at a corner faces none. The
// two hexes block the segment only together.
func (l sightLine) facing(offset Hex) (Hex, bool) {
	return offset.Add(l.across.Multiply(geom.Sign(l.sight.cross(offset)))), !l.across.IsZero()
}

// crossed returns an iterator over the hexes the line can meet strictly between its ends, from
// its start towards the target. It takes a step at a time along the cube coordinate the segment
// covers the most of and yields the hexes of that step lying no further from the segment than a
// hex reaches, one or two of them, so every hex that does not stand clear of it is among them.
// The two ends of a step are quotients of a numerator that moves by the same amount at each
// one, so they are carried from step to step and nothing is divided.
func (l sightLine) crossed() iter.Seq[Hex] {
	return func(yield func(Hex) bool) {
		steps := l.sight.Length()
		q, r, s := l.sight.QRS()

		var forward, lateral Hex
		var side int
		switch steps {
		case geom.Abs(s):
			forward, lateral, side = Hex{0, -geom.Sign(s)}, Hex{1, -1}, q
		case geom.Abs(r):
			forward, lateral, side = Hex{-geom.Sign(r), geom.Sign(r)}, Hex{-1, 0}, s
		default:
			forward, lateral, side = Hex{geom.Sign(q), 0}, Hex{0, 1}, r
		}

		drift, divisor := 3*side, 3*steps
		before, after := quotient{0, l.width}, quotient{0, l.width}
		row := Hex{}
		for i := 1; i < steps; i++ {
			row = row.Add(forward)
			before, after = before.carry(-drift, divisor), after.carry(drift, divisor)

			offset := row.Subtract(lateral.Multiply(before.whole))
			for range before.whole + after.whole + 1 {
				if !yield(offset) {
					return
				}

				offset = offset.Add(lateral)
			}
		}
	}
}

// quotient is a numerator divided by a positive divisor and rounded down, kept with the rest of
// the division so the numerator can move without dividing again.
type quotient struct {
	whole int
	rest  int
}

// carry returns the quotient after its numerator moved by delta, which is no larger than the
// divisor either way.
func (q quotient) carry(delta, divisor int) quotient {
	switch rest := q.rest + delta; {
	case rest >= divisor:
		return quotient{q.whole + 1, rest - divisor}
	case rest < 0:
		return quotient{q.whole - 1, rest + divisor}
	default:
		return quotient{q.whole, rest}
	}
}
