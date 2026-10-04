package hex

// cubeBounds is the span of each cube coordinate over a set of hexes: a range, or the hexes two
// ranges share. A span whose minimum passes its maximum holds nothing.
type cubeBounds struct {
	minQ, maxQ int
	minR, maxR int
	minS, maxS int
}

// meet returns the bounds of the hexes within both b and other.
func (b cubeBounds) meet(other cubeBounds) cubeBounds {
	return cubeBounds{
		max(b.minQ, other.minQ), min(b.maxQ, other.maxQ),
		max(b.minR, other.minR), min(b.maxR, other.maxR),
		max(b.minS, other.minS), min(b.maxS, other.maxS),
	}
}

// row returns the first and the last r of the hexes within the bounds at the given q. The
// first passes the last where the row is empty.
func (b cubeBounds) row(q int) (int, int) {
	return max(b.minR, -q-b.maxS), min(b.maxR, -q-b.minS)
}

// size returns the number of hexes within the bounds, summed row by row without walking them.
func (b cubeBounds) size() int {
	size := 0
	for q := b.minQ; q <= b.maxQ; q++ {
		first, last := b.row(q)
		size += max(last-first+1, 0)
	}

	return size
}
