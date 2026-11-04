package hex

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

// span returns the arc the hex at offset spans as seen from the source, between the two of its
// corners furthest apart, each as three times its offset so it stays in integers. The offset
// lies inside the wedge or on one of its axes.
func (w wedge) span(offset Hex) arc {
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

// view is the part of every ring a sweep covers: the hexes within turns sixth-turns about a
// direction, by position along the ring from the axis first, the one before the start of the cone.
type view struct {
	first Direction
	turns int
}

// newView returns the view of the hexes within turns sixth-turns about the direction, half of
// them either side, as Sector takes them.
func newView(direction Direction, turns int) view {
	turns = min(turns, 6)

	return view{direction.Turn(-(turns + 1) / 2), turns}
}

// runs returns the positions of the ring at distance the view sees, and those whose blocking can
// hide one of them. A border at an odd number of sixth-turns runs along the edge between two hexes
// of an odd ring, and the hex outside closes it together with the one inside; a border on an axis
// or through a center touches no hex outside.
func (v view) runs(distance int) (run, run) {
	ring := 6 * distance
	reach := v.turns * distance / 2
	from := (v.turns+1)/2*distance - reach
	seen := run{from, min(from+2*reach, ring-1)}

	if v.turns%2 == 1 && distance%2 == 1 {
		return seen, run{seen.from - 1, min(seen.to+1, ring-1)}
	}

	return seen, seen
}

// run is the positions along a ring between two, both included.
type run struct {
	from int
	to   int
}

// contains reports whether the position lies in the run.
func (r run) contains(position int) bool {
	return r.from <= position && position <= r.to
}

// shadowCast is a sweep of the hexes around a source, ring by ring outward, keeping the arcs
// the blocking hexes met so far hide. The shadows of the wedges follow one another, those of a
// wedge in order and apart, and the hexes on an axis are kept as a flag: only a hex on the same
// axis hides them.
type shadowCast struct {
	shadows []arc
	ends    [6]int
	cast    []arc
	axes    [6]bool
}

// newShadowCast returns a sweep with room for the shadows of every ring within radius: one for
// each hex a wedge has on the last ring asked about, the two on its axes included. A wedge holds
// no more after ring d: from one axis step to the next, a shadow covers at least the 3/(3k+1)
// of a hex off the axes on ring k, or the 2/(3k+1) of one on an axis, and two at most start on
// an axis, so d+1 of them fill it. The shadows swept against and
// those being cast share one buffer, a half each.
func newShadowCast(radius int) shadowCast {
	room := 6 * radius
	buffer := make([]arc, 2*room)

	return shadowCast{shadows: buffer[:0:room], cast: buffer[room : room : 2*room]}
}

// ring sweeps the part of the ring at the given distance of source the view covers, the same view
// every ring: it yields the hexes of it no shadow hides, by increasing angle from the axis first,
// and adds the shadows of those blocked reports. It asks only about the hexes that can hide one of
// the view; a hex hidden with both ends of its arc casts no shadow of its own and is not asked
// about either, and a nil blocked asks about none, for the last ring. It reports whether every hex
// was yielded.
func (c *shadowCast) ring(source Hex, distance int, part view, blocked func(Hex) bool, yield func(Hex) bool) bool {
	visible, asked := part.runs(distance)

	var seen, blocks [6]bool
	for k := range seen {
		axis := part.first.Turn(k).Hex().Multiply(distance)
		seen[k] = !c.axes[k] && visible.contains(k*distance)
		blocks[k] = !c.axes[k] && asked.contains(k*distance) && blocked != nil && blocked(source.Add(axis))
		c.axes[k] = c.axes[k] || blocks[k]
	}

	start := 0
	for k := range seen {
		direction := part.first.Turn(k)
		sector := wedge{direction.Hex(), direction.Turn(1).Hex()}
		axis, side := sector.first.Multiply(distance), sector.second.Subtract(sector.first)
		if seen[k] && !yield(source.Add(axis)) {
			return false
		}

		shadows := c.shadows[start:c.ends[k]]
		start = c.ends[k]

		merge := shading{shadows: shadows, cast: c.cast}
		edge, joins := Hex{}, false
		if blocks[k] {
			edge, joins = sector.span(axis).to, true
			merge.add(arc{sector.first, edge}, false)
		}

		position := k * distance
		for j := max(1, asked.from-position); j < min(distance, asked.to-position+1); j++ {
			offset := axis.Add(side.Multiply(j))
			span := sector.span(offset)

			for len(shadows) > 0 && offset.cross(shadows[0].to) <= 0 {
				shadows = shadows[1:]
			}

			hidden := len(shadows) > 0 && shadows[0].contains(offset)
			if !hidden && visible.contains(position+j) && !yield(source.Add(offset)) {
				return false
			}
			covered := hidden && shadows[0].covers(span)
			if covered || blocked == nil || !blocked(source.Add(offset)) {
				joins = false
				continue
			}

			merge.add(span, joins && edge.cross(span.from) == 0)
			edge, joins = span.to, true
		}

		if next := (k + 1) % len(blocks); blocks[next] {
			span := sector.span(sector.second.Multiply(distance))
			merge.add(arc{span.from, sector.second}, joins && edge.cross(span.from) == 0)
		}

		c.cast = merge.close()
		c.ends[k] = len(c.cast)
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
