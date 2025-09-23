package hex

import (
	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/ints"
)

// Direction represents one of the six neighbor directions around a hex.
//
// Constants for the directions from a Hex.
// - Q+ increments q and compensates by decrementing r.
// - R+ increments r and compensates by decrementing s (-q-r).
// - S+ increments s (-q-r) compensates by decrementing q.
type Direction int

// NeighborOffset returns the neighbor offset vector for the given direction.
func (d Direction) NeighborOffset() ints.Vector {
	return Directions[d%6]
}

const (
	DirectionSMinus Direction = iota // -S, flat-top SE, pointy-top E
	DirectionQPlus                   // +Q, flat-top NE, pointy-top NE
	DirectionRMinus                  // -R, flat-top N,  pointy-top NW
	DirectionSPlus                   // +S, flat-top NW, pointy-top W
	DirectionQMinus                  // -Q, flat-top SW, pointy-top SW
	DirectionRPlus                   // +R, flat-top S,  pointy-top SE
)

// Direction aliases for flat-top hexes (Axial, OffsetOddQ, OffsetEvenQ, DoubleHeight)
const (
	DirectionFlatTopSE = DirectionSMinus
	DirectionFlatTopNE = DirectionQPlus
	DirectionFlatTopN  = DirectionRMinus
	DirectionFlatTopNW = DirectionSPlus
	DirectionFlatTopSW = DirectionQMinus
	DirectionFlatTopS  = DirectionRPlus
)

// Direction aliases for pointy-top hexes (Axial, OffsetOddR, OffsetEvenR, DoubleWidth)
const (
	DirectionPointyTopE  = DirectionSMinus
	DirectionPointyTopNE = DirectionQPlus
	DirectionPointyTopNW = DirectionRMinus
	DirectionPointyTopW  = DirectionSPlus
	DirectionPointyTopSW = DirectionQMinus
	DirectionPointyTopSE = DirectionRPlus
)

// NeighborOffsets returns the 6 neighbor offsets for the given coordinate index
// in the specified coordinate system. For offset systems this accounts for
// row/column parity when determining neighbor vectors.
func NeighborOffsets(index ints.Point, system CoordinateSystem) []ints.Vector {
	switch system {
	case OffsetOddR:
		return NeighborOffsetsOffsetOddR(index)
	case OffsetEvenR:
		return NeighborOffsetsOffsetEvenR(index)
	case OffsetOddQ:
		return NeighborOffsetsOffsetOddQ(index)
	case OffsetEvenQ:
		return NeighborOffsetsOffsetEvenQ(index)
	case DoubleWidth:
		return NeighborOffsetsDoubleWidth()
	case DoubleHeight:
		return NeighborOffsetsDoubleHeight()
	case Axial:
		return NeighborOffsetsAxial()
	default:
		panic("unsupported coordinate system")
	}
}

// NeighborOffset returns the neighbor offset vector for the given coordinate
// in the specified direction and coordinate system.
func NeighborOffset(index ints.Point, system CoordinateSystem, direction Direction) ints.Vector {
	return NeighborOffsets(index, system)[direction]
}

// NeighborOffsetsAxial returns the neighbor offsets in axial (cube) coordinate system.
func NeighborOffsetsAxial() []ints.Vector {
	return Directions[:]
}

// NeighborOffsetsOffsetOddR returns the neighbor offsets in odd-r offset coordinate system.
func NeighborOffsetsOffsetOddR(index ints.Point) []ints.Vector {
	parity := index.Y & 1

	return DirectionsOffsetOddR[parity][:]
}

// NeighborOffsetsOffsetEvenR returns the neighbor offsets in even-r offset coordinate system.
func NeighborOffsetsOffsetEvenR(index ints.Point) []ints.Vector {
	parity := index.Y & 1

	return DirectionsOffsetEvenR[parity][:]
}

// NeighborOffsetsOffsetOddQ returns the neighbor offsets in odd-q offset coordinate system.
func NeighborOffsetsOffsetOddQ(index ints.Point) []ints.Vector {
	parity := index.X & 1

	return DirectionsOffsetOddQ[parity][:]
}

// NeighborOffsetsOffsetEvenQ returns the neighbor offsets in even-q offset coordinate system.
func NeighborOffsetsOffsetEvenQ(index ints.Point) []ints.Vector {
	parity := index.X & 1

	return DirectionsOffsetEvenQ[parity][:]
}

// NeighborOffsetsDoubleWidth returns the neighbor offsets in double-width coordinate system.
func NeighborOffsetsDoubleWidth() []ints.Vector {
	return DirectionsDoubleWidth[:]
}

// NeighborOffsetsDoubleHeight returns the neighbor offsets in double-height coordinate system.
func NeighborOffsetsDoubleHeight() []ints.Vector {
	return DirectionsDoubleHeight[:]
}

// Directions lists neighbor vectors for axial coordinates in counter-clockwise, starting at (south)-east direction.
var Directions = [6]ints.Vector{
	geom.V(1, 0),  // -S, flat-top SE, pointy-top E
	geom.V(1, -1), // +Q, flat-top NE, pointy-top NE
	geom.V(0, -1), // -R, flat-top N,  pointy-top NW
	geom.V(-1, 0), // +S, flat-top NW, pointy-top W
	geom.V(-1, 1), // -Q, flat-top SW, pointy-top SW
	geom.V(0, 1),  // +R, flat-top S,  pointy-top SE
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

// DirectionsDoubleWidth lists neighbor vectors for double-width coordinates.
var DirectionsDoubleWidth = [6]ints.Vector{
	geom.V(2, 0),
	geom.V(1, -1),
	geom.V(-1, -1),
	geom.V(-2, 0),
	geom.V(-1, 1),
	geom.V(1, 1),
}

// DirectionsDoubleHeight lists neighbor vectors for double-height coordinates.
var DirectionsDoubleHeight = [6]ints.Vector{
	geom.V(1, 1),
	geom.V(1, -1),
	geom.V(0, -2),
	geom.V(-1, -1),
	geom.V(-1, 1),
	geom.V(0, 2),
}
