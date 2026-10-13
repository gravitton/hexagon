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

// arc is the open range of directions between two vectors from the source of a sweep,
// counterclockwise from the first to the second, neither of them inside it: the directions a
// blocking hex hides, between two of its corners.
type arc struct {
	from Hex
	to   Hex
}

// contains reports whether the direction lies strictly inside the arc.
func (a arc) contains(direction Hex) bool {
	return a.from.cross(direction) > 0 && direction.cross(a.to) > 0
}

// covers reports whether the other arc lies strictly inside the arc, both its ends included.
func (a arc) covers(other arc) bool {
	return a.from.cross(other.from) > 0 && other.to.cross(a.to) > 0
}

// wedge is a sixth of the plane around the source of a sweep: the directions from one axis to
// the next, each given as the step to the neighbor on it. A hex off the axes lies inside one
// wedge with all its corners; a hex on an axis is halved by it between two.
type wedge struct {
	first  Hex
	second Hex
}

// arcOf returns the arc the hex at offset spans as seen from the source, between the two of its
// corners furthest apart, each as three times its offset so it stays in integers. The offset
// lies inside the wedge or on one of its axes.
func (w wedge) arcOf(offset Hex) arc {
	center := offset.Multiply(3)

	from, other := center.Add(w.first.Multiply(2)).Subtract(w.second), center.Add(w.first).Subtract(w.second.Multiply(2))
	if from.cross(other) < 0 {
		from = other
	}

	to, other := center.Add(w.second.Multiply(2)).Subtract(w.first), center.Add(w.second).Subtract(w.first.Multiply(2))
	if to.cross(other) > 0 {
		to = other
	}

	return arc{from, to}
}

// shadowCast is a sweep of the hexes around a source, ring by ring outwards, keeping the arcs
// the blocking hexes met so far hide. The shadows of the wedges follow one another, those of a
// wedge in order and apart, and the hexes on an axis are kept as a flag: only a hex on the same
// axis hides them.
type shadowCast struct {
	shadows []arc
	ends    [6]int
	cast    []arc
	axes    [6]bool
}

// newShadowCast returns a sweep with room for the shadows of every ring within radius: three to
// a step of it at most, where every other hex of the last ring asked about blocks. The shadows
// swept against and those being cast share one buffer, a half each.
func newShadowCast(radius int) shadowCast {
	room := 3 * radius
	buffer := make([]arc, 2*room)

	return shadowCast{shadows: buffer[:0:room], cast: buffer[room : room : 2*room]}
}

// ring sweeps the ring at the given distance of source: it yields the hexes of it no shadow
// hides, in the order of Ring, and adds the shadows of those blocked reports. A hex hidden with
// both ends of its arc casts no shadow of its own and is not asked about; a nil blocked asks
// about none, for the last ring. It reports whether every hex was yielded.
func (c *shadowCast) ring(source Hex, distance int, blocked func(Hex) bool, yield func(Hex) bool) bool {
	directions := Directions()

	var seen, blocks [6]bool
	for i, direction := range directions {
		seen[i] = !c.axes[i]
		blocks[i] = seen[i] && blocked != nil && blocked(source.Add(direction.Hex().Multiply(distance)))
		c.axes[i] = c.axes[i] || blocks[i]
	}

	start := 0
	for i, direction := range directions {
		sector := wedge{direction.Hex(), direction.Turn(1).Hex()}
		axis, side := sector.first.Multiply(distance), sector.second.Subtract(sector.first)
		if seen[i] && !yield(source.Add(axis)) {
			return false
		}

		shadows := c.shadows[start:c.ends[i]]
		start = c.ends[i]

		merge := shading{shadows: shadows, cast: c.cast}
		edge, joins := Hex{}, false
		if blocks[i] {
			edge, joins = sector.arcOf(axis).to, true
			merge.add(arc{sector.first, edge}, false)
		}

		for j := 1; j < distance; j++ {
			offset := axis.Add(side.Multiply(j))
			span := sector.arcOf(offset)

			for len(shadows) > 0 && offset.cross(shadows[0].to) <= 0 {
				shadows = shadows[1:]
			}

			hidden := len(shadows) > 0 && shadows[0].contains(offset)
			if !hidden && !yield(source.Add(offset)) {
				return false
			}
			if hidden && shadows[0].covers(span) || blocked == nil || !blocked(source.Add(offset)) {
				joins = false
				continue
			}

			merge.add(span, joins && edge.cross(span.from) == 0)
			edge, joins = span.to, true
		}

		if next := (i + 1) % len(directions); blocks[next] {
			span := sector.arcOf(sector.second.Multiply(distance))
			merge.add(arc{span.from, sector.second}, joins && edge.cross(span.from) == 0)
		}

		c.cast = merge.close()
		c.ends[i] = len(c.cast)
	}

	c.shadows, c.cast = c.cast, c.shadows[:0]

	return true
}

// shading merges the shadows of a wedge with the arcs a ring adds to them, both in order, into
// the shadows the next ring is swept against.
type shading struct {
	shadows []arc
	cast    []arc
	last    arc
	open    bool
}

// add merges the shadows that start before the arc, then the arc itself. Two arcs that only meet
// stay apart, since the direction they meet in runs past a corner of each and is hidden by
// neither, unless joined says they are the two hexes of one ring beside the edge it runs along.
func (s *shading) add(shadow arc, joined bool) {
	for len(s.shadows) > 0 && s.shadows[0].from.cross(shadow.from) > 0 {
		s.merge(s.shadows[0], false)
		s.shadows = s.shadows[1:]
	}

	s.merge(shadow, joined)
}

// merge extends the last shadow by one that overlaps it and starts another otherwise.
func (s *shading) merge(shadow arc, joined bool) {
	gap := s.last.to.cross(shadow.from)
	if !s.open || gap > 0 || gap == 0 && !joined {
		s.flush()
		s.last, s.open = shadow, true

		return
	}

	if s.last.to.cross(shadow.to) > 0 {
		s.last.to = shadow.to
	}
}

// flush appends the last shadow, if there is one.
func (s *shading) flush() {
	if s.open {
		s.cast = append(s.cast, s.last)
	}
}

// close merges the shadows left and returns them all.
func (s *shading) close() []arc {
	for _, shadow := range s.shadows {
		s.merge(shadow, false)
	}

	s.flush()

	return s.cast
}
