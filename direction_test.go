package hex_test

import (
	"encoding/json"
	"fmt"
	"math"
	"testing"

	"github.com/gravitton/assert"
	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/ints"
	. "github.com/gravitton/hexagon"
	"github.com/gravitton/hexagon/hextest"
)

var axialDirection = [6]ints.Vector{geom.Vec(1, 0), geom.Vec(0, 1), geom.Vec(-1, 1), geom.Vec(-1, 0), geom.Vec(0, -1), geom.Vec(1, -1)}

func TestDirections(t *testing.T) {
	t.Run("by increasing angle", func(t *testing.T) {
		assert.Equal(t, Directions(), [6]Direction{SMinus, RPlus, QMinus, SPlus, RMinus, QPlus})
	})
	t.Run("returns a fresh array", func(t *testing.T) {
		list := Directions()
		list[0] = QPlus
		assert.Equal(t, Directions()[0], SMinus)
	})
	t.Run("flat-top aliases", func(t *testing.T) {
		assert.Equal(t, FlatTopSouthEast, SMinus)
		assert.Equal(t, FlatTopSouth, RPlus)
		assert.Equal(t, FlatTopSouthWest, QMinus)
		assert.Equal(t, FlatTopNorthWest, SPlus)
		assert.Equal(t, FlatTopNorth, RMinus)
		assert.Equal(t, FlatTopNorthEast, QPlus)
	})
	t.Run("pointy-top aliases", func(t *testing.T) {
		assert.Equal(t, PointyTopEast, SMinus)
		assert.Equal(t, PointyTopSouthEast, RPlus)
		assert.Equal(t, PointyTopSouthWest, QMinus)
		assert.Equal(t, PointyTopWest, SPlus)
		assert.Equal(t, PointyTopNorthWest, RMinus)
		assert.Equal(t, PointyTopNorthEast, QPlus)
	})
	t.Run("allocates nothing", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkDirections = Directions()
		}), 0.0)
	})
}

func TestDirectionFromAngle(t *testing.T) {
	t.Run("round-trips with Angle", func(t *testing.T) {
		for _, direction := range Directions() {
			assert.Equal(t, DirectionFromAngle(direction.Angle()), direction, direction.String()+": ")
			assert.Equal(t, DirectionFromAngle(direction.Angle()+2*geom.Pi), direction, direction.String()+": ")
		}
	})
	t.Run("rounds to the nearest direction", func(t *testing.T) {
		assert.Equal(t, DirectionFromAngle(geom.Pi/3-0.1), RPlus)
		assert.Equal(t, DirectionFromAngle(2*geom.Pi-0.1), SMinus)
	})
	t.Run("non-finite angles have no direction", func(t *testing.T) {
		assert.Equal(t, DirectionFromAngle(math.NaN()), DirectionNone)
		assert.Equal(t, DirectionFromAngle(math.Inf(1)), DirectionNone)
		assert.Equal(t, DirectionFromAngle(math.Inf(-1)), DirectionNone)
	})
	t.Run("holds across the twelfth of a turn either side", func(t *testing.T) {
		for _, direction := range Directions() {
			assert.Equal(t, DirectionFromAngle(direction.Angle()-geom.Pi/6+0.01), direction, direction.String()+": ")
			assert.Equal(t, DirectionFromAngle(direction.Angle()+geom.Pi/6-0.01), direction, direction.String()+": ")
		}
	})
	t.Run("any number of whole turns away", func(t *testing.T) {
		for _, direction := range Directions() {
			for turns := -3; turns <= 3; turns++ {
				assert.Equal(t, DirectionFromAngle(direction.Angle()+float64(turns)*2*geom.Pi), direction, direction.String()+": ")
			}
		}
	})
}

func BenchmarkDirectionFromAngle(b *testing.B) {
	for b.Loop() {
		sinkDirection = DirectionFromAngle(4.2)
	}
}

func FuzzDirectionFromAngle(f *testing.F) {
	f.Add(0.0)
	f.Add(geom.Pi)
	f.Add(-0.1)
	f.Add(geom.Pi / 6)
	f.Add(-1e6)
	f.Add(math.NaN())
	f.Add(math.Inf(1))

	f.Fuzz(func(t *testing.T, angle float64) {
		direction := DirectionFromAngle(angle)
		message := fmt.Sprintf("%v → %s: ", angle, direction)

		if math.IsNaN(angle) || math.IsInf(angle, 0) {
			assert.Equal(t, direction, DirectionNone, message)
			return
		}
		if math.Abs(angle) > 1e6 {
			t.Skip()
		}

		assert.True(t, direction >= SMinus && direction <= QPlus, message)

		gap := geom.NormalizeAngle(angle-direction.Angle()+geom.Pi) - geom.Pi
		assert.True(t, math.Abs(gap) <= geom.Pi/6+geom.Delta, message)
	})
}

func TestParseDirection(t *testing.T) {
	t.Run("round-trips with String", func(t *testing.T) {
		for _, direction := range Directions() {
			parsed, err := ParseDirection(direction.String())
			assert.NoError(t, err)
			assert.Equal(t, parsed, direction)
		}
	})
	t.Run("none", func(t *testing.T) {
		parsed, err := ParseDirection("None")
		assert.NoError(t, err)
		assert.Equal(t, parsed, DirectionNone)
	})
	t.Run("unknown name", func(t *testing.T) {
		parsed, err := ParseDirection("Northwest")
		assert.Error(t, err)
		assert.Equal(t, parsed, DirectionNone)
	})
}

func TestDirection_Opposite(t *testing.T) {
	t.Run("half a turn", func(t *testing.T) {
		assert.Equal(t, SMinus.Opposite(), SPlus)
		assert.Equal(t, RPlus.Opposite(), RMinus)
		assert.Equal(t, QMinus.Opposite(), QPlus)
		assert.Equal(t, SPlus.Opposite(), SMinus)
		assert.Equal(t, RMinus.Opposite(), RPlus)
		assert.Equal(t, QPlus.Opposite(), QMinus)
	})
	t.Run("twice is the identity", func(t *testing.T) {
		for _, direction := range Directions() {
			assert.Equal(t, direction.Opposite().Opposite(), direction)
		}
	})
	t.Run("the offsets cancel", func(t *testing.T) {
		for _, direction := range Directions() {
			assert.Equal(t, direction.Offset().Add(direction.Opposite().Offset()), ints.Vector{}, direction.String()+": ")
		}
	})
	t.Run("none is its own opposite", func(t *testing.T) {
		assert.Equal(t, DirectionNone.Opposite(), DirectionNone)
	})
}

func TestDirection_Turn(t *testing.T) {
	t.Run("zero and six steps are the identity", func(t *testing.T) {
		for _, direction := range Directions() {
			assert.Equal(t, direction.Turn(0), direction)
			assert.Equal(t, direction.Turn(6), direction)
		}
	})
	t.Run("one step of increasing angle", func(t *testing.T) {
		assert.Equal(t, SMinus.Turn(1), RPlus)
		assert.Equal(t, RPlus.Turn(1), QMinus)
		assert.Equal(t, QPlus.Turn(1), SMinus)
	})
	t.Run("one step of decreasing angle", func(t *testing.T) {
		assert.Equal(t, SMinus.Turn(-1), QPlus)
		assert.Equal(t, RPlus.Turn(-1), SMinus)
		assert.Equal(t, QMinus.Turn(-1), RPlus)
	})
	t.Run("three steps are Opposite", func(t *testing.T) {
		for _, direction := range Directions() {
			assert.Equal(t, direction.Turn(3), direction.Opposite())
		}
	})
	t.Run("there and back returns to the start", func(t *testing.T) {
		for _, direction := range Directions() {
			assert.Equal(t, direction.Turn(2).Turn(-2), direction)
		}
	})
	t.Run("steps add up", func(t *testing.T) {
		for _, direction := range Directions() {
			for first := -6; first <= 6; first++ {
				for second := -6; second <= 6; second++ {
					assert.Equal(t, direction.Turn(first).Turn(second), direction.Turn(first+second), direction.String()+": ")
				}
			}
		}
	})
	t.Run("a step lands on the hex beside", func(t *testing.T) {
		for _, direction := range Directions() {
			assert.Equal(t, direction.Hex().DistanceTo(direction.Turn(1).Hex()), 1, direction.String()+": ")
		}
	})
	t.Run("a step raises the angle by a sixth of a turn", func(t *testing.T) {
		for _, direction := range Directions() {
			assert.EqualDelta(t, geom.NormalizeAngle(direction.Turn(1).Angle()-direction.Angle()), geom.Pi/3, geom.Delta, direction.String()+": ")
		}
	})
	t.Run("any step count, without overflow", func(t *testing.T) {
		assert.Equal(t, QPlus.Turn(math.MaxInt), QPlus.Turn(math.MaxInt%6))
		assert.Equal(t, QPlus.Turn(math.MinInt), QPlus.Turn(math.MinInt%6))
		assert.Equal(t, Direction(math.MaxInt).Turn(1), Direction(math.MaxInt%6).Turn(1))
	})
	t.Run("none turns to itself", func(t *testing.T) {
		assert.Equal(t, DirectionNone.Turn(1), DirectionNone)
	})
}

func TestDirection_Offset(t *testing.T) {
	t.Run("the axial step of each direction", func(t *testing.T) {
		assert.Equal(t, SMinus.Offset(), axialDirection[0])
		assert.Equal(t, RPlus.Offset(), axialDirection[1])
		assert.Equal(t, QMinus.Offset(), axialDirection[2])
		assert.Equal(t, SPlus.Offset(), axialDirection[3])
		assert.Equal(t, RMinus.Offset(), axialDirection[4])
		assert.Equal(t, QPlus.Offset(), axialDirection[5])
	})
	t.Run("out-of-range values wrap", func(t *testing.T) {
		assert.Equal(t, Direction(6).Offset(), axialDirection[0])
		assert.Equal(t, Direction(8).Offset(), axialDirection[2])
		assert.Equal(t, Direction(15).Offset(), axialDirection[3])
		assert.Equal(t, Direction(-2).Offset(), axialDirection[4])
		assert.Equal(t, Direction(-8).Offset(), axialDirection[4])
	})
	t.Run("none steps nowhere", func(t *testing.T) {
		assert.Equal(t, DirectionNone.Offset(), ints.Vector{})
	})
	t.Run("the six cancel out", func(t *testing.T) {
		sum := ints.Vector{}
		for _, direction := range Directions() {
			sum = sum.Add(direction.Offset())
		}
		assert.Equal(t, sum, ints.Vector{})
	})
	t.Run("every step is one hex long and no two are the same", func(t *testing.T) {
		directions := Directions()
		for i, direction := range directions {
			assert.Equal(t, direction.Hex().Length(), 1, direction.String()+": ")
			for _, other := range directions[i+1:] {
				assert.NotEqual(t, direction.Offset(), other.Offset(), direction.String()+" and "+other.String()+": ")
			}
		}
	})
}

func TestDirection_Hex(t *testing.T) {
	t.Run("the offset as a hex", func(t *testing.T) {
		for _, direction := range Directions() {
			offset := direction.Offset()
			hextest.AssertHex(t, direction.Hex(), Pt(offset.X, offset.Y), direction.String()+": ")
		}
	})
	t.Run("none is the zero hex", func(t *testing.T) {
		hextest.AssertHex(t, DirectionNone.Hex(), testHexZero)
	})
}

func TestDirection_Angle(t *testing.T) {
	t.Run("sixths of a turn", func(t *testing.T) {
		assert.EqualDelta(t, SMinus.Angle(), 0, geom.Delta)
		assert.EqualDelta(t, RPlus.Angle(), geom.Pi/3, geom.Delta)
		assert.EqualDelta(t, SPlus.Angle(), geom.Pi, geom.Delta)
		assert.EqualDelta(t, QPlus.Angle(), 5*geom.Pi/3, geom.Delta)
	})
	t.Run("out-of-range values wrap", func(t *testing.T) {
		assert.EqualDelta(t, Direction(6).Angle(), 0, geom.Delta)
		assert.EqualDelta(t, Direction(-2).Angle(), 4*geom.Pi/3, geom.Delta)
	})
	t.Run("is the angle of the offset in a pointy-top layout", func(t *testing.T) {
		for _, direction := range Directions() {
			v := direction.Offset()
			x := geom.Sqrt3 * (float64(v.X) + float64(v.Y)/2)
			y := 1.5 * float64(v.Y)

			assert.EqualDelta(t, geom.NormalizeAngle(math.Atan2(y, x)), direction.Angle(), geom.Delta, direction.String()+": ")
		}
	})
	t.Run("a flat-top layout adds a twelfth of a turn", func(t *testing.T) {
		for _, direction := range Directions() {
			v := direction.Offset()
			x := 1.5 * float64(v.X)
			y := geom.Sqrt3 * (float64(v.Y) + float64(v.X)/2)

			assert.EqualDelta(t, geom.NormalizeAngle(math.Atan2(y, x)), direction.Angle()+geom.Pi/6, geom.Delta, direction.String()+": ")
		}
	})
	t.Run("none has no angle", func(t *testing.T) {
		assert.True(t, math.IsNaN(DirectionNone.Angle()))
	})
}

func TestDirection_IsNone(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		assert.True(t, DirectionNone.IsNone())
	})
	t.Run("the six are not none", func(t *testing.T) {
		for _, direction := range Directions() {
			assert.False(t, direction.IsNone(), direction.String()+": ")
		}
	})
	t.Run("out-of-range values wrap rather than being none", func(t *testing.T) {
		assert.False(t, Direction(6).IsNone())
	})
}

func TestDirection_String(t *testing.T) {
	t.Run("the constant names", func(t *testing.T) {
		assert.Equal(t, SMinus.String(), "SMinus")
		assert.Equal(t, RPlus.String(), "RPlus")
		assert.Equal(t, QMinus.String(), "QMinus")
		assert.Equal(t, SPlus.String(), "SPlus")
		assert.Equal(t, RMinus.String(), "RMinus")
		assert.Equal(t, QPlus.String(), "QPlus")
	})
	t.Run("out-of-range values wrap", func(t *testing.T) {
		assert.Equal(t, Direction(6).String(), "SMinus")
		assert.Equal(t, Direction(8).String(), "QMinus")
		assert.Equal(t, Direction(-2).String(), "RMinus")
		assert.Equal(t, Direction(-8).String(), "RMinus")
	})
	t.Run("none", func(t *testing.T) {
		assert.Equal(t, DirectionNone.String(), "None")
	})
}

func TestDirection_Text(t *testing.T) {
	t.Run("round-trips through JSON as the name", func(t *testing.T) {
		directions := Directions()
		for _, direction := range append(directions[:], DirectionNone) {
			data, err := json.Marshal(direction)
			assert.NoError(t, err)
			assert.Equal(t, string(data), `"`+direction.String()+`"`)

			var decoded Direction
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, direction)
		}
	})
	t.Run("unknown name fails", func(t *testing.T) {
		var decoded Direction
		assert.Error(t, json.Unmarshal([]byte(`"Northwest"`), &decoded))
	})
}

func ExampleDirectionFromAngle() {
	fmt.Println(DirectionFromAngle(0))
	fmt.Println(DirectionFromAngle(math.Pi))
	fmt.Println(DirectionFromAngle(math.NaN()))
	// Output:
	// SMinus
	// SPlus
	// None
}

func ExampleDirection_Turn() {
	fmt.Println(SMinus.Turn(1))
	fmt.Println(SMinus.Turn(-1))
	fmt.Println(SMinus.Opposite())
	// Output:
	// RPlus
	// QPlus
	// SPlus
}

func ExampleDirection_Hex() {
	for _, direction := range Directions() {
		fmt.Println(direction, direction.Hex())
	}
	// Output:
	// SMinus (1,0)
	// RPlus (0,1)
	// QMinus (-1,1)
	// SPlus (-1,0)
	// RMinus (0,-1)
	// QPlus (1,-1)
}
