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
// hex reaches, one or two of them and found by division, so every hex that does not stand clear
// of it is among them.
func (l sightLine) crossed() iter.Seq[Hex] {
	return func(yield func(Hex) bool) {
		steps := l.sight.Length()

		q, r, s := l.sight.QRS()
		cube := [3]int{q, r, s}

		axis := 0
		for i, coordinate := range cube {
			if geom.Abs(coordinate) == steps {
				axis = i
			}
		}

		second, third := (axis+1)%3, (axis+2)%3

		forward, side, divisor := geom.Sign(cube[axis]), cube[second], 3*steps
		for i := 1; i < steps; i++ {
			var offset [3]int
			offset[axis] = forward * i

			first, last := -floorDiv(l.width-3*side*i, divisor), floorDiv(l.width+3*side*i, divisor)
			for lateral := first; lateral <= last; lateral++ {
				offset[second], offset[third] = lateral, -offset[axis]-lateral

				if !yield(Hex{offset[0], offset[1]}) {
					return
				}
			}
		}
	}
}

// floorDiv returns n divided by a positive m and rounded down, where the / operator rounds
// towards zero.
func floorDiv(n, m int) int {
	return (n - geom.Mod(n, m)) / m
}
