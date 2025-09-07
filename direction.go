package hex

import (
	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/ints"
)

// NeighborOffsets returns the 6 neighbor offsets for the given coordinate index
// in the specified coordinate system. For offset systems this accounts for
// row/column parity when determining neighbor vectors.
func NeighborOffsets(index ints.Point, coordType CoordinateType) []ints.Vector {
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

// NeighborOffset returns the neighbor offset vector for the given coordinate
// in the specified direction and coordinate system.
func NeighborOffset(index ints.Point, coordType CoordinateType, direction Direction) ints.Vector {
	return NeighborOffsets(index, coordType)[direction]
}

// Direction represents one of the six neighbor directions around a hex.
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

// Directions are the 6 axial neighbor vectors in clockwise order, starting at E/SE.
var Directions = [6]ints.Vector{
	geom.V(1, 0),
	geom.V(1, -1),
	geom.V(0, -1),
	geom.V(-1, 0),
	geom.V(-1, 1),
	geom.V(0, 1),
}

// DirectionsOffsetOddR lists neighbor vectors for odd-r offset coordinates as
// [parityRow][direction], where parityRow=0 for even rows and 1 for odd rows.
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

// DirectionsOffsetEvenR lists neighbor vectors for even-r offset coordinates as
// [parityRow][direction], where parityRow=0 for even rows and 1 for odd rows.
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

// DirectionsOffsetOddQ lists neighbor vectors for odd-q offset coordinates as
// [parityCol][direction], where parityCol=0 for even columns and 1 for odd columns.
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

// DirectionsOffsetEvenQ lists neighbor vectors for even-q offset coordinates as
// [parityCol][direction], where parityCol=0 for even columns and 1 for odd columns.
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

// DirectionsDoubleWidth are neighbor vectors for double-width coordinates.
var DirectionsDoubleWidth = [6]ints.Vector{
	geom.V(2, 0),
	geom.V(1, -1),
	geom.V(-1, -1),
	geom.V(-2, 0),
	geom.V(-1, 1),
	geom.V(1, 1),
}

// DirectionsDoubleHeight are neighbor vectors for double-height coordinates.
var DirectionsDoubleHeight = [6]ints.Vector{
	geom.V(1, 1),
	geom.V(1, -1),
	geom.V(0, -2),
	geom.V(-1, -1),
	geom.V(-1, 1),
	geom.V(0, 2),
}
