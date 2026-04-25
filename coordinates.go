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
	DoubleWidth

	// DoubleHeight is a flat-top system that doubles the row axis: col = q, row = 2r+q.
	// All six neighbors are reachable with fixed offsets — no per-cell parity check needed.
	DoubleHeight
)

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
		return fmt.Sprintf("CoordinateSystem(%d)", int(s))
	}
}

// To converts an axial hex to the given coordinate system as an ints.Point.
func To(hex Hex, system CoordinateSystem) ints.Point {
	switch system {
	case OffsetOddR:
		return ToOffsetOddR(hex)
	case OffsetEvenR:
		return ToOffsetEvenR(hex)
	case OffsetOddQ:
		return ToOffsetOddQ(hex)
	case OffsetEvenQ:
		return ToOffsetEvenQ(hex)
	case DoubleWidth:
		return ToDoubleWidth(hex)
	case DoubleHeight:
		return ToDoubleHeight(hex)
	case Axial:
		return ToAxial(hex)
	default:
		panic("unsupported coordinate system")
	}
}

// From converts a coordinate in the given system into an axial Hex.
func From(index ints.Point, system CoordinateSystem) Hex {
	switch system {
	case OffsetOddR:
		return FromOffsetOddR(index)
	case OffsetEvenR:
		return FromOffsetEvenR(index)
	case OffsetOddQ:
		return FromOffsetOddQ(index)
	case OffsetEvenQ:
		return FromOffsetEvenQ(index)
	case DoubleWidth:
		return FromDoubleWidth(index)
	case DoubleHeight:
		return FromDoubleHeight(index)
	case Axial:
		return FromAxial(index)
	default:
		panic("unsupported coordinate system")
	}
}

// ToAxial returns the axial (q,r) as an ints.Point.
func ToAxial(hex Hex) ints.Point {
	return geom.Pt(hex.Q, hex.R)
}

// FromAxial converts an ints.Point (q,r) into an axial Hex.
func FromAxial(index ints.Point) Hex {
	return Hex{index.X, index.Y}
}

// ToOffsetOddR converts axial to odd-r offset coordinates.
// Odd rows are shifted right by +1/2 column.
func ToOffsetOddR(hex Hex) ints.Point {
	parity := hex.R & 1

	col := hex.Q + (hex.R-parity)/2
	row := hex.R

	return geom.Pt(col, row)
}

// FromOffsetOddR converts an odd-r offset coordinate to axial.
func FromOffsetOddR(index ints.Point) Hex {
	parity := index.Y & 1

	q := index.X - (index.Y-parity)/2
	r := index.Y

	return Hex{q, r}
}

// ToOffsetEvenR converts axial to even-r offset coordinates.
// Even rows are shifted right by +1/2 column.
func ToOffsetEvenR(hex Hex) ints.Point {
	parity := hex.R & 1

	col := hex.Q + (hex.R+parity)/2
	row := hex.R

	return geom.Pt(col, row)
}

// FromOffsetEvenR converts an even-r offset coordinate to axial.
func FromOffsetEvenR(index ints.Point) Hex {
	parity := index.Y & 1

	q := index.X - (index.Y+parity)/2
	r := index.Y

	return Hex{q, r}
}

// ToOffsetOddQ converts axial to odd-q offset coordinates.
// Odd columns are shifted down by +1/2 row.
func ToOffsetOddQ(hex Hex) ints.Point {
	parity := hex.Q & 1

	col := hex.Q
	row := hex.R + (hex.Q-parity)/2

	return geom.Pt(col, row)
}

// FromOffsetOddQ converts an odd-q offset coordinate to axial.
func FromOffsetOddQ(index ints.Point) Hex {
	parity := index.X & 1

	q := index.X
	r := index.Y - (index.X-parity)/2

	return Hex{q, r}
}

// ToOffsetEvenQ converts axial to even-q offset coordinates.
// Even columns are shifted down by +1/2 row.
func ToOffsetEvenQ(hex Hex) ints.Point {
	parity := hex.Q & 1

	col := hex.Q
	row := hex.R + (hex.Q+parity)/2

	return geom.Pt(col, row)
}

// FromOffsetEvenQ converts an even-q offset coordinate to axial.
func FromOffsetEvenQ(index ints.Point) Hex {
	parity := index.X & 1

	q := index.X
	r := index.Y - (index.X+parity)/2

	return Hex{q, r}
}

// ToDoubleWidth converts axial to double-width coordinates (doubling the q axis).
func ToDoubleWidth(hex Hex) ints.Point {
	col := 2*hex.Q + hex.R
	row := hex.R

	return geom.Pt(col, row)
}

// FromDoubleWidth converts a double-width coordinate to axial.
func FromDoubleWidth(index ints.Point) Hex {
	q := (index.X - index.Y) / 2
	r := index.Y

	return Hex{q, r}
}

// ToDoubleHeight converts axial to double-height coordinates (doubling the r axis).
func ToDoubleHeight(hex Hex) ints.Point {
	col := hex.Q
	row := 2*hex.R + hex.Q

	return geom.Pt(col, row)
}

// FromDoubleHeight converts a double-height coordinate to axial.
func FromDoubleHeight(index ints.Point) Hex {
	q := index.X
	r := (index.Y - index.X) / 2

	return Hex{q, r}
}
