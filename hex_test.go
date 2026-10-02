package hex_test

import (
	"encoding/json"
	"slices"
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
	sinkHex   Hex
	sinkHexes []Hex
	sinkBool  bool
)

func TestHex_Constructor(t *testing.T) {
	hextest.AssertHex(t, Pt(-1, 3), Hex{Q: -1, R: 3})
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
			assert.True(t, h.Length() <= 3, h.String())
		}
	})
}

func TestHex_Add(t *testing.T) {
	hextest.AssertHex(t, testHex.Add(Pt(3, -2)), Pt(2, 1))
}

func TestHex_Subtract(t *testing.T) {
	hextest.AssertHex(t, testHex.Subtract(Pt(3, -2)), Pt(-4, 5))
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
				hextest.AssertHex(t, direction.Hex().Turn(steps), direction.Turn(steps).Hex(), direction.String())
			}
		}
	})
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
	t.Run("every hex is at the radius", func(t *testing.T) {
		for _, h := range testHex.Ring(3) {
			assert.Equal(t, testHex.DistanceTo(h), 3)
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

func TestHex_DistanceTo(t *testing.T) {
	t.Run("steps between hexes", func(t *testing.T) {
		assert.Equal(t, testHex.DistanceTo(Pt(0, 3)), 1)
		assert.Equal(t, testHex.DistanceTo(Pt(1, 6)), 5)
		assert.Equal(t, testHex.DistanceTo(testHex), 0)
	})
	t.Run("symmetric", func(t *testing.T) {
		for _, h := range testHexZero.Spiral(3) {
			assert.Equal(t, testHex.DistanceTo(h), h.DistanceTo(testHex), h.String())
		}
	})
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
			assert.Equal(t, testHexZero.HasLineOfSight(target, []Hex{h}), !between, h.String())
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
			assert.Equal(t, slices.Contains(visible, candidate), testHexZero.HasLineOfSight(candidate, blocking), candidate.String())
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
				assert.Equal(t, h.To(system), system.To(h), h.String())
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
