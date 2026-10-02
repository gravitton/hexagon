package hex

import (
	"fmt"

	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/ints"
)

// CoordinateSystem lists supported hexagonal grid coordinate systems.
type CoordinateSystem int

const (
	// Axial stores q and r directly; s is derived as -q-r.
	// This is the native system used throughout the package.
	Axial CoordinateSystem = iota

	// OffsetOddR is a pointy-top system where odd rows (r & 1 == 1) are shifted
	// right by half a column. Useful for rectangular array storage.
	OffsetOddR

	// OffsetEvenR is a pointy-top system where even rows (r & 1 == 0) are shifted
	// right by half a column. Use OffsetOddR or OffsetEvenR depending on which
	// row parity your data places at the left edge.
	OffsetEvenR

	// OffsetOddQ is a flat-top system where odd columns (q & 1 == 1) are shifted
	// down by half a row. Useful for rectangular array storage.
	OffsetOddQ

	// OffsetEvenQ is a flat-top system where even columns (q & 1 == 0) are shifted
	// down by half a row. Use OffsetOddQ or OffsetEvenQ depending on which
	// column parity your data places at the top edge.
	OffsetEvenQ

	// DoubleWidth is a pointy-top system that doubles the column axis: col = 2q+r, row = r.
	// All six neighbors are reachable with fixed offsets — no per-cell parity check needed.
	// It addresses only the cells whose col and row share a parity; From maps any other
	// coordinate to the cell beside it on the side of lower q, the same at every position.
	DoubleWidth

	// DoubleHeight is a flat-top system that doubles the row axis: col = q, row = 2r+q.
	// All six neighbors are reachable with fixed offsets — no per-cell parity check needed.
	// It addresses only the cells whose col and row share a parity; From maps any other
	// coordinate to the cell beside it on the side of lower r, the same at every position.
	DoubleHeight

	// CoordinateSystemNone is the absence of a coordinate system. It names no layout, so
	// Offsets, Offset, To and From have nothing to convert for it and panic.
	CoordinateSystemNone CoordinateSystem = -1
)

// CoordinateSystems lists the seven coordinate systems in order. It returns a fresh array, so
// a caller cannot alter the list.
func CoordinateSystems() [7]CoordinateSystem {
	return [7]CoordinateSystem{Axial, OffsetOddR, OffsetEvenR, OffsetOddQ, OffsetEvenQ, DoubleWidth, DoubleHeight}
}

// ParseCoordinateSystem returns the coordinate system with the given name, as String prints it,
// and an error for any other string. "None" parses to CoordinateSystemNone.
func ParseCoordinateSystem(name string) (CoordinateSystem, error) {
	if name == CoordinateSystemNone.String() {
		return CoordinateSystemNone, nil
	}

	for _, system := range CoordinateSystems() {
		if system.String() == name {
			return system, nil
		}
	}

	return CoordinateSystemNone, fmt.Errorf("hex: unknown coordinate system %q", name)
}

// Offsets returns the 6 neighbor offsets for the given coordinate index in this system,
// indexed by [Direction]. For offset systems this accounts for row/column parity.
func (s CoordinateSystem) Offsets(index ints.Point) [6]ints.Vector {
	switch s {
	case OffsetOddR:
		return offsetsOddR[index.Y&1]
	case OffsetEvenR:
		return offsetsEvenR[index.Y&1]
	case OffsetOddQ:
		return offsetsOddQ[index.X&1]
	case OffsetEvenQ:
		return offsetsEvenQ[index.X&1]
	case DoubleWidth:
		return offsetsDoubleWidth
	case DoubleHeight:
		return offsetsDoubleHeight
	case Axial:
		return directionOffsets
	default:
		panic(fmt.Sprintf("hex: unknown coordinate system %d", s))
	}
}

// Offset returns the neighbor offset vector for the given coordinate index and direction
// in this system. Out-of-range directions wrap into [SMinus, QPlus], and [DirectionNone]
// gives the zero vector. It panics for a system outside the seven, whatever the direction.
func (s CoordinateSystem) Offset(index ints.Point, direction Direction) ints.Vector {
	offsets := s.Offsets(index)
	if direction.IsNone() {
		return ints.Vector{}
	}

	return offsets[direction.normalize()]
}

// To converts an axial hex into this coordinate system as an ints.Point.
func (s CoordinateSystem) To(hex Hex) ints.Point {
	switch s {
	case OffsetOddR:
		return toOffsetOddR(hex)
	case OffsetEvenR:
		return toOffsetEvenR(hex)
	case OffsetOddQ:
		return toOffsetOddQ(hex)
	case OffsetEvenQ:
		return toOffsetEvenQ(hex)
	case DoubleWidth:
		return toDoubleWidth(hex)
	case DoubleHeight:
		return toDoubleHeight(hex)
	case Axial:
		return toAxial(hex)
	default:
		panic(fmt.Sprintf("hex: unknown coordinate system %d", s))
	}
}

// From converts a coordinate in this system into an axial Hex.
func (s CoordinateSystem) From(index ints.Point) Hex {
	switch s {
	case OffsetOddR:
		return fromOffsetOddR(index)
	case OffsetEvenR:
		return fromOffsetEvenR(index)
	case OffsetOddQ:
		return fromOffsetOddQ(index)
	case OffsetEvenQ:
		return fromOffsetEvenQ(index)
	case DoubleWidth:
		return fromDoubleWidth(index)
	case DoubleHeight:
		return fromDoubleHeight(index)
	case Axial:
		return fromAxial(index)
	default:
		panic(fmt.Sprintf("hex: unknown coordinate system %d", s))
	}
}

// IsNone reports whether the coordinate system is none of the seven. Like geom.Orientation, a
// CoordinateSystem outside the constants is not normalized, so every such value counts as
// CoordinateSystemNone.
func (s CoordinateSystem) IsNone() bool {
	switch s {
	case Axial, OffsetOddR, OffsetEvenR, OffsetOddQ, OffsetEvenQ, DoubleWidth, DoubleHeight:
		return false
	default:
		return true
	}
}

// String returns the name of the coordinate system.
func (s CoordinateSystem) String() string {
	switch s {
	case Axial:
		return "Axial"
	case OffsetOddR:
		return "OffsetOddR"
	case OffsetEvenR:
		return "OffsetEvenR"
	case OffsetOddQ:
		return "OffsetOddQ"
	case OffsetEvenQ:
		return "OffsetEvenQ"
	case DoubleWidth:
		return "DoubleWidth"
	case DoubleHeight:
		return "DoubleHeight"
	default:
		return "None"
	}
}

// MarshalText implements encoding.TextMarshaler with the name String prints, so a coordinate
// system is stored as "OffsetOddR" in JSON and as a map key rather than as its number.
func (s CoordinateSystem) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler, the inverse of MarshalText through
// ParseCoordinateSystem.
func (s *CoordinateSystem) UnmarshalText(text []byte) error {
	system, err := ParseCoordinateSystem(string(text))
	if err != nil {
		return err
	}

	*s = system

	return nil
}

// toAxial returns the axial (q,r) as an ints.Point.
func toAxial(hex Hex) ints.Point {
	return geom.Pt(hex.Q, hex.R)
}

// fromAxial converts an ints.Point (q,r) into an axial Hex.
func fromAxial(index ints.Point) Hex {
	return Hex{index.X, index.Y}
}

// toOffsetOddR converts axial to odd-r offset coordinates.
// Odd rows are shifted right by +1/2 column.
func toOffsetOddR(hex Hex) ints.Point {
	parity := hex.R & 1

	col := hex.Q + (hex.R-parity)/2
	row := hex.R

	return geom.Pt(col, row)
}

// fromOffsetOddR converts an odd-r offset coordinate to axial.
func fromOffsetOddR(index ints.Point) Hex {
	parity := index.Y & 1

	q := index.X - (index.Y-parity)/2
	r := index.Y

	return Hex{q, r}
}

// toOffsetEvenR converts axial to even-r offset coordinates.
// Even rows are shifted right by +1/2 column.
func toOffsetEvenR(hex Hex) ints.Point {
	parity := hex.R & 1

	col := hex.Q + (hex.R+parity)/2
	row := hex.R

	return geom.Pt(col, row)
}

// fromOffsetEvenR converts an even-r offset coordinate to axial.
func fromOffsetEvenR(index ints.Point) Hex {
	parity := index.Y & 1

	q := index.X - (index.Y+parity)/2
	r := index.Y

	return Hex{q, r}
}

// toOffsetOddQ converts axial to odd-q offset coordinates.
// Odd columns are shifted down by +1/2 row.
func toOffsetOddQ(hex Hex) ints.Point {
	parity := hex.Q & 1

	col := hex.Q
	row := hex.R + (hex.Q-parity)/2

	return geom.Pt(col, row)
}

// fromOffsetOddQ converts an odd-q offset coordinate to axial.
func fromOffsetOddQ(index ints.Point) Hex {
	parity := index.X & 1

	q := index.X
	r := index.Y - (index.X-parity)/2

	return Hex{q, r}
}

// toOffsetEvenQ converts axial to even-q offset coordinates.
// Even columns are shifted down by +1/2 row.
func toOffsetEvenQ(hex Hex) ints.Point {
	parity := hex.Q & 1

	col := hex.Q
	row := hex.R + (hex.Q+parity)/2

	return geom.Pt(col, row)
}

// fromOffsetEvenQ converts an even-q offset coordinate to axial.
func fromOffsetEvenQ(index ints.Point) Hex {
	parity := index.X & 1

	q := index.X
	r := index.Y - (index.X+parity)/2

	return Hex{q, r}
}

// toDoubleWidth converts axial to double-width coordinates (doubling the q axis).
func toDoubleWidth(hex Hex) ints.Point {
	col := 2*hex.Q + hex.R
	row := hex.R

	return geom.Pt(col, row)
}

// fromDoubleWidth converts a double-width coordinate to axial, flooring a coordinate of mixed
// parity.
func fromDoubleWidth(index ints.Point) Hex {
	q := (index.X - index.Y) >> 1
	r := index.Y

	return Hex{q, r}
}

// toDoubleHeight converts axial to double-height coordinates (doubling the r axis).
func toDoubleHeight(hex Hex) ints.Point {
	col := hex.Q
	row := 2*hex.R + hex.Q

	return geom.Pt(col, row)
}

// fromDoubleHeight converts a double-height coordinate to axial, flooring a coordinate of mixed
// parity.
func fromDoubleHeight(index ints.Point) Hex {
	q := index.X
	r := (index.Y - index.X) >> 1

	return Hex{q, r}
}

// offsetsOddR lists neighbor vectors for odd-r offset coordinates as
// [parityRow][direction], where parityRow=0 for even rows and 1 for odd rows.
var offsetsOddR = [2][6]ints.Vector{
	{
		geom.Vec(1, 0),   // SMinus
		geom.Vec(0, 1),   // RPlus
		geom.Vec(-1, 1),  // QMinus
		geom.Vec(-1, 0),  // SPlus
		geom.Vec(-1, -1), // RMinus
		geom.Vec(0, -1),  // QPlus
	},
	{
		geom.Vec(1, 0),  // SMinus
		geom.Vec(1, 1),  // RPlus
		geom.Vec(0, 1),  // QMinus
		geom.Vec(-1, 0), // SPlus
		geom.Vec(0, -1), // RMinus
		geom.Vec(1, -1), // QPlus
	},
}

// offsetsEvenR lists neighbor vectors for even-r offset coordinates as
// [parityRow][direction], where parityRow=0 for even rows and 1 for odd rows.
var offsetsEvenR = [2][6]ints.Vector{
	{
		geom.Vec(1, 0),  // SMinus
		geom.Vec(1, 1),  // RPlus
		geom.Vec(0, 1),  // QMinus
		geom.Vec(-1, 0), // SPlus
		geom.Vec(0, -1), // RMinus
		geom.Vec(1, -1), // QPlus
	},
	{
		geom.Vec(1, 0),   // SMinus
		geom.Vec(0, 1),   // RPlus
		geom.Vec(-1, 1),  // QMinus
		geom.Vec(-1, 0),  // SPlus
		geom.Vec(-1, -1), // RMinus
		geom.Vec(0, -1),  // QPlus
	},
}

// offsetsOddQ lists neighbor vectors for odd-q offset coordinates as
// [parityCol][direction], where parityCol=0 for even columns and 1 for odd columns.
var offsetsOddQ = [2][6]ints.Vector{
	{
		geom.Vec(1, 0),   // SMinus
		geom.Vec(0, 1),   // RPlus
		geom.Vec(-1, 0),  // QMinus
		geom.Vec(-1, -1), // SPlus
		geom.Vec(0, -1),  // RMinus
		geom.Vec(1, -1),  // QPlus
	},
	{
		geom.Vec(1, 1),  // SMinus
		geom.Vec(0, 1),  // RPlus
		geom.Vec(-1, 1), // QMinus
		geom.Vec(-1, 0), // SPlus
		geom.Vec(0, -1), // RMinus
		geom.Vec(1, 0),  // QPlus
	},
}

// offsetsEvenQ lists neighbor vectors for even-q offset coordinates as
// [parityCol][direction], where parityCol=0 for even columns and 1 for odd columns.
var offsetsEvenQ = [2][6]ints.Vector{
	{
		geom.Vec(1, 1),  // SMinus
		geom.Vec(0, 1),  // RPlus
		geom.Vec(-1, 1), // QMinus
		geom.Vec(-1, 0), // SPlus
		geom.Vec(0, -1), // RMinus
		geom.Vec(1, 0),  // QPlus
	},
	{
		geom.Vec(1, 0),   // SMinus
		geom.Vec(0, 1),   // RPlus
		geom.Vec(-1, 0),  // QMinus
		geom.Vec(-1, -1), // SPlus
		geom.Vec(0, -1),  // RMinus
		geom.Vec(1, -1),  // QPlus
	},
}

// offsetsDoubleWidth lists neighbor vectors for double-width coordinates.
var offsetsDoubleWidth = [6]ints.Vector{
	geom.Vec(2, 0),   // SMinus
	geom.Vec(1, 1),   // RPlus
	geom.Vec(-1, 1),  // QMinus
	geom.Vec(-2, 0),  // SPlus
	geom.Vec(-1, -1), // RMinus
	geom.Vec(1, -1),  // QPlus
}

// offsetsDoubleHeight lists neighbor vectors for double-height coordinates.
var offsetsDoubleHeight = [6]ints.Vector{
	geom.Vec(1, 1),   // SMinus
	geom.Vec(0, 2),   // RPlus
	geom.Vec(-1, 1),  // QMinus
	geom.Vec(-1, -1), // SPlus
	geom.Vec(0, -2),  // RMinus
	geom.Vec(1, -1),  // QPlus
}
