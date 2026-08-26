package hex

import (
	"math"
	"testing"

	"github.com/gravitton/assert"
	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/ints"
)

var axialDirection = [6]ints.Vector{{X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 1}, {X: -1, Y: 0}, {X: 0, Y: -1}, {X: 1, Y: -1}}
var offsetOddRDirectionOddRow = [6]ints.Vector{{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: 1}, {X: -1, Y: 0}, {X: 0, Y: -1}, {X: 1, Y: -1}}
var offsetOddRDirectionEvenRow = [6]ints.Vector{{X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 1}, {X: -1, Y: 0}, {X: -1, Y: -1}, {X: 0, Y: -1}}
var offsetEvenRDirectionOddRow = [6]ints.Vector{{X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 1}, {X: -1, Y: 0}, {X: -1, Y: -1}, {X: 0, Y: -1}}
var offsetEvenRDirectionEvenRow = [6]ints.Vector{{X: 1, Y: 0}, {X: 1, Y: 1}, {X: 0, Y: 1}, {X: -1, Y: 0}, {X: 0, Y: -1}, {X: 1, Y: -1}}
var offsetOddQDirectionOddCol = [6]ints.Vector{{X: 1, Y: 1}, {X: 0, Y: 1}, {X: -1, Y: 1}, {X: -1, Y: 0}, {X: 0, Y: -1}, {X: 1, Y: 0}}
var offsetOddQDirectionEvenCol = [6]ints.Vector{{X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 0}, {X: -1, Y: -1}, {X: 0, Y: -1}, {X: 1, Y: -1}}
var offsetEvenQDirectionOddCol = [6]ints.Vector{{X: 1, Y: 0}, {X: 0, Y: 1}, {X: -1, Y: 0}, {X: -1, Y: -1}, {X: 0, Y: -1}, {X: 1, Y: -1}}
var offsetEvenQDirectionEvenCol = [6]ints.Vector{{X: 1, Y: 1}, {X: 0, Y: 1}, {X: -1, Y: 1}, {X: -1, Y: 0}, {X: 0, Y: -1}, {X: 1, Y: 0}}
var doubleWidthDirection = [6]ints.Vector{{X: 2, Y: 0}, {X: 1, Y: 1}, {X: -1, Y: 1}, {X: -2, Y: 0}, {X: -1, Y: -1}, {X: 1, Y: -1}}
var doubleHeightDirection = [6]ints.Vector{{X: 1, Y: 1}, {X: 0, Y: 2}, {X: -1, Y: 1}, {X: -1, Y: -1}, {X: 0, Y: -2}, {X: 1, Y: -1}}

func TestDirection_Offsets(t *testing.T) {
	assert.Equal(t, SMinus.Offset(), axialDirection[0])
	assert.Equal(t, RPlus.Offset(), axialDirection[1])
	assert.Equal(t, QMinus.Offset(), axialDirection[2])
	assert.Equal(t, SPlus.Offset(), axialDirection[3])
	assert.Equal(t, RMinus.Offset(), axialDirection[4])
	assert.Equal(t, QPlus.Offset(), axialDirection[5])

	// out-of-range directions wrap, negatives included
	assert.Equal(t, Direction(6).Offset(), axialDirection[0])
	assert.Equal(t, Direction(8).Offset(), axialDirection[2])
	assert.Equal(t, Direction(15).Offset(), axialDirection[3])
	assert.Equal(t, Direction(-1).Offset(), axialDirection[5])
	assert.Equal(t, Direction(-8).Offset(), axialDirection[4])

	assert.Equal(t, FlatTopSouthEast, SMinus)
	assert.Equal(t, FlatTopNorthEast, QPlus)
	assert.Equal(t, FlatTopNorth, RMinus)
	assert.Equal(t, FlatTopNorthWest, SPlus)
	assert.Equal(t, FlatTopSouthWest, QMinus)
	assert.Equal(t, FlatTopSouth, RPlus)

	assert.Equal(t, PointyTopEast, SMinus)
	assert.Equal(t, PointyTopNorthEast, QPlus)
	assert.Equal(t, PointyTopNorthWest, RMinus)
	assert.Equal(t, PointyTopWest, SPlus)
	assert.Equal(t, PointyTopSouthWest, QMinus)
	assert.Equal(t, PointyTopSouthEast, RPlus)
}

// TestDirection_AngleConvention locks the direction order to increasing angle in pixel space,
// matching geom.Direction. It uses the pointy-top mapping; flat-top differs only by a constant
// rotation, so the ordering is the same either way.
func TestDirection_AngleConvention(t *testing.T) {
	for i, direction := range Directions {
		v := direction.Offset()
		x := geom.Sqrt3 * (float64(v.X) + float64(v.Y)/2)
		y := 1.5 * float64(v.Y)

		assert.EqualDelta(t, geom.NormalizeAngle(math.Atan2(y, x)), geom.NormalizeAngle(float64(i)*geom.Pi/3), geom.Delta, direction.String())
	}
}

func TestDirection_Opposite(t *testing.T) {
	assert.Equal(t, SMinus.Opposite(), SPlus)
	assert.Equal(t, QPlus.Opposite(), QMinus)
	assert.Equal(t, RMinus.Opposite(), RPlus)
	assert.Equal(t, SPlus.Opposite(), SMinus)
	assert.Equal(t, QMinus.Opposite(), QPlus)
	assert.Equal(t, RPlus.Opposite(), RMinus)

	// double opposite returns the original direction
	for _, d := range []Direction{SMinus, QPlus, RMinus, SPlus, QMinus, RPlus} {
		assert.Equal(t, d.Opposite().Opposite(), d)
	}

	// verify neighbor offset vectors are truly opposite
	for _, d := range []Direction{SMinus, QPlus, RMinus} {
		v := d.Offset()
		opp := d.Opposite().Offset()
		assert.Equal(t, v.X+opp.X, 0)
		assert.Equal(t, v.Y+opp.Y, 0)
	}
}

func TestDirection_Rotate(t *testing.T) {
	// Rotate(0) is identity
	assert.Equal(t, SMinus.Rotate(0), SMinus)
	assert.Equal(t, RPlus.Rotate(0), RPlus)

	// Rotate(1) advances one step of increasing angle (increments index)
	assert.Equal(t, SMinus.Rotate(1), RPlus) // 0 → 1
	assert.Equal(t, RPlus.Rotate(1), QMinus) // 1 → 2
	assert.Equal(t, QPlus.Rotate(1), SMinus) // 5 → 0 (wrap)

	// Rotate(-1) advances one step of decreasing angle (decrements index)
	assert.Equal(t, SMinus.Rotate(-1), QPlus) // 0 → 5 (wrap)
	assert.Equal(t, RPlus.Rotate(-1), SMinus) // 1 → 0
	assert.Equal(t, QMinus.Rotate(-1), RPlus) // 2 → 1

	// Rotate(3) == Opposite()
	for _, d := range []Direction{SMinus, QPlus, RMinus, SPlus, QMinus, RPlus} {
		assert.Equal(t, d.Rotate(3), d.Opposite())
	}

	// Rotate(6) is identity
	for _, d := range []Direction{SMinus, QPlus, RMinus, SPlus, QMinus, RPlus} {
		assert.Equal(t, d.Rotate(6), d)
	}

	// Rotate and its inverse cancel out
	for _, d := range []Direction{SMinus, QPlus, RMinus, SPlus, QMinus, RPlus} {
		assert.Equal(t, d.Rotate(2).Rotate(-2), d)
	}
}

func Test_String(t *testing.T) {
	assert.Equal(t, SMinus.String(), "SMinus")
	assert.Equal(t, QPlus.String(), "QPlus")
	assert.Equal(t, RMinus.String(), "RMinus")
	assert.Equal(t, SPlus.String(), "SPlus")
	assert.Equal(t, QMinus.String(), "QMinus")
	assert.Equal(t, RPlus.String(), "RPlus")
	// out-of-range values wrap into [SMinus, QPlus], negatives included
	assert.Equal(t, Direction(6).String(), "SMinus")
	assert.Equal(t, Direction(8).String(), "QMinus")
	assert.Equal(t, Direction(-1).String(), "QPlus")
	assert.Equal(t, Direction(-8).String(), "RMinus")
}
