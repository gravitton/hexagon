package hex

import (
	"encoding/json"
	"fmt"
	"iter"
	"slices"
	"strings"
)

// Region is a set of hexes, for area queries: the hexes a unit can reach, a room, the overlap of
// two ranges. It holds each hex once, ordered by q and then by r like [Hex.Compare], and every
// method returns a new region, so a Region is safe to copy and share. The zero value is the
// empty region.
type Region struct {
	hexes []Hex
}

// Reg returns the region of the given hexes, in any order and with any repeats. It copies the
// slice, so the region does not change with it.
func Reg(hexes []Hex) Region {
	if len(hexes) == 0 {
		return Region{}
	}

	sorted := slices.Clone(hexes)
	slices.SortFunc(sorted, Hex.Compare)

	return Region{slices.Compact(sorted)}
}

// Hexes returns the hexes of the region as a fresh slice, ordered by q and then by r like
// [Hex.Compare]. It returns nil for the empty region.
func (r Region) Hexes() []Hex {
	return slices.Clone(r.hexes)
}

// AppendHexes appends the hexes Hexes returns to dst and returns the extended slice, so a
// caller reusing dst allocates nothing once it has room.
func (r Region) AppendHexes(dst []Hex) []Hex {
	return append(dst, r.hexes...)
}

// All returns an iterator over the hexes of the region, in the order Hexes returns them, so a
// caller walks them without a buffer.
func (r Region) All() iter.Seq[Hex] {
	return slices.Values(r.hexes)
}

// Len returns the number of hexes in the region.
func (r Region) Len() int {
	return len(r.hexes)
}

// Union returns the region of the hexes in r, in other or in both.
func (r Region) Union(other Region) Region {
	return r.merge(other, func(int) bool {
		return true
	})
}

// Intersection returns the region of the hexes in both r and other.
func (r Region) Intersection(other Region) Region {
	return r.merge(other, func(side int) bool {
		return side == 0
	})
}

// Difference returns the region of the hexes in r that are not in other.
func (r Region) Difference(other Region) Region {
	return r.merge(other, func(side int) bool {
		return side < 0
	})
}

// merge returns the region of the hexes of r and other on a side to keep, walked once to count
// them and once to place them, so the result is allocated once at its exact size and needs no
// sort.
func (r Region) merge(other Region, keep func(side int) bool) Region {
	size := 0
	r.walk(other, func(_ Hex, side int) {
		if keep(side) {
			size++
		}
	})
	if size == 0 {
		return Region{}
	}

	hexes := make([]Hex, 0, size)
	r.walk(other, func(hex Hex, side int) {
		if keep(side) {
			hexes = append(hexes, hex)
		}
	})

	return Region{hexes}
}

// walk visits every hex of r and other once, in order, with the side it lies on: negative in r
// alone, zero in both, positive in other alone. It steps through the two regions together, so
// no hex is searched for.
func (r Region) walk(other Region, visit func(hex Hex, side int)) {
	i, j := 0, 0
	for i < len(r.hexes) && j < len(other.hexes) {
		switch side := r.hexes[i].Compare(other.hexes[j]); {
		case side < 0:
			visit(r.hexes[i], side)
			i++
		case side > 0:
			visit(other.hexes[j], side)
			j++
		default:
			visit(r.hexes[i], side)
			i++
			j++
		}
	}

	for _, hex := range r.hexes[i:] {
		visit(hex, -1)
	}

	for _, hex := range other.hexes[j:] {
		visit(hex, 1)
	}
}

// mapHexes applies fn to every hex in order and returns the results in a new slice.
func (r Region) mapHexes[R any](fn func(Hex) R) []R {
	mapped := make([]R, len(r.hexes))
	for i, hex := range r.hexes {
		mapped[i] = fn(hex)
	}

	return mapped
}

// Border returns the region of the hexes of r with at least one neighbor outside it: its
// outline, the inner side of every edge the region has, around a hole as well.
func (r Region) Border() Region {
	size := 0
	r.sweep(func(Hex) {
		size++
	})
	if size == 0 {
		return Region{}
	}

	hexes := make([]Hex, 0, size)
	r.sweep(func(hex Hex) {
		hexes = append(hexes, hex)
	})

	return Region{hexes}
}

// sweep visits the hexes of r with a neighbor outside it, in order. The neighbors along q are
// the hexes stored beside the hex, and the two in each column beside it are stored together, at
// a position that only moves forward as the sweep does, so no neighbor is searched for.
func (r Region) sweep(visit func(Hex)) {
	before, after := 0, 0
	for i, hex := range r.hexes {
		before = r.seek(before, Hex{hex.Q - 1, hex.R})
		after = r.seek(after, Hex{hex.Q + 1, hex.R - 1})

		enclosed := r.holds(i-1, Hex{hex.Q, hex.R - 1}) &&
			r.holds(i+1, Hex{hex.Q, hex.R + 1}) &&
			r.holds(before, Hex{hex.Q - 1, hex.R}) &&
			r.holds(before+1, Hex{hex.Q - 1, hex.R + 1}) &&
			r.holds(after, Hex{hex.Q + 1, hex.R - 1}) &&
			r.holds(after+1, Hex{hex.Q + 1, hex.R})
		if !enclosed {
			visit(hex)
		}
	}
}

// seek returns the first index from the given one whose hex does not come before the hex, or
// the length of the region where every hex does.
func (r Region) seek(index int, hex Hex) int {
	for index < len(r.hexes) && r.hexes[index].Compare(hex) < 0 {
		index++
	}

	return index
}

// holds reports whether the region stores the hex at the index.
func (r Region) holds(index int, hex Hex) bool {
	return index >= 0 && index < len(r.hexes) && r.hexes[index] == hex
}

// Contains reports whether the hex is in the region.
func (r Region) Contains(hex Hex) bool {
	_, found := slices.BinarySearchFunc(r.hexes, hex, Hex.Compare)

	return found
}

// Equal reports whether r and other hold the same hexes.
func (r Region) Equal(other Region) bool {
	return slices.Equal(r.hexes, other.hexes)
}

// IsEmpty reports whether the region holds no hex.
func (r Region) IsEmpty() bool {
	return len(r.hexes) == 0
}

// String returns the hexes of the region in order, as Reg((q,r);(q,r)).
func (r Region) String() string {
	return fmt.Sprintf("Reg(%s)", strings.Join(r.mapHexes(Hex.String), ";"))
}

// MarshalJSON implements json.Marshaler with the hexes of the region as an array, in order, and
// an empty array for the empty region.
func (r Region) MarshalJSON() ([]byte, error) {
	if r.IsEmpty() {
		return []byte("[]"), nil
	}

	return json.Marshal(r.hexes)
}

// UnmarshalJSON implements json.Unmarshaler, reading an array of hexes in any order and with
// any repeats, as Reg does.
func (r *Region) UnmarshalJSON(data []byte) error {
	var hexes []Hex
	if err := json.Unmarshal(data, &hexes); err != nil {
		return err
	}

	*r = Reg(hexes)

	return nil
}
