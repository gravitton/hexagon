package hex

import (
	"fmt"

	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/ints"
)

// Direction represents one of the six neighbor directions around a hex.
// The named constants follow cube coordinate axes:
// - Q+ increments q and compensates by decrementing r.
// - R+ increments r and compensates by decrementing s (-q-r).
// - S+ increments s (-q-r) and compensates by decrementing q.
type Direction int

const (
	SMinus Direction = iota // -S, flat-top SE, pointy-top E
	QPlus                   // +Q, flat-top NE, pointy-top NE
	RMinus                  // -R, flat-top N,  pointy-top NW
	SPlus                   // +S, flat-top NW, pointy-top W
	QMinus                  // -Q, flat-top SW, pointy-top SW
	RPlus                   // +R, flat-top S,  pointy-top SE
)

// Direction aliases for flat-top hexes (Axial, OffsetOddQ, OffsetEvenQ, DoubleHeight)
const (
	FlatTopSE = SMinus
	FlatTopNE = QPlus
	FlatTopN  = RMinus
	FlatTopNW = SPlus
	FlatTopSW = QMinus
	FlatTopS  = RPlus
)

// Direction aliases for pointy-top hexes (Axial, OffsetOddR, OffsetEvenR, DoubleWidth)
const (
	PointyTopE  = SMinus
	PointyTopNE = QPlus
	PointyTopNW = RMinus
	PointyTopW  = SPlus
	PointyTopSW = QMinus
	PointyTopSE = RPlus
)

// NeighborOffset returns the neighbor offset vector for the given direction.
func (d Direction) NeighborOffset() ints.Vector {
	return Directions[d%6]
}

// Opposite returns the direction directly opposite to d (rotated 180°, three steps away).
func (d Direction) Opposite() Direction {
	return Direction(mod6(int(d) + 3))
}

// Rotate advances d by steps positions counterclockwise in the standard mathematical
// sense. On typical screens with Y pointing down, positive steps appear clockwise.
// Negative steps go the other way. Rotate(3) is equivalent to Opposite().
func (d Direction) Rotate(steps int) Direction {
	return Direction(mod6(int(d) + steps))
}

// String returns the name of the direction constant.
func (d Direction) String() string {
	switch d % 6 {
	case SMinus:
		return "SMinus"
	case QPlus:
		return "QPlus"
	case RMinus:
		return "RMinus"
	case SPlus:
		return "SPlus"
	case QMinus:
		return "QMinus"
	case RPlus:
		return "RPlus"
	default:
		return fmt.Sprintf("Direction(%d)", int(d))
	}
}

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

// mod6 wraps n into [0, 6) correctly for negative values.
func mod6(n int) int {
	return ((n % 6) + 6) % 6
}

// Directions lists the 6 neighbor vectors for axial coordinates, ordered counterclockwise from the (south)-east direction.
var Directions = [6]ints.Vector{
	geom.Vec(1, 0),  // -S, flat-top SE, pointy-top E
	geom.Vec(1, -1), // +Q, flat-top NE, pointy-top NE
	geom.Vec(0, -1), // -R, flat-top N,  pointy-top NW
	geom.Vec(-1, 0), // +S, flat-top NW, pointy-top W
	geom.Vec(-1, 1), // -Q, flat-top SW, pointy-top SW
	geom.Vec(0, 1),  // +R, flat-top S,  pointy-top SE
}

// DirectionsOffsetOddR lists neighbor vectors for odd-r offset coordinates as
// [parityRow][direction], where parityRow=0 for even rows and 1 for odd rows.
var DirectionsOffsetOddR = [2][6]ints.Vector{
	{
		geom.Vec(1, 0),
		geom.Vec(0, -1),
		geom.Vec(-1, -1),
		geom.Vec(-1, 0),
		geom.Vec(-1, 1),
		geom.Vec(0, 1),
	},
	{
		geom.Vec(1, 0),
		geom.Vec(1, -1),
		geom.Vec(0, -1),
		geom.Vec(-1, 0),
		geom.Vec(0, 1),
		geom.Vec(1, 1),
	},
}

// DirectionsOffsetEvenR lists neighbor vectors for even-r offset coordinates as
// [parityRow][direction], where parityRow=0 for even rows and 1 for odd rows.
var DirectionsOffsetEvenR = [2][6]ints.Vector{
	{
		geom.Vec(1, 0),
		geom.Vec(1, -1),
		geom.Vec(0, -1),
		geom.Vec(-1, 0),
		geom.Vec(0, 1),
		geom.Vec(1, 1),
	},
	{
		geom.Vec(1, 0),
		geom.Vec(0, -1),
		geom.Vec(-1, -1),
		geom.Vec(-1, 0),
		geom.Vec(-1, 1),
		geom.Vec(0, 1),
	},
}

// DirectionsOffsetOddQ lists neighbor vectors for odd-q offset coordinates as
// [parityCol][direction], where parityCol=0 for even columns and 1 for odd columns.
var DirectionsOffsetOddQ = [2][6]ints.Vector{
	{
		geom.Vec(1, 0),
		geom.Vec(1, -1),
		geom.Vec(0, -1),
		geom.Vec(-1, -1),
		geom.Vec(-1, 0),
		geom.Vec(0, 1),
	},
	{
		geom.Vec(1, 1),
		geom.Vec(1, 0),
		geom.Vec(0, -1),
		geom.Vec(-1, 0),
		geom.Vec(-1, 1),
		geom.Vec(0, 1),
	},
}

// DirectionsOffsetEvenQ lists neighbor vectors for even-q offset coordinates as
// [parityCol][direction], where parityCol=0 for even columns and 1 for odd columns.
var DirectionsOffsetEvenQ = [2][6]ints.Vector{
	{
		geom.Vec(1, 1),
		geom.Vec(1, 0),
		geom.Vec(0, -1),
		geom.Vec(-1, 0),
		geom.Vec(-1, 1),
		geom.Vec(0, 1),
	},
	{
		geom.Vec(1, 0),
		geom.Vec(1, -1),
		geom.Vec(0, -1),
		geom.Vec(-1, -1),
		geom.Vec(-1, 0),
		geom.Vec(0, 1),
	},
}

// DirectionsDoubleWidth lists neighbor vectors for double-width coordinates.
var DirectionsDoubleWidth = [6]ints.Vector{
	geom.Vec(2, 0),
	geom.Vec(1, -1),
	geom.Vec(-1, -1),
	geom.Vec(-2, 0),
	geom.Vec(-1, 1),
	geom.Vec(1, 1),
}

// DirectionsDoubleHeight lists neighbor vectors for double-height coordinates.
var DirectionsDoubleHeight = [6]ints.Vector{
	geom.Vec(1, 1),
	geom.Vec(1, -1),
	geom.Vec(0, -2),
	geom.Vec(-1, -1),
	geom.Vec(-1, 1),
	geom.Vec(0, 2),
}
