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
	sinkPoint      geom.Point[int]
	sinkVector     geom.Vector[int]
	sinkVectors    [6]geom.Vector[int]
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

func TestRangeLen(t *testing.T) {
	t.Run("the center and its rings", func(t *testing.T) {
		assert.Equal(t, RangeLen(0), 1)
		assert.Equal(t, RangeLen(1), 7)
		assert.Equal(t, RangeLen(2), 19)
		assert.Equal(t, RangeLen(10), 331)
	})
	t.Run("negative radius holds nothing", func(t *testing.T) {
		assert.Equal(t, RangeLen(-1), 0)
		assert.Equal(t, RangeLen(-5), 0)
	})
	t.Run("the length of Range and Spiral", func(t *testing.T) {
		for n := -2; n <= 12; n++ {
			assert.Equal(t, RangeLen(n), len(testHex.Range(n)), strconv.Itoa(n)+": ")
			assert.Equal(t, RangeLen(n), len(testHex.Spiral(n)), strconv.Itoa(n)+": ")
		}
	})
	t.Run("grows by a ring each step", func(t *testing.T) {
		for n := 1; n <= 12; n++ {
			assert.Equal(t, RangeLen(n)-RangeLen(n-1), len(testHex.Ring(n)), strconv.Itoa(n)+": ")
		}
	})
}

func TestRingLen(t *testing.T) {
	t.Run("six hexes a step out", func(t *testing.T) {
		assert.Equal(t, RingLen(1), 6)
		assert.Equal(t, RingLen(2), 12)
		assert.Equal(t, RingLen(10), 60)
	})
	t.Run("zero radius is the center alone", func(t *testing.T) {
		assert.Equal(t, RingLen(0), 1)
	})
	t.Run("negative radius holds nothing", func(t *testing.T) {
		assert.Equal(t, RingLen(-1), 0)
		assert.Equal(t, RingLen(-5), 0)
	})
	t.Run("the length of Ring", func(t *testing.T) {
		for radius := -2; radius <= 12; radius++ {
			assert.Equal(t, RingLen(radius), len(testHex.Ring(radius)), strconv.Itoa(radius)+": ")
		}
	})
	t.Run("the rings add up to the range", func(t *testing.T) {
		total := 0
		for radius := 0; radius <= 12; radius++ {
			total += RingLen(radius)
			assert.Equal(t, total, RangeLen(radius), strconv.Itoa(radius)+": ")
		}
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
			assert.Equal(t, len(hexes), RangeLen(n))
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
		buffer := make([]Hex, 0, RangeLen(3))
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkHexes = testHex.AppendRange(buffer, 3)
		}), 0.0)
	})
}

func BenchmarkHex_AppendRange(b *testing.B) {
	buffer := make([]Hex, 0, RangeLen(10))

	for b.Loop() {
		sinkHexes = testHex.AppendRange(buffer, 10)
	}
}

func TestHex_RangeSeq(t *testing.T) {
	t.Run("yields the hexes of Range in order", func(t *testing.T) {
		for _, n := range []int{0, 1, 2, 5} {
			assert.Equal(t, slices.Collect(testHex.RangeSeq(n)), testHex.Range(n), strconv.Itoa(n)+": ")
		}
	})
	t.Run("negative radius yields nothing", func(t *testing.T) {
		for range testHex.RangeSeq(-1) {
			assert.Fail(t, "yielded a hex")
		}
	})
	t.Run("stops when the loop breaks", func(t *testing.T) {
		for _, n := range []int{0, 1, 3} {
			for stop := range len(testHex.Range(n)) {
				walked := 0
				for range testHex.RangeSeq(n) {
					if walked == stop {
						break
					}
					walked++
				}
				assert.Equal(t, walked, stop, strconv.Itoa(n)+": ")
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			for h := range testHex.RangeSeq(3) {
				sinkHex = h
			}
		}), 0.0)
	})
}

func BenchmarkHex_RangeSeq(b *testing.B) {
	for b.Loop() {
		for h := range testHex.RangeSeq(10) {
			sinkHex = h
		}
	}
}

func TestHex_RangeIntersection(t *testing.T) {
	other := Pt(2, 2)

	t.Run("the hexes two ranges share", func(t *testing.T) {
		assert.Equal(t, testHexZero.RangeIntersection(1, Pt(2, 0), 1), []Hex{Pt(1, 0)})
		assert.Equal(t, testHexZero.RangeIntersection(2, Pt(2, 0), 1), []Hex{Pt(1, 0), Pt(1, 1), Pt(2, -1), Pt(2, 0)})
	})
	t.Run("holds exactly the hexes within both radii", func(t *testing.T) {
		for _, n := range []int{0, 1, 3} {
			for _, m := range []int{0, 2, 4} {
				shared := testHex.RangeIntersection(n, other, m)
				for _, h := range testHex.Range(8) {
					assert.Equal(t, slices.Contains(shared, h), testHex.DistanceTo(h) <= n && other.DistanceTo(h) <= m, h.String()+": ")
				}
			}
		}
	})
	t.Run("ordered like Compare", func(t *testing.T) {
		assert.True(t, slices.IsSortedFunc(testHex.RangeIntersection(3, other, 4), Hex.Compare))
	})
	t.Run("the same from either hex", func(t *testing.T) {
		assert.Equal(t, testHex.RangeIntersection(3, other, 4), other.RangeIntersection(4, testHex, 3))
	})
	t.Run("with itself is the smaller range", func(t *testing.T) {
		assert.Equal(t, testHex.RangeIntersection(3, testHex, 2), testHex.Range(2))
	})
	t.Run("ranges apart share nothing", func(t *testing.T) {
		assert.Equal(t, testHexZero.RangeIntersection(1, Pt(5, 0), 1), nil)
	})
	t.Run("negative radius is nil", func(t *testing.T) {
		assert.Equal(t, testHexZero.RangeIntersection(-1, testHexZero, 2), nil)
		assert.Equal(t, testHexZero.RangeIntersection(2, testHexZero, -1), nil)
	})
	t.Run("fills its capacity", func(t *testing.T) {
		shared := testHex.RangeIntersection(3, other, 4)
		assert.Equal(t, cap(shared), len(shared))
	})
	t.Run("allocates once", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkHexes = testHex.RangeIntersection(3, other, 4)
		}), 1.0)
	})
}

func BenchmarkHex_RangeIntersection(b *testing.B) {
	other := Pt(4, -2)

	for b.Loop() {
		sinkHexes = testHex.RangeIntersection(10, other, 10)
	}
}

func FuzzHex_RangeIntersection(f *testing.F) {
	f.Add(0, 0, 1, 2, 0, 1)
	f.Add(-1, 3, 3, 2, 2, 4)
	f.Add(5, -9, 0, 0, 0, 6)
	f.Add(-1000000, 999999, 7, -3, 4, -1)

	f.Fuzz(func(t *testing.T, q, r, n, stepQ, stepR, m int) {
		if max(geom.Abs(q), geom.Abs(r)) > 1e6 || max(geom.Abs(n), geom.Abs(m), geom.Abs(stepQ), geom.Abs(stepR)) > 12 {
			t.Skip()
		}

		source := Pt(q, r)
		other := source.Add(Pt(stepQ, stepR))
		message := source.String() + " within " + strconv.Itoa(n) + " and " + other.String() + " within " + strconv.Itoa(m) + ": "

		var expected []Hex
		for _, h := range source.Range(n) {
			if other.DistanceTo(h) <= m {
				expected = append(expected, h)
			}
		}

		assert.Equal(t, source.RangeIntersection(n, other, m), expected, message)
	})
}

func TestHex_AppendRangeIntersection(t *testing.T) {
	other := Pt(2, 2)

	t.Run("appends after dst", func(t *testing.T) {
		assert.Equal(t, testHex.AppendRangeIntersection([]Hex{testHexZero}, 3, other, 4), append([]Hex{testHexZero}, testHex.RangeIntersection(3, other, 4)...))
	})
	t.Run("ranges apart append nothing", func(t *testing.T) {
		assert.Equal(t, testHexZero.AppendRangeIntersection([]Hex{testHex}, 1, Pt(5, 0), 1), []Hex{testHex})
	})
	t.Run("allocates nothing with room", func(t *testing.T) {
		buffer := make([]Hex, 0, RangeLen(3))
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkHexes = testHex.AppendRangeIntersection(buffer, 3, other, 4)
		}), 0.0)
	})
}

func TestHex_RangeIntersectionSeq(t *testing.T) {
	other := Pt(2, 2)

	t.Run("yields the hexes of RangeIntersection in order", func(t *testing.T) {
		assert.Equal(t, slices.Collect(testHex.RangeIntersectionSeq(3, other, 4)), testHex.RangeIntersection(3, other, 4))
	})
	t.Run("ranges apart yield nothing", func(t *testing.T) {
		for range testHexZero.RangeIntersectionSeq(1, Pt(5, 0), 1) {
			assert.Fail(t, "yielded a hex")
		}
	})
	t.Run("stops when the loop breaks", func(t *testing.T) {
		for stop := range len(testHex.RangeIntersection(3, other, 4)) {
			walked := 0
			for range testHex.RangeIntersectionSeq(3, other, 4) {
				if walked == stop {
					break
				}
				walked++
			}
			assert.Equal(t, walked, stop)
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			for h := range testHex.RangeIntersectionSeq(3, other, 4) {
				sinkHex = h
			}
		}), 0.0)
	})
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
			assert.Equal(t, len(ring), RingLen(radius))
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
	buffer := make([]Hex, 0, RingLen(10))

	for b.Loop() {
		sinkHexes = testHex.AppendRing(buffer, 10)
	}
}

func TestHex_RingSeq(t *testing.T) {
	t.Run("yields the hexes of Ring in order", func(t *testing.T) {
		for _, radius := range []int{0, 1, 2, 5} {
			assert.Equal(t, slices.Collect(testHex.RingSeq(radius)), testHex.Ring(radius), strconv.Itoa(radius)+": ")
		}
	})
	t.Run("negative radius yields nothing", func(t *testing.T) {
		for range testHex.RingSeq(-1) {
			assert.Fail(t, "yielded a hex")
		}
	})
	t.Run("stops when the loop breaks", func(t *testing.T) {
		for _, radius := range []int{0, 1, 3} {
			for stop := range len(testHex.Ring(radius)) {
				walked := 0
				for range testHex.RingSeq(radius) {
					if walked == stop {
						break
					}
					walked++
				}
				assert.Equal(t, walked, stop, strconv.Itoa(radius)+": ")
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			for h := range testHex.RingSeq(3) {
				sinkHex = h
			}
		}), 0.0)
	})
}

func BenchmarkHex_RingSeq(b *testing.B) {
	for b.Loop() {
		for h := range testHex.RingSeq(10) {
			sinkHex = h
		}
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
		buffer := make([]Hex, 0, RangeLen(3))
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkHexes = testHex.AppendSpiral(buffer, 3)
		}), 0.0)
	})
}

func BenchmarkHex_AppendSpiral(b *testing.B) {
	buffer := make([]Hex, 0, RangeLen(10))

	for b.Loop() {
		sinkHexes = testHex.AppendSpiral(buffer, 10)
	}
}

func TestHex_SpiralSeq(t *testing.T) {
	t.Run("yields the hexes of Spiral in order", func(t *testing.T) {
		for _, radius := range []int{0, 1, 2, 5} {
			assert.Equal(t, slices.Collect(testHex.SpiralSeq(radius)), testHex.Spiral(radius), strconv.Itoa(radius)+": ")
		}
	})
	t.Run("negative radius yields nothing", func(t *testing.T) {
		for range testHex.SpiralSeq(-1) {
			assert.Fail(t, "yielded a hex")
		}
	})
	t.Run("stops when the loop breaks", func(t *testing.T) {
		for _, radius := range []int{0, 1, 3} {
			for stop := range len(testHex.Spiral(radius)) {
				walked := 0
				for range testHex.SpiralSeq(radius) {
					if walked == stop {
						break
					}
					walked++
				}
				assert.Equal(t, walked, stop, strconv.Itoa(radius)+": ")
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			for h := range testHex.SpiralSeq(3) {
				sinkHex = h
			}
		}), 0.0)
	})
}

func BenchmarkHex_SpiralSeq(b *testing.B) {
	for b.Loop() {
		for h := range testHex.SpiralSeq(10) {
			sinkHex = h
		}
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

func TestHex_LineSeq(t *testing.T) {
	t.Run("yields the hexes of Line in order", func(t *testing.T) {
		for _, target := range testHex.Spiral(4) {
			assert.Equal(t, slices.Collect(testHex.LineSeq(target)), testHex.Line(target), target.String()+": ")
		}
	})
	t.Run("to itself yields the one hex", func(t *testing.T) {
		assert.Equal(t, slices.Collect(testHex.LineSeq(testHex)), []Hex{testHex})
	})
	t.Run("stops when the loop breaks", func(t *testing.T) {
		target := Pt(4, -1)
		for stop := range len(testHexZero.Line(target)) {
			walked := 0
			for range testHexZero.LineSeq(target) {
				if walked == stop {
					break
				}
				walked++
			}
			assert.Equal(t, walked, stop)
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			for h := range testHexZero.LineSeq(Pt(4, -1)) {
				sinkHex = h
			}
		}), 0.0)
	})
}

func BenchmarkHex_LineSeq(b *testing.B) {
	target := Pt(10, -4)

	for b.Loop() {
		for h := range testHexZero.LineSeq(target) {
			sinkHex = h
		}
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
	t.Run("blocked exactly by the hexes the segment runs through", func(t *testing.T) {
		for _, target := range testHexZero.Spiral(4) {
			for _, h := range testHexZero.Range(5) {
				if h == testHexZero || h == target {
					continue
				}
				assert.Equal(t, testHexZero.HasLineOfSight(target, []Hex{h}), !runsThrough(testHexZero, target, h), target.String()+" past "+h.String()+": ")
			}
		}
	})
	t.Run("a hex clipped at a corner blocks though it is off the Line", func(t *testing.T) {
		assert.False(t, slices.Contains(testHexZero.Line(Pt(3, 1)), Pt(2, 0)))
		assert.False(t, testHexZero.HasLineOfSight(Pt(3, 1), []Hex{Pt(2, 0)}))
	})
	t.Run("a hex met at one corner does not block", func(t *testing.T) {
		assert.True(t, testHexZero.HasLineOfSight(Pt(4, 1), []Hex{Pt(1, 1)}))
		assert.False(t, testHexZero.HasLineOfSight(Pt(4, 1), []Hex{Pt(1, 0)}))
		assert.False(t, testHexZero.HasLineOfSight(Pt(4, 1), []Hex{Pt(2, 0)}))
	})
	t.Run("along an edge only both hexes beside it block", func(t *testing.T) {
		for _, direction := range Directions() {
			target := testHex.DiagonalNeighbor(direction)
			before, after := testHex.Neighbor(direction), testHex.Neighbor(direction.Turn(1))

			assert.True(t, testHex.HasLineOfSight(target, []Hex{before}), direction.String()+": ")
			assert.True(t, testHex.HasLineOfSight(target, []Hex{after}), direction.String()+": ")
			assert.False(t, testHex.HasLineOfSight(target, []Hex{before, after}), direction.String()+": ")
			assert.False(t, testHex.HasLineOfSight(target, []Hex{after, before}), direction.String()+": ")
		}
	})
	t.Run("two hexes beside different edges leave the gap open", func(t *testing.T) {
		for _, direction := range Directions() {
			middle := testHex.DiagonalNeighbor(direction)
			target := middle.DiagonalNeighbor(direction)
			near, far := testHex.Neighbor(direction), middle.Neighbor(direction.Turn(1))

			assert.True(t, testHex.HasLineOfSight(target, []Hex{near, far}), direction.String()+": ")
			assert.False(t, testHex.HasLineOfSight(target, []Hex{middle}), direction.String()+": ")
			assert.False(t, testHex.HasLineOfSight(target, []Hex{far, middle.Neighbor(direction)}), direction.String()+": ")
		}
	})
	t.Run("another blocker never opens the view", func(t *testing.T) {
		target := Pt(4, -1)
		for _, first := range testHexZero.Range(4) {
			for _, second := range testHexZero.Range(4) {
				if !testHexZero.HasLineOfSight(target, []Hex{first}) {
					assert.False(t, testHexZero.HasLineOfSight(target, []Hex{first, second}), first.String()+" and "+second.String()+": ")
				}
			}
		}
	})
	t.Run("exact far from the origin", func(t *testing.T) {
		far := Pt(1<<28, -1<<27)
		for _, target := range testHexZero.Spiral(3) {
			for _, h := range testHexZero.Range(4) {
				assert.Equal(t, far.Add(testHexZero).HasLineOfSight(far.Add(target), []Hex{far.Add(h)}), testHexZero.HasLineOfSight(target, []Hex{h}), target.String()+" past "+h.String()+": ")
			}
		}
	})
	t.Run("exact across a long segment", func(t *testing.T) {
		target := Pt(1<<14, 0)

		assert.False(t, testHexZero.HasLineOfSight(target, []Hex{Pt(1<<13, 0)}))
		assert.True(t, testHexZero.HasLineOfSight(target, []Hex{Pt(1<<13, 1)}))
	})
	t.Run("sees the same both ways", func(t *testing.T) {
		for _, source := range testHexZero.Spiral(2) {
			for _, target := range testHexZero.Spiral(2) {
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

func FuzzHex_HasLineOfSight(f *testing.F) {
	f.Add(0, 0, 4, -1, 2, -1)
	f.Add(-1, 3, 3, 1, 2, 0)
	f.Add(7, -12, 2, 2, 1, 0)
	f.Add(-1000000, 999999, -40, 17, -20, 9)

	f.Fuzz(func(t *testing.T, q, r, stepQ, stepR, blockQ, blockR int) {
		if max(geom.Abs(q), geom.Abs(r)) > 1e6 || max(geom.Abs(stepQ), geom.Abs(stepR), geom.Abs(blockQ), geom.Abs(blockR)) > 60 {
			t.Skip()
		}

		source := Pt(q, r)
		step, block := Pt(stepQ, stepR), Pt(blockQ, blockR)
		target, blocker := source.Add(step), source.Add(block)
		sees := source.HasLineOfSight(target, []Hex{blocker})
		message := source.String() + " to " + target.String() + " past " + blocker.String() + ": "

		assert.Equal(t, target.HasLineOfSight(source, []Hex{blocker}), sees, message)
		assert.Equal(t, source.HasLineOfSightFunc(target, among([]Hex{blocker})), sees, message)
		assert.Equal(t, testHexZero.HasLineOfSight(step, []Hex{block}), sees, message)
		assert.Equal(t, testHexZero.HasLineOfSight(step.ReflectQ(), []Hex{block.ReflectQ()}), sees, message)
		for steps := 1; steps < 6; steps++ {
			assert.Equal(t, testHexZero.HasLineOfSight(step.Turn(steps), []Hex{block.Turn(steps)}), sees, message)
		}

		if !sees {
			assert.True(t, source.DistanceTo(blocker)+blocker.DistanceTo(target) <= source.DistanceTo(target)+1, message)
		}
	})
}

func TestHex_HasLineOfSightFunc(t *testing.T) {
	t.Run("a nil blocked blocks nothing", func(t *testing.T) {
		assert.True(t, testHexZero.HasLineOfSightFunc(Pt(3, 0), nil))
	})
	t.Run("clear to itself and to a neighbor whatever blocks", func(t *testing.T) {
		everything := func(Hex) bool {
			return true
		}

		assert.True(t, testHex.HasLineOfSightFunc(testHex, everything))
		for _, neighbor := range testHex.Neighbors() {
			assert.True(t, testHex.HasLineOfSightFunc(neighbor, everything), neighbor.String()+": ")
		}
	})
	t.Run("blocked between", func(t *testing.T) {
		assert.False(t, testHexZero.HasLineOfSightFunc(Pt(3, 0), among([]Hex{Pt(1, 0)})))
		assert.False(t, testHexZero.HasLineOfSightFunc(Pt(3, 0), among([]Hex{Pt(2, 0)})))
	})
	t.Run("the endpoints and a blocker beyond the target are ignored", func(t *testing.T) {
		assert.True(t, testHexZero.HasLineOfSightFunc(Pt(2, 0), among([]Hex{testHexZero, Pt(2, 0), Pt(3, 0)})))
	})
	t.Run("along an edge only both hexes beside it block", func(t *testing.T) {
		for _, direction := range Directions() {
			target := testHex.DiagonalNeighbor(direction)
			before, after := testHex.Neighbor(direction), testHex.Neighbor(direction.Turn(1))
			assert.True(t, testHex.HasLineOfSightFunc(target, among([]Hex{before})), direction.String()+": ")
			assert.True(t, testHex.HasLineOfSightFunc(target, among([]Hex{after})), direction.String()+": ")
			assert.False(t, testHex.HasLineOfSightFunc(target, among([]Hex{before, after})), direction.String()+": ")
		}
	})
	t.Run("agrees with HasLineOfSight", func(t *testing.T) {
		for _, stride := range []int{2, 3, 5, 7, 11} {
			blocking := everyNth(testHexZero.Range(7), stride)
			blocked := among(blocking)
			for _, source := range testHexZero.Spiral(1) {
				for _, target := range testHexZero.Spiral(6) {
					assert.Equal(t, source.HasLineOfSightFunc(target, blocked), source.HasLineOfSight(target, blocking), source.String()+" to "+target.String()+" at stride "+strconv.Itoa(stride)+": ")
				}
			}
		}
	})
	t.Run("asks only about hexes between the two, a few each step", func(t *testing.T) {
		for _, target := range testHexZero.Spiral(6) {
			asked := 0
			testHexZero.HasLineOfSightFunc(target, func(h Hex) bool {
				asked++
				assert.True(t, h != testHexZero && h != target, target.String()+" asked about "+h.String()+": ")
				assert.True(t, testHexZero.DistanceTo(h)+h.DistanceTo(target) <= testHexZero.DistanceTo(target)+1, target.String()+" asked about "+h.String()+": ")

				return false
			})
			assert.True(t, asked <= 2*testHexZero.DistanceTo(target), target.String()+": ")
		}
	})
	t.Run("stops at the first blocker", func(t *testing.T) {
		asked := 0
		assert.False(t, testHexZero.HasLineOfSightFunc(Pt(6, 0), func(h Hex) bool {
			asked++

			return true
		}))
		assert.Equal(t, asked, 1)
	})
	t.Run("exact far from the origin", func(t *testing.T) {
		far := Pt(1<<28, -1<<27)
		for _, target := range testHexZero.Spiral(3) {
			for _, h := range testHexZero.Range(4) {
				assert.Equal(t, far.HasLineOfSightFunc(far.Add(target), among([]Hex{far.Add(h)})), testHexZero.HasLineOfSight(target, []Hex{h}), target.String()+" past "+h.String()+": ")
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		blocked := among([]Hex{Pt(9, 9)})
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkBool = testHexZero.HasLineOfSightFunc(Pt(4, -1), blocked)
		}), 0.0)
	})
}

func BenchmarkHex_HasLineOfSightFunc(b *testing.B) {
	blocked := among(testHexZero.Ring(5))

	for b.Loop() {
		sinkBool = testHexZero.HasLineOfSightFunc(Pt(10, -4), blocked)
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

func TestHex_FieldOfViewSeq(t *testing.T) {
	candidates := testHexZero.Range(3)
	blocking := []Hex{Pt(1, 0), Pt(-1, 2)}

	t.Run("yields the hexes of FieldOfView in order", func(t *testing.T) {
		assert.Equal(t, slices.Collect(testHexZero.FieldOfViewSeq(candidates, blocking)), testHexZero.FieldOfView(candidates, blocking))
		assert.Equal(t, slices.Collect(testHexZero.FieldOfViewSeq(candidates, nil)), candidates)
	})
	t.Run("no candidates yield nothing", func(t *testing.T) {
		for range testHexZero.FieldOfViewSeq(nil, blocking) {
			assert.Fail(t, "yielded a hex")
		}
	})
	t.Run("stops when the loop breaks", func(t *testing.T) {
		for stop := range len(testHexZero.FieldOfView(candidates, blocking)) {
			walked := 0
			for range testHexZero.FieldOfViewSeq(candidates, blocking) {
				if walked == stop {
					break
				}
				walked++
			}
			assert.Equal(t, walked, stop)
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			for h := range testHexZero.FieldOfViewSeq(candidates, blocking) {
				sinkHex = h
			}
		}), 0.0)
	})
}

func TestHex_FieldOfViewFunc(t *testing.T) {
	candidates := testHexZero.Range(5)

	t.Run("a nil blocked shows every candidate", func(t *testing.T) {
		assert.Equal(t, testHexZero.FieldOfViewFunc(candidates, nil), candidates)
	})
	t.Run("no candidates see nothing", func(t *testing.T) {
		assert.Equal(t, len(testHexZero.FieldOfViewFunc(nil, among([]Hex{Pt(1, 0)}))), 0)
	})
	t.Run("agrees with FieldOfView", func(t *testing.T) {
		for _, stride := range []int{2, 3, 5, 7, 11} {
			blocking := everyNth(candidates, stride)
			assert.Equal(t, testHex.FieldOfViewFunc(candidates, among(blocking)), testHex.FieldOfView(candidates, blocking), "stride "+strconv.Itoa(stride)+": ")
		}
	})
	t.Run("allocates once at the number of candidates", func(t *testing.T) {
		blocked := among([]Hex{Pt(1, 0), Pt(-1, 2)})
		assert.Equal(t, cap(testHexZero.FieldOfViewFunc(candidates, blocked)), len(candidates))
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkHexes = testHexZero.FieldOfViewFunc(candidates, blocked)
		}), 1.0)
	})
}

func BenchmarkHex_FieldOfViewFunc(b *testing.B) {
	candidates := testHexZero.Range(10)
	blocked := among(testHexZero.Ring(5)[:10])

	for b.Loop() {
		sinkHexes = testHexZero.FieldOfViewFunc(candidates, blocked)
	}
}

func TestHex_AppendFieldOfViewFunc(t *testing.T) {
	candidates := testHexZero.Range(3)
	blocked := among([]Hex{Pt(1, 0), Pt(-1, 2)})

	t.Run("appends the hexes of FieldOfViewFunc after dst", func(t *testing.T) {
		assert.Equal(t, testHexZero.AppendFieldOfViewFunc([]Hex{testHex}, candidates, blocked), append([]Hex{testHex}, testHexZero.FieldOfViewFunc(candidates, blocked)...))
	})
	t.Run("appends nothing when all are blocked", func(t *testing.T) {
		assert.Equal(t, testHexZero.AppendFieldOfViewFunc([]Hex{testHex}, []Hex{Pt(2, 0), Pt(3, 0)}, among([]Hex{Pt(1, 0)})), []Hex{testHex})
	})
	t.Run("allocates nothing with room", func(t *testing.T) {
		buffer := make([]Hex, 0, len(candidates))
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkHexes = testHexZero.AppendFieldOfViewFunc(buffer, candidates, blocked)
		}), 0.0)
	})
}

func TestHex_FieldOfViewFuncSeq(t *testing.T) {
	candidates := testHexZero.Range(3)
	blocked := among([]Hex{Pt(1, 0), Pt(-1, 2)})

	t.Run("yields the hexes of FieldOfViewFunc in order", func(t *testing.T) {
		assert.Equal(t, slices.Collect(testHexZero.FieldOfViewFuncSeq(candidates, blocked)), testHexZero.FieldOfViewFunc(candidates, blocked))
	})
	t.Run("stops when the loop breaks", func(t *testing.T) {
		for stop := range len(testHexZero.FieldOfViewFunc(candidates, blocked)) {
			walked := 0
			for range testHexZero.FieldOfViewFuncSeq(candidates, blocked) {
				if walked == stop {
					break
				}
				walked++
			}
			assert.Equal(t, walked, stop)
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			for h := range testHexZero.FieldOfViewFuncSeq(candidates, blocked) {
				sinkHex = h
			}
		}), 0.0)
	})
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

// runsThrough reports whether the segment between the centers of source and target passes
// through the inside of hex. It clips the segment against the three bands of cube space the hex
// is the meet of, a formulation that shares nothing with the corner test it checks.
func runsThrough(source, target, hex Hex) bool {
	start, step := source.Subtract(hex), target.Subtract(source)
	startQ, startR, startS := start.QRS()
	stepQ, stepR, stepS := step.QRS()

	enter, leave := 0.0, 1.0
	for _, band := range [3][2]int{{startQ - startR, stepQ - stepR}, {startR - startS, stepR - stepS}, {startS - startQ, stepS - stepQ}} {
		offset, slope := float64(band[0]), float64(band[1])
		if slope == 0 {
			if math.Abs(offset) >= 1 {
				return false
			}
			continue
		}

		enter = max(enter, min((-1-offset)/slope, (1-offset)/slope))
		leave = min(leave, max((-1-offset)/slope, (1-offset)/slope))
	}

	return leave-enter > geom.Delta
}

// among returns the predicate holding for the given hexes and no other.
func among(hexes []Hex) func(Hex) bool {
	return func(hex Hex) bool {
		return slices.Contains(hexes, hex)
	}
}

// everyNth returns the hexes at every stride-th place of the given ones.
func everyNth(hexes []Hex, stride int) []Hex {
	var result []Hex
	for i := stride - 1; i < len(hexes); i += stride {
		result = append(result, hexes[i])
	}

	return result
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

func ExampleRangeLen() {
	buffer := make([]Hex, 0, RangeLen(2))
	buffer = Pt(0, 0).AppendRange(buffer, 2)

	fmt.Println(RangeLen(2), len(buffer), cap(buffer))
	// Output: 19 19 19
}

func ExampleRingLen() {
	fmt.Println(RingLen(0), RingLen(1), RingLen(2))
	// Output: 1 6 12
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

func ExampleHex_RangeIntersection() {
	fmt.Println(Pt(0, 0).RangeIntersection(2, Pt(2, 0), 1))
	// Output: [(1,0) (1,1) (2,-1) (2,0)]
}

func ExampleHex_Ring() {
	fmt.Println(Pt(0, 0).Ring(1))
	// Output: [(1,0) (0,1) (-1,1) (-1,0) (0,-1) (1,-1)]
}

func ExampleHex_Spiral() {
	fmt.Println(Pt(0, 0).Spiral(1))
	// Output: [(0,0) (1,0) (0,1) (-1,1) (-1,0) (0,-1) (1,-1)]
}

func ExampleHex_SpiralSeq() {
	occupied := []Hex{Pt(2, 0), Pt(0, -1)}

	for h := range Pt(0, 0).SpiralSeq(3) {
		if slices.Contains(occupied, h) {
			fmt.Println("nearest occupied hex:", h)
			break
		}
	}
	// Output: nearest occupied hex: (0,-1)
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

func ExampleHex_LineSeq() {
	walls := []Hex{Pt(3, -1)}

	for h := range Pt(0, 0).LineSeq(Pt(4, -2)) {
		if slices.Contains(walls, h) {
			break
		}
		fmt.Println(h)
	}
	// Output:
	// (0,0)
	// (1,0)
	// (2,-1)
}

func ExampleHex_HasLineOfSight() {
	source, target := Pt(0, 0), Pt(3, 0)

	fmt.Println(source.HasLineOfSight(target, []Hex{Pt(2, 0)}))
	fmt.Println(source.HasLineOfSight(target, []Hex{Pt(0, 2)}))
	// Output:
	// false
	// true
}

func ExampleHex_HasLineOfSightFunc() {
	walls := map[Hex]bool{Pt(2, 0): true}
	blocked := func(h Hex) bool {
		return walls[h]
	}

	fmt.Println(Pt(0, 0).HasLineOfSightFunc(Pt(3, 0), blocked))
	fmt.Println(Pt(0, 0).HasLineOfSightFunc(Pt(0, 3), blocked))
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

func ExampleHex_FieldOfViewSeq() {
	walls := []Hex{Pt(1, 0)}

	for h := range Pt(0, 0).FieldOfViewSeq(Pt(0, 0).Ring(2), walls) {
		fmt.Println("first visible hex of the ring:", h)
		break
	}
	// Output: first visible hex of the ring: (1,1)
}

func ExampleHex_FieldOfViewFunc() {
	walls := map[Hex]bool{Pt(1, 0): true}
	candidates := []Hex{Pt(2, 0), Pt(0, 2), Pt(1, 0)}

	fmt.Println(Pt(0, 0).FieldOfViewFunc(candidates, func(h Hex) bool {
		return walls[h]
	}))
	// Output: [(0,2) (1,0)]
}

func ExampleHex_Compare() {
	hexes := []Hex{Pt(1, 2), Pt(-1, 0), Pt(1, -3)}
	slices.SortFunc(hexes, Hex.Compare)

	fmt.Println(hexes)
	// Output: [(-1,0) (1,-3) (1,2)]
}
