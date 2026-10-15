package hex

import (
	"fmt"

	geom "github.com/gravitton/geometry"
)

// CoordinateSystem is a way of addressing the hexes of a grid by two integers. Every value outside
// the seven counts as CoordinateSystemNone.
type CoordinateSystem int

const (
	// Axial addresses a hex by its own q and r.
	Axial CoordinateSystem = iota

	// OffsetOddR addresses pointy-top rows, the odd ones shifted right by half a hex.
	OffsetOddR

	// OffsetEvenR addresses pointy-top rows, the even ones shifted right by half a hex.
	OffsetEvenR

	// OffsetOddQ addresses flat-top columns, the odd ones shifted down by half a hex.
	OffsetOddQ

	// OffsetEvenQ addresses flat-top columns, the even ones shifted down by half a hex.
	OffsetEvenQ

	// DoubleWidth addresses pointy-top rows with the column doubled (col = 2q+r, row = r), so
	// every cell has the same six steps. It addresses only the cells whose col and row share a
	// parity; From takes any other to the hex beside it on the side of lower q.
	DoubleWidth

	// DoubleHeight addresses flat-top columns with the row doubled (col = q, row = 2r+q), so every
	// cell has the same six steps. It addresses only the cells whose col and row share a parity;
	// From takes any other to the hex beside it on the side of lower r.
	DoubleHeight

	// CoordinateSystemNone is the absence of a coordinate system, for which Offsets, Offset,
	// Neighbor, Neighbors, To and From panic.
	CoordinateSystemNone CoordinateSystem = -1
)

// CoordinateSystems returns the seven coordinate systems in order.
func CoordinateSystems() [7]CoordinateSystem {
	return [7]CoordinateSystem{Axial, OffsetOddR, OffsetEvenR, OffsetOddQ, OffsetEvenQ, DoubleWidth, DoubleHeight}
}

// ParseCoordinateSystem returns the coordinate system with the name String prints, and an error for
// any other string.
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

// IsNone reports whether the system is none of the seven, so every other value counts as
// CoordinateSystemNone.
func (s CoordinateSystem) IsNone() bool {
	switch s {
	case Axial, OffsetOddR, OffsetEvenR, OffsetOddQ, OffsetEvenQ, DoubleWidth, DoubleHeight:
		return false
	default:
		return true
	}
}

// Offsets returns the step to each neighbor of the cell at index, in the order of Directions; in an
// offset system they depend on the parity of the row or column. It panics for a system outside the
// seven.
func (s CoordinateSystem) Offsets(index geom.Point[int]) [6]geom.Vector[int] {
	switch s {
	case Axial:
		return directionOffsets
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
	default:
		panic(fmt.Sprintf("hex: unknown coordinate system %d", s))
	}
}

// Offset returns the step to the neighbor of the cell at index in the direction, and the zero
// vector for DirectionNone. It panics for a system outside the seven, whatever the direction.
func (s CoordinateSystem) Offset(index geom.Point[int], direction Direction) geom.Vector[int] {
	offsets := s.Offsets(index)
	if direction.IsNone() {
		return geom.Vector[int]{}
	}

	return offsets[direction.normalize()]
}

// Neighbor returns the cell beside index in the direction, and index itself for DirectionNone. It
// panics for a system outside the seven, whatever the direction.
func (s CoordinateSystem) Neighbor(index geom.Point[int], direction Direction) geom.Point[int] {
	return index.Add(s.Offset(index, direction))
}

// Neighbors returns the six cells beside index in the order of Directions. It panics for a system
// outside the seven.
func (s CoordinateSystem) Neighbors(index geom.Point[int]) [6]geom.Point[int] {
	var neighbors [6]geom.Point[int]
	for i, offset := range s.Offsets(index) {
		neighbors[i] = index.Add(offset)
	}

	return neighbors
}

// To returns the coordinates of the hex in the system. It panics for a system outside the seven.
func (s CoordinateSystem) To(hex Hex) geom.Point[int] {
	switch s {
	case Axial:
		return toAxial(hex)
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
	default:
		panic(fmt.Sprintf("hex: unknown coordinate system %d", s))
	}
}

// From returns the hex at the coordinates in the system. It panics for a system outside the seven.
func (s CoordinateSystem) From(index geom.Point[int]) Hex {
	switch s {
	case Axial:
		return fromAxial(index)
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
	default:
		panic(fmt.Sprintf("hex: unknown coordinate system %d", s))
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

// MarshalText implements encoding.TextMarshaler with the name String prints.
func (s CoordinateSystem) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler through ParseCoordinateSystem.
func (s *CoordinateSystem) UnmarshalText(text []byte) error {
	system, err := ParseCoordinateSystem(string(text))
	if err != nil {
		return err
	}

	*s = system

	return nil
}

// toAxial returns the axial coordinates of the hex.
func toAxial(hex Hex) geom.Point[int] {
	return hex.Point()
}

// fromAxial returns the hex at the axial coordinates.
func fromAxial(index geom.Point[int]) Hex {
	return Hex{index.X, index.Y}
}

// toOffsetOddR returns the odd-r coordinates of the hex.
func toOffsetOddR(hex Hex) geom.Point[int] {
	parity := hex.R & 1

	col := hex.Q + (hex.R-parity)/2
	row := hex.R

	return geom.Pt(col, row)
}

// fromOffsetOddR returns the hex at the odd-r coordinates.
func fromOffsetOddR(index geom.Point[int]) Hex {
	parity := index.Y & 1

	q := index.X - (index.Y-parity)/2
	r := index.Y

	return Hex{q, r}
}

// toOffsetEvenR returns the even-r coordinates of the hex.
func toOffsetEvenR(hex Hex) geom.Point[int] {
	parity := hex.R & 1

	col := hex.Q + (hex.R+parity)/2
	row := hex.R

	return geom.Pt(col, row)
}

// fromOffsetEvenR returns the hex at the even-r coordinates.
func fromOffsetEvenR(index geom.Point[int]) Hex {
	parity := index.Y & 1

	q := index.X - (index.Y+parity)/2
	r := index.Y

	return Hex{q, r}
}

// toOffsetOddQ returns the odd-q coordinates of the hex.
func toOffsetOddQ(hex Hex) geom.Point[int] {
	parity := hex.Q & 1

	col := hex.Q
	row := hex.R + (hex.Q-parity)/2

	return geom.Pt(col, row)
}

// fromOffsetOddQ returns the hex at the odd-q coordinates.
func fromOffsetOddQ(index geom.Point[int]) Hex {
	parity := index.X & 1

	q := index.X
	r := index.Y - (index.X-parity)/2

	return Hex{q, r}
}

// toOffsetEvenQ returns the even-q coordinates of the hex.
func toOffsetEvenQ(hex Hex) geom.Point[int] {
	parity := hex.Q & 1

	col := hex.Q
	row := hex.R + (hex.Q+parity)/2

	return geom.Pt(col, row)
}

// fromOffsetEvenQ returns the hex at the even-q coordinates.
func fromOffsetEvenQ(index geom.Point[int]) Hex {
	parity := index.X & 1

	q := index.X
	r := index.Y - (index.X+parity)/2

	return Hex{q, r}
}

// toDoubleWidth returns the double-width coordinates of the hex.
func toDoubleWidth(hex Hex) geom.Point[int] {
	col := 2*hex.Q + hex.R
	row := hex.R

	return geom.Pt(col, row)
}

// fromDoubleWidth returns the hex at the double-width coordinates, flooring a mixed parity.
func fromDoubleWidth(index geom.Point[int]) Hex {
	q := (index.X - index.Y) >> 1
	r := index.Y

	return Hex{q, r}
}

// toDoubleHeight returns the double-height coordinates of the hex.
func toDoubleHeight(hex Hex) geom.Point[int] {
	col := hex.Q
	row := 2*hex.R + hex.Q

	return geom.Pt(col, row)
}

// fromDoubleHeight returns the hex at the double-height coordinates, flooring a mixed parity.
func fromDoubleHeight(index geom.Point[int]) Hex {
	q := index.X
	r := (index.Y - index.X) >> 1

	return Hex{q, r}
}

// offsetsOddR holds the odd-r steps to each neighbor, by row parity and then direction.
var offsetsOddR = [2][6]geom.Vector[int]{
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

// offsetsEvenR holds the even-r steps to each neighbor, by row parity and then direction.
var offsetsEvenR = [2][6]geom.Vector[int]{
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

// offsetsOddQ holds the odd-q steps to each neighbor, by column parity and then direction.
var offsetsOddQ = [2][6]geom.Vector[int]{
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

// offsetsEvenQ holds the even-q steps to each neighbor, by column parity and then direction.
var offsetsEvenQ = [2][6]geom.Vector[int]{
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

// offsetsDoubleWidth holds the double-width step to each neighbor, by direction.
var offsetsDoubleWidth = [6]geom.Vector[int]{
	geom.Vec(2, 0),   // SMinus
	geom.Vec(1, 1),   // RPlus
	geom.Vec(-1, 1),  // QMinus
	geom.Vec(-2, 0),  // SPlus
	geom.Vec(-1, -1), // RMinus
	geom.Vec(1, -1),  // QPlus
}

// offsetsDoubleHeight holds the double-height step to each neighbor, by direction.
var offsetsDoubleHeight = [6]geom.Vector[int]{
	geom.Vec(1, 1),   // SMinus
	geom.Vec(0, 2),   // RPlus
	geom.Vec(-1, 1),  // QMinus
	geom.Vec(-1, -1), // SPlus
	geom.Vec(0, -2),  // RMinus
	geom.Vec(1, -1),  // QPlus
}
