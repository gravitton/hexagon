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
//
// Directions are numbered by increasing angle in pixel space, matching geom.Direction:
// counterclockwise in the standard math convention where Y grows upward, which appears
// clockwise as drawn on a screen with Y pointing down. A negative step is therefore
// counterclockwise on screen.
type Direction int

const (
	SMinus Direction = iota // -S, flat-top SE, pointy-top E
	RPlus                   // +R, flat-top S,  pointy-top SE
	QMinus                  // -Q, flat-top SW, pointy-top SW
	SPlus                   // +S, flat-top NW, pointy-top W
	RMinus                  // -R, flat-top N,  pointy-top NW
	QPlus                   // +Q, flat-top NE, pointy-top NE
)

// Direction aliases for flat-top hexes (Axial, OffsetOddQ, OffsetEvenQ, DoubleHeight)
const (
	FlatTopSouthEast = SMinus
	FlatTopSouth     = RPlus
	FlatTopSouthWest = QMinus
	FlatTopNorthWest = SPlus
	FlatTopNorth     = RMinus
	FlatTopNorthEast = QPlus
)

// Direction aliases for pointy-top hexes (Axial, OffsetOddR, OffsetEvenR, DoubleWidth)
const (
	PointyTopEast      = SMinus
	PointyTopSouthEast = RPlus
	PointyTopSouthWest = QMinus
	PointyTopWest      = SPlus
	PointyTopNorthWest = RMinus
	PointyTopNorthEast = QPlus
)

// Directions lists the 6 directions ordered by increasing angle from the (south)-east direction.
var Directions = [6]Direction{SMinus, RPlus, QMinus, SPlus, RMinus, QPlus}

// directionOffsets lists the axial neighbor vector of each direction, indexed by direction.
var directionOffsets = [6]ints.Vector{
	geom.Vec(1, 0),  // -S, flat-top SE, pointy-top E
	geom.Vec(0, 1),  // +R, flat-top S,  pointy-top SE
	geom.Vec(-1, 1), // -Q, flat-top SW, pointy-top SW
	geom.Vec(-1, 0), // +S, flat-top NW, pointy-top W
	geom.Vec(0, -1), // -R, flat-top N,  pointy-top NW
	geom.Vec(1, -1), // +Q, flat-top NE, pointy-top NE
}

// Offset returns the neighbor offset vector for the given direction.
func (d Direction) Offset() ints.Vector {
	return directionOffsets[d.normalize()]
}

// Hex returns the unit Hex step in the direction, the [Hex] counterpart of [Direction.Offset].
func (d Direction) Hex() Hex {
	return Pt(d.Offset().XY())
}

// Opposite returns the direction directly opposite to d (rotated 180°, three steps away).
func (d Direction) Opposite() Direction {
	return d.Rotate(3)
}

// Rotate advances d by steps sixths of a turn of increasing angle, the same sense as a
// positive geom.Direction step: counterclockwise in math coordinates, clockwise as drawn
// on a screen with Y pointing down. Negative steps go the other way.
// Rotate(3) is equivalent to Opposite().
func (d Direction) Rotate(steps int) Direction {
	return geom.Mod(d+Direction(steps), 6)
}

// String returns the name of the direction constant.
func (d Direction) String() string {
	switch d.normalize() {
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

// normalize wraps d into [SMinus, QPlus], correctly for negative values.
func (d Direction) normalize() Direction {
	return geom.Mod(d, 6)
}
