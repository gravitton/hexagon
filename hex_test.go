package hex_test

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"
	"strconv"
	"testing"

	"github.com/gravitton/assert"
	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/ints"
	. "github.com/gravitton/hexagon"
	"github.com/gravitton/hexagon/hextest"
)

var (
	testHex     = Pt(-1, 3)
	testHexZero = Pt(0, 0)
)

var (
	sinkHex        Hex
	sinkHexes      []Hex
	sinkBool       bool
	sinkInt        int
	sinkFracHex    FractionalHex
	sinkDirection  Direction
	sinkDirections [6]Direction
	sinkSystems    [7]CoordinateSystem
	sinkPoint      ints.Point
	sinkVector     ints.Vector
	sinkVectors    [6]ints.Vector
)

func TestHex_Constructor(t *testing.T) {
	hextest.AssertHex(t, Pt(-1, 3), Hex{Q: -1, R: 3})
}

func TestParseHex(t *testing.T) {
	t.Run("the form String prints", func(t *testing.T) {
		h, err := ParseHex("(-1,3)")
		assert.NoError(t, err)
		hextest.AssertHex(t, h, testHex)
	})
	t.Run("rejects fractional values", func(t *testing.T) {
		_, err := ParseHex("(1.5,2)")
		assert.ErrorContains(t, err, "hex: invalid q value")
	})
	t.Run("the parse error is wrapped", func(t *testing.T) {
		_, err := ParseHex("(99999999999999999999,1)")
		assert.ErrorContains(t, err, "hex: invalid q value")
		assert.ErrorIs(t, err, strconv.ErrRange)

		_, err = ParseHex("(1,b)")
		assert.ErrorContains(t, err, "hex: invalid r value")
		assert.ErrorIs(t, err, strconv.ErrSyntax)
	})
	t.Run("malformed input", func(t *testing.T) {
		for _, s := range []string{"", "()", "(1)", "(1,2,3)", "1,2", "(1,2", "1,2)", "((1,2))", "( 1,2)", "(1, 2)"} {
			_, err := ParseHex(s)
			assert.Error(t, err, s+": ")
		}
	})
	t.Run("names the type of a malformed value", func(t *testing.T) {
		_, err := ParseHex("1,2")
		assert.ErrorContains(t, err, `hex: invalid hex format "1,2"`)
	})
	t.Run("round-trips with String", func(t *testing.T) {
		for _, h := range testHex.Spiral(2) {
			parsed, err := ParseHex(h.String())
			assert.NoError(t, err)
			hextest.AssertHex(t, parsed, h, h.String()+": ")
		}
	})
}

func BenchmarkParseHex(b *testing.B) {
	for b.Loop() {
		sinkHex, _ = ParseHex("(-12,345)")
	}
}

func FuzzParseHex(f *testing.F) {
	for _, seed := range []string{"(-1,3)", "(0,0)", "(+1,-0)", "(1.5,2)", "(1,2,3)", "((1,2))", "( 1,2)", "", "(99999999999999999999,1)"} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, s string) {
		message := strconv.Quote(s) + ": "

		h, err := ParseHex(s)
		if err != nil {
			hextest.AssertHex(t, h, testHexZero, message)
			return
		}

		parsed, err := ParseHex(h.String())
		assert.NoError(t, err, message)
		hextest.AssertHex(t, parsed, h, message)
	})
}

func TestHex_S(t *testing.T) {
	assert.Equal(t, testHex.S(), -2)
}

func TestHex_QR(t *testing.T) {
	q, r := testHex.QR()
	assert.Equal(t, q, -1)
	assert.Equal(t, r, 3)
}

func TestHex_QRS(t *testing.T) {
	q, r, s := testHex.QRS()
	assert.Equal(t, q, -1)
	assert.Equal(t, r, 3)
	assert.Equal(t, s, -2)
}

func TestHex_Length(t *testing.T) {
	t.Run("steps from the origin", func(t *testing.T) {
		assert.Equal(t, testHex.Length(), 3)
		assert.Equal(t, testHexZero.Length(), 0)
	})
	t.Run("counts the rings of a spiral", func(t *testing.T) {
		for _, h := range testHexZero.Spiral(3) {
			assert.True(t, h.Length() <= 3, h.String()+": ")
		}
	})
}

func TestHex_Add(t *testing.T) {
	t.Run("vector sum", func(t *testing.T) {
		hextest.AssertHex(t, testHex.Add(Pt(3, -2)), Pt(2, 1))
	})
	t.Run("commutative", func(t *testing.T) {
		for _, h := range testHexZero.Spiral(2) {
			hextest.AssertHex(t, testHex.Add(h), h.Add(testHex), h.String()+": ")
		}
	})
	t.Run("the origin changes nothing", func(t *testing.T) {
		hextest.AssertHex(t, testHex.Add(testHexZero), testHex)
	})
}

func TestHex_Subtract(t *testing.T) {
	t.Run("vector difference", func(t *testing.T) {
		hextest.AssertHex(t, testHex.Subtract(Pt(3, -2)), Pt(-4, 5))
	})
	t.Run("undoes Add", func(t *testing.T) {
		for _, h := range testHexZero.Spiral(2) {
			hextest.AssertHex(t, testHex.Add(h).Subtract(h), testHex, h.String()+": ")
		}
	})
	t.Run("of itself is the origin", func(t *testing.T) {
		hextest.AssertHex(t, testHex.Subtract(testHex), testHexZero)
	})
}

func TestHex_Multiply(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		hextest.AssertHex(t, testHex.Multiply(3), Pt(-3, 9))
	})
	t.Run("zero", func(t *testing.T) {
		hextest.AssertHex(t, testHex.Multiply(0), testHexZero)
	})
	t.Run("negative", func(t *testing.T) {
		hextest.AssertHex(t, testHex.Multiply(-1), Pt(1, -3))
		hextest.AssertHex(t, testHex.Multiply(-2), Pt(2, -6))
	})
	t.Run("scales the length", func(t *testing.T) {
		for _, h := range testHexZero.Spiral(2) {
			for factor := -3; factor <= 3; factor++ {
				assert.Equal(t, h.Multiply(factor).Length(), geom.Abs(factor)*h.Length(), h.String()+" by "+strconv.Itoa(factor)+": ")
			}
		}
	})
}

func TestHex_Lerp(t *testing.T) {
	a := Pt(0, 0)
	b := Pt(4, -2)

	t.Run("ends and midpoint", func(t *testing.T) {
		hextest.AssertHex(t, a.Lerp(b, 0), a)
		hextest.AssertHex(t, a.Lerp(b, 1), b)
		hextest.AssertHex(t, a.Lerp(b, 0.5), Pt(2, -1))
	})
	t.Run("stays on the straight line between the two hexes", func(t *testing.T) {
		for i := range 11 {
			h := a.Lerp(b, float64(i)/10)
			assert.Equal(t, a.DistanceTo(h)+h.DistanceTo(b), a.DistanceTo(b))
		}
	})
	t.Run("panics for a non-finite t", func(t *testing.T) {
		assert.Panics(t, func() {
			sinkHex = a.Lerp(b, math.NaN())
		})
		assert.Panics(t, func() {
			sinkHex = a.Lerp(b, math.Inf(1))
		})
	})
	t.Run("allocates nothing", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkHex = a.Lerp(b, 0.3)
		}), 0.0)
	})
}

func BenchmarkHex_Lerp(b *testing.B) {
	target := Pt(10, -4)

	for b.Loop() {
		sinkHex = testHexZero.Lerp(target, 0.3)
	}
}

func TestHex_Turn(t *testing.T) {
	h := Pt(3, 0)

	t.Run("zero and six steps are the identity", func(t *testing.T) {
		hextest.AssertHex(t, h.Turn(0), h)
		hextest.AssertHex(t, h.Turn(6), h)
	})
	t.Run("one step forward", func(t *testing.T) {
		hextest.AssertHex(t, h.Turn(1), Pt(-h.R, h.Q+h.R))
	})
	t.Run("one step back", func(t *testing.T) {
		hextest.AssertHex(t, h.Turn(-1), Pt(h.Q+h.R, -h.Q))
	})
	t.Run("there and back returns to the start", func(t *testing.T) {
		hextest.AssertHex(t, h.Turn(1).Turn(-1), h)
	})
	t.Run("agrees with turning the direction", func(t *testing.T) {
		for _, direction := range Directions() {
			for steps := -6; steps <= 6; steps++ {
				hextest.AssertHex(t, direction.Hex().Turn(steps), direction.Turn(steps).Hex(), direction.String()+": ")
			}
		}
	})
	t.Run("keeps the length", func(t *testing.T) {
		for _, h := range testHexZero.Spiral(3) {
			for steps := -6; steps <= 6; steps++ {
				assert.Equal(t, h.Turn(steps).Length(), h.Length(), h.String()+" by "+strconv.Itoa(steps)+": ")
			}
		}
	})
	t.Run("steps add up", func(t *testing.T) {
		for first := -6; first <= 6; first++ {
			for second := -6; second <= 6; second++ {
				hextest.AssertHex(t, testHex.Turn(first).Turn(second), testHex.Turn(first+second), strconv.Itoa(first)+" then "+strconv.Itoa(second)+": ")
			}
		}
	})
	t.Run("any step count, without overflow", func(t *testing.T) {
		hextest.AssertHex(t, testHex.Turn(math.MaxInt), testHex.Turn(math.MaxInt%6))
		hextest.AssertHex(t, testHex.Turn(math.MinInt), testHex.Turn(math.MinInt%6))
	})
}

func BenchmarkHex_Turn(b *testing.B) {
	for b.Loop() {
		sinkHex = testHex.Turn(5)
	}
}

func TestHex_TurnAround(t *testing.T) {
	center := Pt(1, 1)
	h := Pt(3, 0)

	t.Run("zero and six steps are the identity", func(t *testing.T) {
		hextest.AssertHex(t, h.TurnAround(center, 0), h)
		hextest.AssertHex(t, h.TurnAround(center, 6), h)
	})
	t.Run("keeps the distance from the center", func(t *testing.T) {
		for steps := 1; steps <= 5; steps++ {
			assert.Equal(t, center.DistanceTo(h.TurnAround(center, steps)), center.DistanceTo(h))
		}
	})
	t.Run("there and back returns to the start", func(t *testing.T) {
		hextest.AssertHex(t, h.TurnAround(center, 1).TurnAround(center, -1), h)
	})
	t.Run("around the origin is Turn", func(t *testing.T) {
		hextest.AssertHex(t, h.TurnAround(testHexZero, 1), h.Turn(1))
		hextest.AssertHex(t, h.TurnAround(testHexZero, -1), h.Turn(-1))
	})
	t.Run("the center stays put", func(t *testing.T) {
		for steps := -6; steps <= 6; steps++ {
			hextest.AssertHex(t, center.TurnAround(center, steps), center, strconv.Itoa(steps)+": ")
		}
	})
	t.Run("a step carries the ring one side along", func(t *testing.T) {
		for _, radius := range []int{1, 2, 3} {
			ring := center.Ring(radius)
			for i, on := range ring {
				hextest.AssertHex(t, on.TurnAround(center, 1), ring[(i+radius)%len(ring)], on.String()+": ")
			}
		}
	})
}

func TestHex_ReflectQ(t *testing.T) {
	t.Run("swaps r and s", func(t *testing.T) {
		hextest.AssertHex(t, Pt(0, 0).ReflectQ(), Pt(0, 0))
		hextest.AssertHex(t, Pt(1, 0).ReflectQ(), Pt(1, -1))
		hextest.AssertHex(t, Pt(2, 0).ReflectQ(), Pt(2, -2))
		hextest.AssertHex(t, Pt(2, -2).ReflectQ(), Pt(2, 0))
	})
	t.Run("keeps q", func(t *testing.T) {
		assert.Equal(t, testHex.ReflectQ().Q, testHex.Q)
	})
	t.Run("twice is the identity", func(t *testing.T) {
		hextest.AssertHex(t, testHex.ReflectQ().ReflectQ(), testHex)
	})
	t.Run("keeps the length", func(t *testing.T) {
		for _, h := range testHexZero.Spiral(3) {
			assert.Equal(t, h.ReflectQ().Length(), h.Length(), h.String()+": ")
		}
	})
	t.Run("followed by ReflectR is a third of a turn", func(t *testing.T) {
		for _, h := range testHexZero.Spiral(3) {
			hextest.AssertHex(t, h.ReflectQ().ReflectR(), h.Turn(4), h.String()+": ")
		}
	})
}

func TestHex_ReflectR(t *testing.T) {
	t.Run("swaps q and s", func(t *testing.T) {
		hextest.AssertHex(t, Pt(0, 0).ReflectR(), Pt(0, 0))
		hextest.AssertHex(t, Pt(1, 0).ReflectR(), Pt(-1, 0))
		hextest.AssertHex(t, Pt(2, 0).ReflectR(), Pt(-2, 0))
		hextest.AssertHex(t, Pt(-2, 0).ReflectR(), Pt(2, 0))
	})
	t.Run("keeps r", func(t *testing.T) {
		assert.Equal(t, testHex.ReflectR().R, testHex.R)
	})
	t.Run("twice is the identity", func(t *testing.T) {
		hextest.AssertHex(t, testHex.ReflectR().ReflectR(), testHex)
	})
	t.Run("keeps the length", func(t *testing.T) {
		for _, h := range testHexZero.Spiral(3) {
			assert.Equal(t, h.ReflectR().Length(), h.Length(), h.String()+": ")
		}
	})
	t.Run("followed by ReflectS is a third of a turn", func(t *testing.T) {
		for _, h := range testHexZero.Spiral(3) {
			hextest.AssertHex(t, h.ReflectR().ReflectS(), h.Turn(4), h.String()+": ")
		}
	})
}

func TestHex_ReflectS(t *testing.T) {
	t.Run("swaps q and r", func(t *testing.T) {
		hextest.AssertHex(t, Pt(0, 0).ReflectS(), Pt(0, 0))
		hextest.AssertHex(t, Pt(1, 0).ReflectS(), Pt(0, 1))
		hextest.AssertHex(t, Pt(3, -1).ReflectS(), Pt(-1, 3))
		hextest.AssertHex(t, Pt(-1, 3).ReflectS(), Pt(3, -1))
	})
	t.Run("keeps s", func(t *testing.T) {
		assert.Equal(t, testHex.ReflectS().S(), testHex.S())
	})
	t.Run("twice is the identity", func(t *testing.T) {
		hextest.AssertHex(t, testHex.ReflectS().ReflectS(), testHex)
	})
	t.Run("keeps the length", func(t *testing.T) {
		for _, h := range testHexZero.Spiral(3) {
			assert.Equal(t, h.ReflectS().Length(), h.Length(), h.String()+": ")
		}
	})
	t.Run("followed by ReflectQ is a third of a turn", func(t *testing.T) {
		for _, h := range testHexZero.Spiral(3) {
			hextest.AssertHex(t, h.ReflectS().ReflectQ(), h.Turn(4), h.String()+": ")
		}
	})
}

func TestHex_Neighbor(t *testing.T) {
	t.Run("one step in each direction", func(t *testing.T) {
		hextest.AssertHex(t, testHex.Neighbor(SMinus), Pt(0, 3))
		hextest.AssertHex(t, testHex.Neighbor(RPlus), Pt(-1, 4))
		hextest.AssertHex(t, testHex.Neighbor(QMinus), Pt(-2, 4))
		hextest.AssertHex(t, testHex.Neighbor(SPlus), Pt(-2, 3))
		hextest.AssertHex(t, testHex.Neighbor(RMinus), Pt(-1, 2))
		hextest.AssertHex(t, testHex.Neighbor(QPlus), Pt(0, 2))
	})
	t.Run("out-of-range directions wrap", func(t *testing.T) {
		hextest.AssertHex(t, testHex.Neighbor(Direction(6)), Pt(0, 3))
		hextest.AssertHex(t, testHex.Neighbor(Direction(-2)), Pt(-1, 2))
	})
	t.Run("none steps nowhere", func(t *testing.T) {
		hextest.AssertHex(t, testHex.Neighbor(DirectionNone), testHex)
	})
	t.Run("the opposite direction steps back", func(t *testing.T) {
		for _, direction := range Directions() {
			hextest.AssertHex(t, testHex.Neighbor(direction).Neighbor(direction.Opposite()), testHex, direction.String()+": ")
		}
	})
}

func TestHex_Neighbors(t *testing.T) {
	t.Run("by increasing angle", func(t *testing.T) {
		assert.Equal(t, testHex.Neighbors(), [6]Hex{Pt(0, 3), Pt(-1, 4), Pt(-2, 4), Pt(-2, 3), Pt(-1, 2), Pt(0, 2)})
	})
	t.Run("every neighbor is one step away", func(t *testing.T) {
		for _, neighbor := range testHex.Neighbors() {
			assert.Equal(t, testHex.DistanceTo(neighbor), 1)
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkHex = testHex.Neighbors()[0]
		}), 0.0)
	})
}

func BenchmarkHex_Neighbors(b *testing.B) {
	for b.Loop() {
		sinkHex = testHex.Neighbors()[0]
	}
}

func TestHex_DiagonalNeighbor(t *testing.T) {
	t.Run("two steps between two neighbors", func(t *testing.T) {
		hextest.AssertHex(t, testHexZero.DiagonalNeighbor(SMinus), Pt(1, 1))
		hextest.AssertHex(t, testHexZero.DiagonalNeighbor(RPlus), Pt(-1, 2))
		hextest.AssertHex(t, testHexZero.DiagonalNeighbor(QMinus), Pt(-2, 1))
		hextest.AssertHex(t, testHexZero.DiagonalNeighbor(SPlus), Pt(-1, -1))
		hextest.AssertHex(t, testHexZero.DiagonalNeighbor(RMinus), Pt(1, -2))
		hextest.AssertHex(t, testHexZero.DiagonalNeighbor(QPlus), Pt(2, -1))
	})
	t.Run("beside the neighbor of the direction and of the one after", func(t *testing.T) {
		for _, direction := range Directions() {
			diagonal := testHex.DiagonalNeighbor(direction)
			assert.Equal(t, testHex.DistanceTo(diagonal), 2, direction.String()+": ")
			assert.Equal(t, testHex.Neighbor(direction).DistanceTo(diagonal), 1, direction.String()+": ")
			assert.Equal(t, testHex.Neighbor(direction.Turn(1)).DistanceTo(diagonal), 1, direction.String()+": ")
		}
	})
	t.Run("a turn carries it to the next diagonal", func(t *testing.T) {
		for _, direction := range Directions() {
			hextest.AssertHex(t, testHex.DiagonalNeighbor(direction).TurnAround(testHex, 1), testHex.DiagonalNeighbor(direction.Turn(1)), direction.String()+": ")
		}
	})
	t.Run("the opposite direction steps back", func(t *testing.T) {
		for _, direction := range Directions() {
			hextest.AssertHex(t, testHex.DiagonalNeighbor(direction).DiagonalNeighbor(direction.Opposite()), testHex, direction.String()+": ")
		}
	})
	t.Run("out-of-range directions wrap", func(t *testing.T) {
		hextest.AssertHex(t, testHexZero.DiagonalNeighbor(Direction(6)), Pt(1, 1))
		hextest.AssertHex(t, testHexZero.DiagonalNeighbor(Direction(-2)), Pt(1, -2))
	})
	t.Run("none steps nowhere", func(t *testing.T) {
		hextest.AssertHex(t, testHex.DiagonalNeighbor(DirectionNone), testHex)
	})
}

func TestHex_DiagonalNeighbors(t *testing.T) {
	t.Run("by increasing angle", func(t *testing.T) {
		assert.Equal(t, testHexZero.DiagonalNeighbors(), [6]Hex{Pt(1, 1), Pt(-1, 2), Pt(-2, 1), Pt(-1, -1), Pt(1, -2), Pt(2, -1)})
	})
	t.Run("reads DiagonalNeighbor by direction", func(t *testing.T) {
		diagonals := testHex.DiagonalNeighbors()
		for i, direction := range Directions() {
			hextest.AssertHex(t, diagonals[i], testHex.DiagonalNeighbor(direction), direction.String()+": ")
		}
	})
	t.Run("every other hex of the second ring, from the second", func(t *testing.T) {
		ring := testHex.Ring(2)
		for i, diagonal := range testHex.DiagonalNeighbors() {
			hextest.AssertHex(t, diagonal, ring[2*i+1])
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkHex = testHex.DiagonalNeighbors()[0]
		}), 0.0)
	})
}

func BenchmarkHex_DiagonalNeighbors(b *testing.B) {
	for b.Loop() {
		sinkHex = testHex.DiagonalNeighbors()[0]
	}
}

func TestHex_Range(t *testing.T) {
	t.Run("negative radius is nil", func(t *testing.T) {
		assert.Equal(t, testHexZero.Range(-1), nil)
	})
	t.Run("zero radius is the center", func(t *testing.T) {
		assert.Equal(t, testHexZero.Range(0), []Hex{Pt(0, 0)})
	})
	t.Run("radius one", func(t *testing.T) {
		assert.Equal(t, testHexZero.Range(1), []Hex{Pt(-1, 0), Pt(-1, 1), Pt(0, -1), Pt(0, 0), Pt(0, 1), Pt(1, -1), Pt(1, 0)})
	})
	t.Run("counts and fills its capacity", func(t *testing.T) {
		for _, n := range []int{1, 2, 3, 5, 10, 20} {
			hexes := testHex.Range(n)
			assert.Equal(t, len(hexes), 1+3*n*(n+1))
			assert.Equal(t, cap(hexes), len(hexes))
		}
	})
	t.Run("every hex is within the radius", func(t *testing.T) {
		for _, h := range testHex.Range(3) {
			assert.True(t, testHex.DistanceTo(h) <= 3)
		}
	})
	t.Run("holds exactly the hexes within the radius", func(t *testing.T) {
		hexes := testHex.Range(3)
		for _, h := range testHex.Spiral(5) {
			assert.Equal(t, slices.Contains(hexes, h), testHex.DistanceTo(h) <= 3, h.String()+": ")
		}
	})
	t.Run("ordered like Compare", func(t *testing.T) {
		assert.True(t, slices.IsSortedFunc(testHex.Range(3), Hex.Compare))
	})
	t.Run("allocates once", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkHexes = testHex.Range(3)
		}), 1.0)
	})
}

func BenchmarkHex_Range(b *testing.B) {
	for b.Loop() {
		sinkHexes = testHex.Range(10)
	}
}

func TestHex_AppendRange(t *testing.T) {
	t.Run("appends after dst", func(t *testing.T) {
		assert.Equal(t, testHex.AppendRange([]Hex{testHexZero}, 1), append([]Hex{testHexZero}, testHex.Range(1)...))
	})
	t.Run("negative radius appends nothing", func(t *testing.T) {
		assert.Equal(t, testHex.AppendRange([]Hex{testHexZero}, -1), []Hex{testHexZero})
	})
	t.Run("allocates nothing with room", func(t *testing.T) {
		buffer := make([]Hex, 0, 37)
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkHexes = testHex.AppendRange(buffer, 3)
		}), 0.0)
	})
}

func BenchmarkHex_AppendRange(b *testing.B) {
	buffer := make([]Hex, 0, 331)

	for b.Loop() {
		sinkHexes = testHex.AppendRange(buffer, 10)
	}
}

func TestHex_Ring(t *testing.T) {
	t.Run("negative radius is nil", func(t *testing.T) {
		assert.Equal(t, testHexZero.Ring(-1), nil)
	})
	t.Run("zero radius is the center", func(t *testing.T) {
		assert.Equal(t, testHexZero.Ring(0), []Hex{Pt(0, 0)})
	})
	t.Run("radius one from the SMinus corner", func(t *testing.T) {
		assert.Equal(t, testHexZero.Ring(1), []Hex{Pt(1, 0), Pt(0, 1), Pt(-1, 1), Pt(-1, 0), Pt(0, -1), Pt(1, -1)})
	})
	t.Run("radius one is the neighbors in order", func(t *testing.T) {
		neighbors := testHex.Neighbors()
		assert.Equal(t, testHex.Ring(1), neighbors[:])
	})
	t.Run("counts and fills its capacity", func(t *testing.T) {
		for _, radius := range []int{1, 2, 3, 5, 10} {
			ring := testHex.Ring(radius)
			assert.Equal(t, len(ring), 6*radius)
			assert.Equal(t, cap(ring), len(ring))
		}
	})
	t.Run("holds exactly the hexes at the radius", func(t *testing.T) {
		ring := testHex.Ring(3)
		for _, h := range testHex.Range(5) {
			assert.Equal(t, slices.Contains(ring, h), testHex.DistanceTo(h) == 3, h.String()+": ")
		}
	})
	t.Run("steps between neighbors and closes on itself", func(t *testing.T) {
		ring := testHex.Ring(3)
		for i, h := range ring {
			assert.Equal(t, h.DistanceTo(ring[(i+1)%len(ring)]), 1, h.String()+": ")
		}
	})
	t.Run("allocates once", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkHexes = testHex.Ring(3)
		}), 1.0)
	})
}

func BenchmarkHex_Ring(b *testing.B) {
	for b.Loop() {
		sinkHexes = testHex.Ring(10)
	}
}

func TestHex_AppendRing(t *testing.T) {
	t.Run("appends after dst", func(t *testing.T) {
		assert.Equal(t, testHex.AppendRing([]Hex{testHexZero}, 2), append([]Hex{testHexZero}, testHex.Ring(2)...))
	})
	t.Run("zero radius appends the center", func(t *testing.T) {
		assert.Equal(t, testHex.AppendRing([]Hex{testHexZero}, 0), []Hex{testHexZero, testHex})
	})
	t.Run("negative radius appends nothing", func(t *testing.T) {
		assert.Equal(t, testHex.AppendRing([]Hex{testHexZero}, -1), []Hex{testHexZero})
	})
	t.Run("allocates nothing with room", func(t *testing.T) {
		buffer := make([]Hex, 0, 18)
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkHexes = testHex.AppendRing(buffer, 3)
		}), 0.0)
	})
}

func BenchmarkHex_AppendRing(b *testing.B) {
	buffer := make([]Hex, 0, 60)

	for b.Loop() {
		sinkHexes = testHex.AppendRing(buffer, 10)
	}
}

func TestHex_Spiral(t *testing.T) {
	t.Run("negative radius is nil", func(t *testing.T) {
		assert.Equal(t, testHexZero.Spiral(-1), nil)
	})
	t.Run("zero radius is the center", func(t *testing.T) {
		assert.Equal(t, testHexZero.Spiral(0), []Hex{Pt(0, 0)})
	})
	t.Run("radius one is the center and its ring", func(t *testing.T) {
		assert.Equal(t, testHexZero.Spiral(1), []Hex{Pt(0, 0), Pt(1, 0), Pt(0, 1), Pt(-1, 1), Pt(-1, 0), Pt(0, -1), Pt(1, -1)})
	})
	t.Run("holds the hexes of Range", func(t *testing.T) {
		for _, n := range []int{1, 2, 3, 5} {
			spiral := testHex.Spiral(n)
			assert.Equal(t, cap(spiral), len(spiral))

			slices.SortFunc(spiral, Hex.Compare)
			assert.Equal(t, spiral, testHex.Range(n))
		}
	})
	t.Run("nearest first", func(t *testing.T) {
		spiral := testHex.Spiral(3)
		for i := 1; i < len(spiral); i++ {
			assert.True(t, testHex.DistanceTo(spiral[i-1]) <= testHex.DistanceTo(spiral[i]))
		}
	})
	t.Run("allocates once", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkHexes = testHex.Spiral(3)
		}), 1.0)
	})
}

func BenchmarkHex_Spiral(b *testing.B) {
	for b.Loop() {
		sinkHexes = testHex.Spiral(10)
	}
}

func TestHex_AppendSpiral(t *testing.T) {
	t.Run("appends after dst", func(t *testing.T) {
		assert.Equal(t, testHex.AppendSpiral([]Hex{testHexZero}, 2), append([]Hex{testHexZero}, testHex.Spiral(2)...))
	})
	t.Run("negative radius appends nothing", func(t *testing.T) {
		assert.Equal(t, testHex.AppendSpiral([]Hex{testHexZero}, -1), []Hex{testHexZero})
	})
	t.Run("allocates nothing with room", func(t *testing.T) {
		buffer := make([]Hex, 0, 37)
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkHexes = testHex.AppendSpiral(buffer, 3)
		}), 0.0)
	})
}

func BenchmarkHex_AppendSpiral(b *testing.B) {
	buffer := make([]Hex, 0, 331)

	for b.Loop() {
		sinkHexes = testHex.AppendSpiral(buffer, 10)
	}
}

func TestHex_DistanceTo(t *testing.T) {
	t.Run("steps between hexes", func(t *testing.T) {
		assert.Equal(t, testHex.DistanceTo(Pt(0, 3)), 1)
		assert.Equal(t, testHex.DistanceTo(Pt(1, 6)), 5)
		assert.Equal(t, testHex.DistanceTo(testHex), 0)
	})
	t.Run("symmetric", func(t *testing.T) {
		for _, h := range testHexZero.Spiral(3) {
			assert.Equal(t, testHex.DistanceTo(h), h.DistanceTo(testHex), h.String()+": ")
		}
	})
	t.Run("unchanged by moving both hexes", func(t *testing.T) {
		for _, h := range testHexZero.Spiral(3) {
			assert.Equal(t, testHex.Add(h).DistanceTo(h), testHex.Length(), h.String()+": ")
		}
	})
	t.Run("never longer through a third hex", func(t *testing.T) {
		target := Pt(2, -3)
		for _, via := range testHexZero.Spiral(4) {
			assert.True(t, testHex.DistanceTo(target) <= testHex.DistanceTo(via)+via.DistanceTo(target), via.String()+": ")
		}
	})
}

func BenchmarkHex_DistanceTo(b *testing.B) {
	target := Pt(10, -4)

	for b.Loop() {
		sinkInt = testHex.DistanceTo(target)
	}
}

func TestHex_DirectionTo(t *testing.T) {
	t.Run("the direction of a neighbor", func(t *testing.T) {
		for _, direction := range Directions() {
			assert.Equal(t, testHex.DirectionTo(testHex.Neighbor(direction)), direction, direction.String()+": ")
		}
	})
	t.Run("the nearest direction of a distant hex", func(t *testing.T) {
		assert.Equal(t, testHexZero.DirectionTo(Pt(5, -1)), SMinus)
		assert.Equal(t, testHexZero.DirectionTo(Pt(5, 1)), SMinus)
		assert.Equal(t, testHexZero.DirectionTo(Pt(1, 5)), RPlus)
		assert.Equal(t, testHexZero.DirectionTo(Pt(-4, -1)), SPlus)
		assert.Equal(t, testHexZero.DirectionTo(Pt(3, -7)), RMinus)
	})
	t.Run("a diagonal takes the direction it is named by", func(t *testing.T) {
		for _, direction := range Directions() {
			assert.Equal(t, testHex.DirectionTo(testHex.DiagonalNeighbor(direction)), direction, direction.String()+": ")
		}
	})
	t.Run("the hex itself has no direction", func(t *testing.T) {
		assert.Equal(t, testHex.DirectionTo(testHex), DirectionNone)
	})
	t.Run("the same along the whole ray", func(t *testing.T) {
		for _, h := range testHexZero.Spiral(4) {
			for _, factor := range []int{2, 7, 1000} {
				assert.Equal(t, testHexZero.DirectionTo(h.Multiply(factor)), testHexZero.DirectionTo(h), h.String()+" by "+strconv.Itoa(factor)+": ")
			}
		}
	})
	t.Run("turns with the target", func(t *testing.T) {
		for _, h := range testHexZero.Spiral(4)[1:] {
			for steps := -6; steps <= 6; steps++ {
				assert.Equal(t, testHexZero.DirectionTo(h.Turn(steps)), testHexZero.DirectionTo(h).Turn(steps), h.String()+" by "+strconv.Itoa(steps)+": ")
			}
		}
	})
	t.Run("within a twelfth of a turn of the angle to the target", func(t *testing.T) {
		for _, h := range testHexZero.Spiral(4)[1:] {
			x := geom.Sqrt3 * (float64(h.Q) + float64(h.R)/2)
			y := 1.5 * float64(h.R)
			gap := geom.NormalizeAngle(math.Atan2(y, x)-testHexZero.DirectionTo(h).Angle()+geom.Pi) - geom.Pi

			assert.True(t, math.Abs(gap) <= geom.Pi/6+geom.Delta, h.String()+": ")
		}
	})
	t.Run("the first step of the line is that way or beside it", func(t *testing.T) {
		for _, h := range testHex.Spiral(4)[1:] {
			first := testHex.DirectionTo(testHex.Line(h)[1])
			towards := testHex.DirectionTo(h)

			assert.True(t, first == towards || first == towards.Turn(1) || first == towards.Turn(-1), h.String()+": ")
		}
	})
}

func BenchmarkHex_DirectionTo(b *testing.B) {
	target := Pt(10, -4)

	for b.Loop() {
		sinkDirection = testHex.DirectionTo(target)
	}
}

func TestHex_Line(t *testing.T) {
	t.Run("straight lines from the origin", func(t *testing.T) {
		assert.Equal(t, testHexZero.Line(Pt(3, 0)), []Hex{Pt(0, 0), Pt(1, 0), Pt(2, 0), Pt(3, 0)})
		assert.Equal(t, testHexZero.Line(Pt(2, -1)), []Hex{Pt(0, 0), Pt(1, 0), Pt(2, -1)})
		assert.Equal(t, testHexZero.Line(Pt(-2, 1)), []Hex{Pt(0, 0), Pt(-1, 1), Pt(-2, 1)})
		assert.Equal(t, testHexZero.Line(Pt(4, -2)), []Hex{Pt(0, 0), Pt(1, 0), Pt(2, -1), Pt(3, -1), Pt(4, -2)})
		assert.Equal(t, testHexZero.Line(Pt(-4, 2)), []Hex{Pt(0, 0), Pt(-1, 1), Pt(-2, 1), Pt(-3, 2), Pt(-4, 2)})
	})
	t.Run("to itself is one hex", func(t *testing.T) {
		assert.Equal(t, testHexZero.Line(testHexZero), []Hex{Pt(0, 0)})
		assert.Equal(t, testHex.Line(testHex), []Hex{testHex})
	})
	t.Run("steps between neighbors from end to end", func(t *testing.T) {
		for _, target := range testHex.Spiral(4) {
			line := testHex.Line(target)
			assert.Equal(t, len(line), testHex.DistanceTo(target)+1)
			assert.Equal(t, cap(line), len(line))
			assert.Equal(t, line[0], testHex)
			assert.Equal(t, line[len(line)-1], target)

			for i := 1; i < len(line); i++ {
				assert.Equal(t, line[i-1].DistanceTo(line[i]), 1)
			}
		}
	})
	t.Run("reversed is the line back", func(t *testing.T) {
		for _, source := range testHexZero.Spiral(4) {
			for _, target := range testHexZero.Spiral(4) {
				back := target.Line(source)
				slices.Reverse(back)
				assert.Equal(t, source.Line(target), back, source.String()+" to "+target.String()+": ")
			}
		}
	})
	t.Run("allocates once", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkHexes = testHexZero.Line(Pt(4, -1))
		}), 1.0)
	})
}

func BenchmarkHex_Line(b *testing.B) {
	for b.Loop() {
		sinkHexes = testHexZero.Line(Pt(10, -4))
	}
}

func FuzzHex_Line(f *testing.F) {
	f.Add(0, 0, 4, -2)
	f.Add(-1, 3, 3, -3)
	f.Add(7, -12, 0, 0)
	f.Add(-1000000, 999999, -40, 17)

	f.Fuzz(func(t *testing.T, q, r, stepQ, stepR int) {
		if max(geom.Abs(q), geom.Abs(r)) > 1e6 || max(geom.Abs(stepQ), geom.Abs(stepR)) > 60 {
			t.Skip()
		}

		source := Pt(q, r)
		target := source.Add(Pt(stepQ, stepR))
		line := source.Line(target)
		distance := source.DistanceTo(target)
		message := source.String() + " to " + target.String() + ": "

		assert.Equal(t, len(line), distance+1, message)
		assert.Equal(t, line[0], source, message)
		assert.Equal(t, line[distance], target, message)

		for i, h := range line {
			assert.Equal(t, source.DistanceTo(h), i, message)
			assert.Equal(t, h.DistanceTo(target), distance-i, message)
			assert.Equal(t, source.HasLineOfSight(target, []Hex{h}), i == 0 || i == distance, message)
		}

		back := target.Line(source)
		slices.Reverse(back)
		assert.Equal(t, line, back, message)
	})
}

func TestHex_AppendLine(t *testing.T) {
	t.Run("appends after dst", func(t *testing.T) {
		assert.Equal(t, testHex.AppendLine([]Hex{testHexZero}, Pt(2, 0)), append([]Hex{testHexZero}, testHex.Line(Pt(2, 0))...))
	})
	t.Run("allocates nothing with room", func(t *testing.T) {
		buffer := make([]Hex, 0, 5)
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkHexes = testHexZero.AppendLine(buffer, Pt(4, -1))
		}), 0.0)
	})
}

func BenchmarkHex_AppendLine(b *testing.B) {
	buffer := make([]Hex, 0, 11)

	for b.Loop() {
		sinkHexes = testHexZero.AppendLine(buffer, Pt(10, -4))
	}
}

func TestHex_HasLineOfSight(t *testing.T) {
	t.Run("clear without blockers", func(t *testing.T) {
		assert.True(t, testHexZero.HasLineOfSight(Pt(3, 0), nil))
		assert.True(t, testHexZero.HasLineOfSight(Pt(3, 0), []Hex{}))
		assert.True(t, testHexZero.HasLineOfSight(testHexZero, nil))
	})
	t.Run("blocked between", func(t *testing.T) {
		assert.False(t, testHexZero.HasLineOfSight(Pt(3, 0), []Hex{Pt(1, 0)}))
		assert.False(t, testHexZero.HasLineOfSight(Pt(3, 0), []Hex{Pt(2, 0)}))
	})
	t.Run("a blocker beyond the target is ignored", func(t *testing.T) {
		assert.True(t, testHexZero.HasLineOfSight(Pt(2, 0), []Hex{Pt(3, 0)}))
	})
	t.Run("the endpoints never block", func(t *testing.T) {
		assert.True(t, testHexZero.HasLineOfSight(Pt(3, 0), []Hex{Pt(3, 0)}))
		assert.True(t, testHexZero.HasLineOfSight(Pt(3, 0), []Hex{testHexZero}))
		assert.True(t, testHexZero.HasLineOfSight(testHexZero, []Hex{testHexZero}))
		assert.True(t, testHexZero.HasLineOfSight(Pt(1, 0), []Hex{testHexZero, Pt(1, 0)}))
	})
	t.Run("blocked exactly by the hexes of Line", func(t *testing.T) {
		target := Pt(4, -1)
		line := testHexZero.Line(target)
		for _, h := range testHexZero.Range(5) {
			between := slices.Contains(line[1:len(line)-1], h)
			assert.Equal(t, testHexZero.HasLineOfSight(target, []Hex{h}), !between, h.String()+": ")
		}
	})
	t.Run("sees the same both ways", func(t *testing.T) {
		for _, source := range testHexZero.Spiral(3) {
			for _, target := range testHexZero.Spiral(3) {
				for _, h := range testHexZero.Spiral(3) {
					blocking := []Hex{h}
					assert.Equal(t, source.HasLineOfSight(target, blocking), target.HasLineOfSight(source, blocking), source.String()+" to "+target.String()+" past "+h.String()+": ")
				}
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		blocking := []Hex{Pt(9, 9)}
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkBool = testHexZero.HasLineOfSight(Pt(4, -1), blocking)
		}), 0.0)
	})
}

func BenchmarkHex_HasLineOfSight(b *testing.B) {
	blocking := testHexZero.Ring(5)

	for b.Loop() {
		sinkBool = testHexZero.HasLineOfSight(Pt(10, -4), blocking)
	}
}

func TestHex_FieldOfView(t *testing.T) {
	candidates := testHexZero.Range(3)

	t.Run("no candidates see nothing", func(t *testing.T) {
		assert.Equal(t, len(testHexZero.FieldOfView(nil, nil)), 0)
		assert.Equal(t, len(testHexZero.FieldOfView([]Hex{}, nil)), 0)
	})
	t.Run("everything is visible without blockers", func(t *testing.T) {
		assert.Equal(t, testHexZero.FieldOfView(candidates, nil), candidates)
		assert.Equal(t, testHexZero.FieldOfView(candidates, []Hex{}), candidates)
	})
	t.Run("a closed ring leaves only the ring and the center", func(t *testing.T) {
		neighbors := testHexZero.Neighbors()
		for _, v := range testHexZero.FieldOfView(candidates, neighbors[:]) {
			assert.True(t, testHexZero.DistanceTo(v) <= 1)
		}
	})
	t.Run("a gap in the ring lets the line through", func(t *testing.T) {
		blocking := []Hex{Pt(0, -1), Pt(1, -1), Pt(-1, 0), Pt(-1, 1), Pt(0, 1)}
		assert.Equal(t, testHexZero.FieldOfView([]Hex{Pt(2, 0), Pt(3, 0), Pt(0, -2)}, blocking), []Hex{Pt(2, 0), Pt(3, 0)})
	})
	t.Run("agrees with HasLineOfSight", func(t *testing.T) {
		blocking := []Hex{Pt(1, 0), Pt(-1, 2)}
		visible := testHexZero.FieldOfView(candidates, blocking)
		for _, candidate := range candidates {
			assert.Equal(t, slices.Contains(visible, candidate), testHexZero.HasLineOfSight(candidate, blocking), candidate.String()+": ")
		}
	})
	t.Run("keeps the order of the candidates", func(t *testing.T) {
		visible := testHexZero.FieldOfView(candidates, []Hex{Pt(1, 0), Pt(-1, 2)})
		next := 0
		for _, candidate := range candidates {
			if next < len(visible) && visible[next] == candidate {
				next++
			}
		}
		assert.Equal(t, next, len(visible))
	})
	t.Run("another blocker never shows more", func(t *testing.T) {
		blocking := []Hex{Pt(1, 0), Pt(-1, 2)}
		visible := testHexZero.FieldOfView(candidates, blocking)
		for _, blocker := range testHexZero.Ring(2) {
			for _, seen := range testHexZero.FieldOfView(candidates, append(blocking, blocker)) {
				assert.True(t, slices.Contains(visible, seen), seen.String()+" past "+blocker.String()+": ")
			}
		}
	})
	t.Run("allocates once", func(t *testing.T) {
		blocking := []Hex{Pt(1, 0), Pt(-1, 2)}
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkHexes = testHexZero.FieldOfView(candidates, blocking)
		}), 1.0)
	})
}

func BenchmarkHex_FieldOfView(b *testing.B) {
	candidates := testHexZero.Range(10)
	blocking := testHexZero.Ring(5)[:10]

	for b.Loop() {
		sinkHexes = testHexZero.FieldOfView(candidates, blocking)
	}
}

func TestHex_AppendFieldOfView(t *testing.T) {
	candidates := []Hex{Pt(2, 0), Pt(3, 0)}

	t.Run("appends the visible after dst", func(t *testing.T) {
		assert.Equal(t, testHexZero.AppendFieldOfView([]Hex{testHex}, candidates, nil), []Hex{testHex, Pt(2, 0), Pt(3, 0)})
	})
	t.Run("appends nothing when all are blocked", func(t *testing.T) {
		assert.Equal(t, testHexZero.AppendFieldOfView([]Hex{testHex}, candidates, []Hex{Pt(1, 0)}), []Hex{testHex})
	})
	t.Run("allocates nothing with room", func(t *testing.T) {
		buffer := make([]Hex, 0, len(candidates))
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkHexes = testHexZero.AppendFieldOfView(buffer, candidates, nil)
		}), 0.0)
	})
}

func BenchmarkHex_AppendFieldOfView(b *testing.B) {
	candidates := testHexZero.Range(10)
	blocking := testHexZero.Ring(5)[:10]
	buffer := make([]Hex, 0, len(candidates))

	for b.Loop() {
		sinkHexes = testHexZero.AppendFieldOfView(buffer, candidates, blocking)
	}
}

func TestHex_Equal(t *testing.T) {
	t.Run("same coordinates", func(t *testing.T) {
		assert.True(t, testHex.Equal(Pt(-1, 3)))
	})
	t.Run("one coordinate differs", func(t *testing.T) {
		assert.False(t, testHex.Equal(Pt(-1, 4)))
		assert.False(t, testHex.Equal(Pt(1, 3)))
	})
}

func TestHex_Compare(t *testing.T) {
	t.Run("orders by q first", func(t *testing.T) {
		assert.Equal(t, Pt(1, 9).Compare(Pt(2, 0)), -1)
		assert.Equal(t, Pt(2, 0).Compare(Pt(1, 9)), 1)
	})
	t.Run("falls back to r", func(t *testing.T) {
		assert.Equal(t, Pt(1, 1).Compare(Pt(1, 2)), -1)
		assert.Equal(t, Pt(1, 2).Compare(Pt(1, 1)), 1)
	})
	t.Run("equal", func(t *testing.T) {
		assert.Equal(t, testHex.Compare(testHex), 0)
		assert.Equal(t, testHexZero.Compare(Pt(0, 0)), 0)
	})
	t.Run("antisymmetric and transitive", func(t *testing.T) {
		hexes := Pt(0, 0).Spiral(2)
		for _, a := range hexes {
			for _, b := range hexes {
				assert.Equal(t, a.Compare(b), -b.Compare(a))

				for _, c := range hexes {
					if a.Compare(b) < 0 && b.Compare(c) < 0 {
						assert.True(t, a.Compare(c) < 0)
					}
				}
			}
		}
	})
	t.Run("zero exactly when Equal", func(t *testing.T) {
		for _, h := range testHex.Spiral(2) {
			assert.Equal(t, testHex.Compare(h) == 0, testHex.Equal(h), h.String()+": ")
		}
	})
	t.Run("sorts", func(t *testing.T) {
		hexes := []Hex{Pt(1, 2), Pt(-1, 0), Pt(1, -3), Pt(0, 7)}
		slices.SortFunc(hexes, Hex.Compare)

		assert.Equal(t, hexes, []Hex{Pt(-1, 0), Pt(0, 7), Pt(1, -3), Pt(1, 2)})

		index, found := slices.BinarySearchFunc(hexes, Pt(1, -3), Hex.Compare)
		assert.True(t, found)
		assert.Equal(t, index, 2)
	})
}

func TestHex_IsZero(t *testing.T) {
	t.Run("origin", func(t *testing.T) {
		assert.True(t, testHexZero.IsZero())
	})
	t.Run("any other hex", func(t *testing.T) {
		assert.False(t, testHex.IsZero())
		assert.False(t, Pt(0, 1).IsZero())
		assert.False(t, Pt(1, 0).IsZero())
	})
}

func TestHex_To(t *testing.T) {
	for _, system := range CoordinateSystems() {
		t.Run(system.String(), func(t *testing.T) {
			for _, h := range testHex.Spiral(2) {
				assert.Equal(t, h.To(system), system.To(h), h.String()+": ")
			}
		})
	}
}

func TestHex_Point(t *testing.T) {
	assert.Equal(t, testHex.Point(), geom.Pt(-1, 3))
}

func TestHex_Float(t *testing.T) {
	t.Run("same coordinates", func(t *testing.T) {
		hextest.AssertFractionalHex(t, testHex.Float(), FracPt(-1, 3))
	})
	t.Run("round-trips through Round", func(t *testing.T) {
		for _, h := range testHexZero.Spiral(2) {
			hextest.AssertHex(t, h.Float().Round(), h)
		}
	})
}

func TestHex_String(t *testing.T) {
	assert.Equal(t, testHex.String(), "(-1,3)")
}

func TestHex_JSON(t *testing.T) {
	data, err := json.Marshal(testHex)
	assert.NoError(t, err)
	assert.Equal(t, string(data), `{"q":-1,"r":3}`)

	var decoded Hex
	assert.NoError(t, json.Unmarshal(data, &decoded))
	hextest.AssertHex(t, decoded, testHex)
}

func ExamplePt() {
	h := Pt(-1, 3)

	fmt.Println(h)
	fmt.Println(h.QRS())
	// Output:
	// (-1,3)
	// -1 3 -2
}

func ExampleParseHex() {
	h, err := ParseHex("(-1,3)")
	fmt.Println(h, err)

	_, err = ParseHex("(1.5,2)")
	fmt.Println(err != nil)
	// Output:
	// (-1,3) <nil>
	// true
}

func ExampleHex_Turn() {
	h := Pt(3, 0)

	fmt.Println(h.Turn(1))
	fmt.Println(h.Turn(-1))
	fmt.Println(h.Turn(3))
	// Output:
	// (0,3)
	// (3,-3)
	// (-3,0)
}

func ExampleHex_Neighbor() {
	h := Pt(-1, 3)

	fmt.Println(h.Neighbor(PointyTopEast))
	fmt.Println(h.Neighbor(FlatTopNorth))
	fmt.Println(h.Neighbor(DirectionNone))
	// Output:
	// (0,3)
	// (-1,2)
	// (-1,3)
}

func ExampleHex_DiagonalNeighbors() {
	fmt.Println(Pt(0, 0).DiagonalNeighbors())
	// Output: [(1,1) (-1,2) (-2,1) (-1,-1) (1,-2) (2,-1)]
}

func ExampleHex_Range() {
	fmt.Println(Pt(0, 0).Range(1))
	// Output: [(-1,0) (-1,1) (0,-1) (0,0) (0,1) (1,-1) (1,0)]
}

func ExampleHex_AppendRange() {
	buffer := make([]Hex, 0, 7)
	for _, center := range []Hex{Pt(0, 0), Pt(5, -2)} {
		buffer = center.AppendRange(buffer[:0], 1)
		fmt.Println(len(buffer), buffer[0])
	}
	// Output:
	// 7 (-1,0)
	// 7 (4,-2)
}

func ExampleHex_Ring() {
	fmt.Println(Pt(0, 0).Ring(1))
	// Output: [(1,0) (0,1) (-1,1) (-1,0) (0,-1) (1,-1)]
}

func ExampleHex_Spiral() {
	fmt.Println(Pt(0, 0).Spiral(1))
	// Output: [(0,0) (1,0) (0,1) (-1,1) (-1,0) (0,-1) (1,-1)]
}

func ExampleHex_DirectionTo() {
	h := Pt(0, 0)

	fmt.Println(h.DirectionTo(Pt(1, 0)))
	fmt.Println(h.DirectionTo(Pt(-4, -1)))
	fmt.Println(h.DirectionTo(h))
	// Output:
	// SMinus
	// SPlus
	// None
}

func ExampleHex_Line() {
	fmt.Println(Pt(0, 0).Line(Pt(4, -2)))
	// Output: [(0,0) (1,0) (2,-1) (3,-1) (4,-2)]
}

func ExampleHex_HasLineOfSight() {
	source, target := Pt(0, 0), Pt(3, 0)

	fmt.Println(source.HasLineOfSight(target, []Hex{Pt(2, 0)}))
	fmt.Println(source.HasLineOfSight(target, []Hex{Pt(0, 2)}))
	// Output:
	// false
	// true
}

func ExampleHex_FieldOfView() {
	walls := []Hex{Pt(1, 0)}
	candidates := []Hex{Pt(2, 0), Pt(0, 2), Pt(1, 0)}

	fmt.Println(Pt(0, 0).FieldOfView(candidates, walls))
	// Output: [(0,2) (1,0)]
}

func ExampleHex_Compare() {
	hexes := []Hex{Pt(1, 2), Pt(-1, 0), Pt(1, -3)}
	slices.SortFunc(hexes, Hex.Compare)

	fmt.Println(hexes)
	// Output: [(-1,0) (1,-3) (1,2)]
}
