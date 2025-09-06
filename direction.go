package hex

import (
	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/ints"
)

// DirectionsFor returns the neighbors for the given coordinate.
func DirectionsFor(index ints.Point, coordType CoordinateType) []ints.Vector {
	parityQ := index.X & 1 // col
	parityR := index.Y & 1 // row

	switch coordType {
	case OffsetOddR:
		return DirectionsOffsetOddR[parityR][:]
	case OffsetEvenR:
		return DirectionsOffsetEvenR[parityR][:]
	case OffsetOddQ:
		return DirectionsOffsetOddQ[parityQ][:]
	case OffsetEvenQ:
		return DirectionsOffsetEvenQ[parityQ][:]
	case DoubleWidth:
		return DirectionsDoubleWidth[:]
	case DoubleHeight:
		return DirectionsDoubleHeight[:]
	case Axial:
		return Directions[:]
	default:
		return Directions[:]
	}
}

// DirectionFor returns the neighbor for the given coordinate in the given direction.
func DirectionFor(index ints.Point, coordType CoordinateType, direction Direction) ints.Vector {
	return DirectionsFor(index, coordType)[direction]
}

// Direction enum
type Direction int

// Direction for flat-top hexes (Axial, OffsetOddQ, OffsetEvenQ, DoubleHeight)
const (
	DirectionFlatSE Direction = 0
	DirectionFlatNE Direction = 1
	DirectionFlatN  Direction = 2
	DirectionFlatNW Direction = 3
	DirectionFlatSW Direction = 4
	DirectionFlatS  Direction = 5
)

// Direction for pointy-top hexes (Axial, OffsetOddR, OffsetEvenR, DoubleWidth)
const (
	DirectionPointyE  Direction = 0
	DirectionPointyNE Direction = 1
	DirectionPointyNW Direction = 2
	DirectionPointyW  Direction = 3
	DirectionPointySW Direction = 4
	DirectionPointySE Direction = 5
)

// Directions for axial (cube) coordinates
var Directions = [6]ints.Vector{
	geom.V(1, 0),
	geom.V(1, -1),
	geom.V(0, -1),
	geom.V(-1, 0),
	geom.V(-1, 1),
	geom.V(0, 1),
}

// Directions for offset odd-r coordinates (even row, odd row)
var DirectionsOffsetOddR = [2][6]ints.Vector{
	{
		geom.V(1, 0),
		geom.V(0, -1),
		geom.V(-1, -1),
		geom.V(-1, 0),
		geom.V(-1, 1),
		geom.V(0, 1),
	},
	{
		geom.V(1, 0),
		geom.V(1, -1),
		geom.V(0, -1),
		geom.V(-1, 0),
		geom.V(0, 1),
		geom.V(1, 1),
	},
}

// Directions for offset even-r coordinates (even row, odd row)
var DirectionsOffsetEvenR = [2][6]ints.Vector{
	{
		geom.V(1, 0),
		geom.V(1, -1),
		geom.V(0, -1),
		geom.V(-1, 0),
		geom.V(0, 1),
		geom.V(1, 1),
	},
	{
		geom.V(1, 0),
		geom.V(0, -1),
		geom.V(-1, -1),
		geom.V(-1, 0),
		geom.V(-1, 1),
		geom.V(0, 1),
	},
}

// Directions for offset odd-q coordinates (even col, odd col)
var DirectionsOffsetOddQ = [2][6]ints.Vector{
	{
		geom.V(1, 0),
		geom.V(1, -1),
		geom.V(0, -1),
		geom.V(-1, -1),
		geom.V(-1, 0),
		geom.V(0, 1),
	},
	{
		geom.V(1, 1),
		geom.V(1, 0),
		geom.V(0, -1),
		geom.V(-1, 0),
		geom.V(-1, 1),
		geom.V(0, 1),
	},
}

// Directions for offset even-q coordinates (even col, odd col)
var DirectionsOffsetEvenQ = [2][6]ints.Vector{
	{
		geom.V(1, 1),
		geom.V(1, 0),
		geom.V(0, -1),
		geom.V(-1, 0),
		geom.V(-1, 1),
		geom.V(0, 1),
	},
	{
		geom.V(1, 0),
		geom.V(1, -1),
		geom.V(0, -1),
		geom.V(-1, -1),
		geom.V(-1, 0),
		geom.V(0, 1),
	},
}

// Directions for double-width coordinates
var DirectionsDoubleWidth = [6]ints.Vector{
	geom.V(2, 0),
	geom.V(1, -1),
	geom.V(-1, -1),
	geom.V(-2, 0),
	geom.V(-1, 1),
	geom.V(1, 1),
}

// Directions for double-height coordinates
var DirectionsDoubleHeight = [6]ints.Vector{
	geom.V(1, 1),
	geom.V(1, -1),
	geom.V(0, -2),
	geom.V(-1, -1),
	geom.V(-1, 1),
	geom.V(0, 2),
}
