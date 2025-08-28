package hex

import (
	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/ints"
)

// CoordinateType is a type of hexagonal grid coordinates.
type CoordinateType int

const (
	Axial        CoordinateType = iota // Axial (Cube) coordinates (q, r)
	OffsetOddR                         // Odd rows are offset (pointy-top hexes)
	OffsetEvenR                        // Even rows are offset (pointy-top hexes)
	OffsetOddQ                         // Odd columns are offset (flat-top hexes)
	OffsetEvenQ                        // Even columns are offset (flat-top hexes)
	DoubleWidth                        // Double cols (pointy-top hexes)
	DoubleHeight                       // Double rows (flat-top hexes)
)

func HexTo(hex Hex, coordType CoordinateType) ints.Point {
	switch coordType {
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
		return hex.ToPoint()
	default:
		return hex.ToPoint()
	}
}

func HexFrom(index ints.Point, coordType CoordinateType) Hex {
	switch coordType {
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
		return FromPoint(index)
	default:
		return FromPoint(index)
	}
}

func FromPoint(index ints.Point) Hex {
	return Hex{index.X, index.Y}
}

// shoves odd rows by +1/2 column
func ToOffsetOddR(hex Hex) ints.Point {
	parity := hex.R & 1

	col := hex.Q + (hex.R-parity)/2
	row := hex.R

	return geom.P(col, row)
}

func FromOffsetOddR(index ints.Point) Hex {
	parity := index.Y & 1

	q := index.X - (index.Y-parity)/2
	r := index.Y

	return Hex{q, r}
}

// shoves even rows by +1/2 column
func ToOffsetEvenR(hex Hex) ints.Point {
	parity := hex.R & 1

	col := hex.Q + (hex.R+parity)/2
	row := hex.R

	return geom.P(col, row)
}

func FromOffsetEvenR(index ints.Point) Hex {
	parity := index.Y & 1

	q := index.X - (index.Y+parity)/2
	r := index.Y

	return Hex{q, r}
}

// shoves odd columns by +1/2 row
func ToOffsetOddQ(hex Hex) ints.Point {
	parity := hex.Q & 1

	col := hex.Q
	row := hex.R + (hex.Q-parity)/2

	return geom.P(col, row)
}

func FromOffsetOddQ(index ints.Point) Hex {
	parity := index.X & 1

	q := index.X
	r := index.Y - (index.X-parity)/2

	return Hex{q, r}
}

// shoves even columns by +1/2 row
func ToOffsetEvenQ(hex Hex) ints.Point {
	parity := hex.Q & 1

	col := hex.Q
	row := hex.R + (hex.Q+parity)/2

	return geom.P(col, row)
}

func FromOffsetEvenQ(index ints.Point) Hex {
	parity := index.X & 1

	q := index.X
	r := index.Y - (index.X+parity)/2

	return Hex{q, r}
}

func ToDoubleWidth(hex Hex) ints.Point {
	col := 2*hex.Q + hex.R
	row := hex.R

	return geom.P(col, row)
}

func FromDoubleWidth(index ints.Point) Hex {
	q := (index.X - index.Y) / 2
	r := index.Y

	return Hex{q, r}
}

func ToDoubleHeight(hex Hex) ints.Point {
	col := hex.Q
	row := 2*hex.R + hex.Q

	return geom.P(col, row)
}

func FromDoubleHeight(index ints.Point) Hex {
	q := index.X
	r := (index.Y - index.X) / 2.0

	return Hex{q, r}
}
