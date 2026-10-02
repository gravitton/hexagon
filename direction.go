package hex

import (
	"fmt"
	"math"

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
// counterclockwise on screen. Any other value wraps into [SMinus, QPlus]; only
// DirectionNone stands outside the six.
type Direction int

const (
	SMinus Direction = iota // -S, flat-top SE, pointy-top E,  angle 0
	RPlus                   // +R, flat-top S,  pointy-top SE, angle π/3
	QMinus                  // -Q, flat-top SW, pointy-top SW, angle 2π/3
	SPlus                   // +S, flat-top NW, pointy-top W,  angle π
	RMinus                  // -R, flat-top N,  pointy-top NW, angle 4π/3
	QPlus                   // +Q, flat-top NE, pointy-top NE, angle 5π/3

	// DirectionNone is the absence of a direction. It names no neighbor, so Offset and Hex
	// give a zero step for it and Angle has no angle to return.
	DirectionNone Direction = -1
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

// directionOffsets lists the axial neighbor vector of each direction, indexed by direction.
var directionOffsets = [6]ints.Vector{
	geom.Vec(1, 0),  // -S, flat-top SE, pointy-top E
	geom.Vec(0, 1),  // +R, flat-top S,  pointy-top SE
	geom.Vec(-1, 1), // -Q, flat-top SW, pointy-top SW
	geom.Vec(-1, 0), // +S, flat-top NW, pointy-top W
	geom.Vec(0, -1), // -R, flat-top N,  pointy-top NW
	geom.Vec(1, -1), // +Q, flat-top NE, pointy-top NE
}

// Directions lists the six directions ordered by increasing angle from SMinus.
// It returns a fresh array, so a caller cannot alter the list.
func Directions() [6]Direction {
	return [6]Direction{SMinus, RPlus, QMinus, SPlus, RMinus, QPlus}
}

// DirectionFromAngle returns the direction nearest to the given angle in radians,
// or DirectionNone for NaN and ±Inf.
func DirectionFromAngle(angle float64) Direction {
	if math.IsNaN(angle) || math.IsInf(angle, 0) {
		return DirectionNone
	}

	return geom.Mod(Direction(math.Round(geom.NormalizeAngle(angle)/(geom.Pi/3))), 6)
}

// ParseDirection returns the direction with the given name, as String prints it, and an error
// for any other string. "None" parses to DirectionNone.
func ParseDirection(name string) (Direction, error) {
	if name == DirectionNone.String() {
		return DirectionNone, nil
	}

	for _, direction := range Directions() {
		if direction.String() == name {
			return direction, nil
		}
	}

	return DirectionNone, fmt.Errorf("hex: unknown direction %q", name)
}

// Opposite returns the direction directly opposite to d (rotated 180°, three steps away).
func (d Direction) Opposite() Direction {
	return d.Turn(3)
}

// Turn advances d by steps sixths of a turn of increasing angle, the same sense as a
// positive geom.Direction step: counterclockwise in math coordinates, clockwise as drawn
// on a screen with Y pointing down. Negative steps go the other way.
// Turn(3) is equivalent to Opposite(), and DirectionNone turns to itself.
func (d Direction) Turn(steps int) Direction {
	if d.IsNone() {
		return DirectionNone
	}

	return geom.Mod(d+Direction(steps), 6)
}

// Offset returns the axial neighbor offset vector for the given direction,
// and the zero vector for DirectionNone.
func (d Direction) Offset() ints.Vector {
	if d.IsNone() {
		return ints.Vector{}
	}

	return directionOffsets[d.normalize()]
}

// Hex returns the unit Hex step in the direction, the [Hex] counterpart of [Direction.Offset],
// and the zero hex for DirectionNone.
func (d Direction) Hex() Hex {
	offset := d.Offset()

	return Hex{offset.X, offset.Y}
}

// Angle returns the angle of the direction in radians, in [0, 2π), the order the constants
// follow: SMinus is 0 and QPlus is 5π/3. It is NaN for DirectionNone, which has no angle.
// Angle and [DirectionFromAngle] round-trip for every direction, DirectionNone included.
func (d Direction) Angle() float64 {
	if d.IsNone() {
		return math.NaN()
	}

	return float64(d.normalize()) * geom.Pi / 3
}

// normalize wraps d into [SMinus, QPlus], correctly for negative values, and leaves
// DirectionNone alone.
func (d Direction) normalize() Direction {
	if d.IsNone() {
		return DirectionNone
	}

	return geom.Mod(d, 6)
}

// IsNone reports whether the direction is DirectionNone.
func (d Direction) IsNone() bool {
	return d == DirectionNone
}

// String returns the name of the direction constant.
func (d Direction) String() string {
	switch d.normalize() {
	case SMinus:
		return "SMinus"
	case RPlus:
		return "RPlus"
	case QMinus:
		return "QMinus"
	case SPlus:
		return "SPlus"
	case RMinus:
		return "RMinus"
	case QPlus:
		return "QPlus"
	default:
		return "None"
	}
}

// MarshalText implements encoding.TextMarshaler with the name String prints, so a direction
// is stored as "QPlus" in JSON and as a map key rather than as its number.
func (d Direction) MarshalText() ([]byte, error) {
	return []byte(d.String()), nil
}

// UnmarshalText implements encoding.TextUnmarshaler, the inverse of MarshalText through ParseDirection.
func (d *Direction) UnmarshalText(text []byte) error {
	direction, err := ParseDirection(string(text))
	if err != nil {
		return err
	}

	*d = direction

	return nil
}
