package hex

import (
	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/ints"
)

// CoordinateSystem enumerates supported hexagonal grid coordinate systems.
type CoordinateSystem int

const (
	Axial        CoordinateSystem = iota // Axial (Cube) coordinates (q, r, s)
	OffsetOddR                           // Odd rows are offset (pointy-top hexes)
	OffsetEvenR                          // Even rows are offset (pointy-top hexes)
	OffsetOddQ                           // Odd columns are offset (flat-top hexes)
	OffsetEvenQ                          // Even columns are offset (flat-top hexes)
	DoubleWidth                          // Double cols (pointy-top hexes)
	DoubleHeight                         // Double rows (flat-top hexes)
)

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
	return hex.ToPoint()
}

// FromAxial converts an ints.Point (q,r) into an axial Hex.
func FromAxial(index ints.Point) Hex {
	return Hex{index.X, index.Y}
}

// FromPoint is alias to FromAxial.
func FromPoint(index ints.Point) Hex {
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
