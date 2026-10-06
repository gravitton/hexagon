package hex_test

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"testing"

	"github.com/gravitton/assert"
	geom "github.com/gravitton/geometry"
	. "github.com/gravitton/hexagon"
	"github.com/gravitton/hexagon/hextest"
)

var testFracHex = FracPt(10.9, -1.2)

// tenth is 0.1 read at run time: a test pinning exact bits takes an operand from it, since the
// compiler folds literal arithmetic into a constant without fusing.
var tenth = 0.1

func TestFractionalHex_Constructor(t *testing.T) {
	hextest.AssertFractionalHex(t, FracPt(10.9, -1.2), FractionalHex{Q: 10.9, R: -1.2})
}

func TestParseFractionalHex(t *testing.T) {
	t.Run("the form String prints", func(t *testing.T) {
		h, err := ParseFractionalHex("(10.90,-1.20)")
		assert.NoError(t, err)
		hextest.AssertFractionalHex(t, h, testFracHex)
	})
	t.Run("accepts whole coordinates", func(t *testing.T) {
		h, err := ParseFractionalHex("(2,-1)")
		assert.NoError(t, err)
		hextest.AssertFractionalHex(t, h, FracPt(2, -1))
	})
	t.Run("the parse error is wrapped", func(t *testing.T) {
		_, err := ParseFractionalHex("(a,1)")
		assert.ErrorContains(t, err, "hex: invalid q value")
		assert.ErrorIs(t, err, strconv.ErrSyntax)

		_, err = ParseFractionalHex("(1,1e999)")
		assert.ErrorContains(t, err, "hex: invalid r value")
		assert.ErrorIs(t, err, strconv.ErrRange)
	})
	t.Run("names the type of a malformed value", func(t *testing.T) {
		_, err := ParseFractionalHex("(1.5)")
		assert.ErrorContains(t, err, `hex: invalid fractional hex format "(1.5)"`)
	})
	t.Run("reprints the string it parsed", func(t *testing.T) {
		for _, h := range []FractionalHex{testFracHex, FracPt(0.25, -1.5), FracPt(-0.001, 3.14159)} {
			parsed, err := ParseFractionalHex(h.String())
			assert.NoError(t, err)
			assert.Equal(t, parsed.String(), h.String())
		}
	})
	t.Run("parses the non-finite values String prints", func(t *testing.T) {
		h, err := ParseFractionalHex(FracPt(math.NaN(), math.Inf(1)).String())
		assert.NoError(t, err)
		assert.True(t, math.IsNaN(h.Q))
		assert.True(t, math.IsInf(h.R, 1))
	})
}

func TestFractionalHex_S(t *testing.T) {
	assert.EqualDelta(t, testFracHex.S(), -9.7, geom.Delta)
}

func TestFractionalHex_QR(t *testing.T) {
	q, r := testFracHex.QR()
	assert.EqualDelta(t, q, 10.9, geom.Delta)
	assert.EqualDelta(t, r, -1.2, geom.Delta)
}

func TestFractionalHex_QRS(t *testing.T) {
	q, r, s := testFracHex.QRS()
	assert.EqualDelta(t, q, 10.9, geom.Delta)
	assert.EqualDelta(t, r, -1.2, geom.Delta)
	assert.EqualDelta(t, s, -9.7, geom.Delta)
}

func TestFractionalHex_Length(t *testing.T) {
	t.Run("steps from the origin", func(t *testing.T) {
		assert.EqualDelta(t, FracPt(3, 0).Length(), 3, geom.Delta)
		assert.EqualDelta(t, FracPt(0.5, -0.5).Length(), 0.5, geom.Delta)
		assert.EqualDelta(t, FracPt(0, 0).Length(), 0, geom.Delta)
	})
	t.Run("agrees with Hex.Length on whole coordinates", func(t *testing.T) {
		for _, h := range Pt(0, 0).Spiral(3) {
			assert.EqualDelta(t, h.Float().Length(), float64(h.Length()), geom.Delta)
		}
	})
}

func TestFractionalHex_Add(t *testing.T) {
	t.Run("vector sum", func(t *testing.T) {
		hextest.AssertFractionalHex(t, testFracHex.Add(FracPt(0.1, 1.2)), FracPt(11, 0))
	})
	t.Run("agrees with Hex.Add on whole coordinates", func(t *testing.T) {
		for _, whole := range Pt(0, 0).Spiral(2) {
			hextest.AssertFractionalHex(t, testHex.Float().Add(whole.Float()), testHex.Add(whole).Float(), whole.String()+": ")
		}
	})
}

func TestFractionalHex_Subtract(t *testing.T) {
	t.Run("vector difference", func(t *testing.T) {
		hextest.AssertFractionalHex(t, testFracHex.Subtract(FracPt(0.9, -0.2)), FracPt(10, -1))
	})
	t.Run("undoes Add", func(t *testing.T) {
		hextest.AssertFractionalHex(t, testFracHex.Add(FracPt(0.3, -2.7)).Subtract(FracPt(0.3, -2.7)), testFracHex)
	})
}

func TestFractionalHex_Multiply(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		hextest.AssertFractionalHex(t, testFracHex.Multiply(2), FracPt(21.8, -2.4))
	})
	t.Run("zero", func(t *testing.T) {
		hextest.AssertFractionalHex(t, testFracHex.Multiply(0), FracPt(0, 0))
	})
	t.Run("negative", func(t *testing.T) {
		hextest.AssertFractionalHex(t, testFracHex.Multiply(-1), FracPt(-10.9, 1.2))
	})
	t.Run("rounds the product before a caller adds it", func(t *testing.T) {
		sum := FracPt(-0.01, -0.01).Add(FracPt(tenth, tenth).Multiply(tenth))
		assert.Equal(t, sum.Q, 1.734723475976807e-18)
		assert.Equal(t, sum.R, 1.734723475976807e-18)
	})
}

func TestFractionalHex_Lerp(t *testing.T) {
	t.Run("interpolates both coordinates", func(t *testing.T) {
		hextest.AssertFractionalHex(t, testFracHex.Lerp(FracPt(12, 0), 0.1), FracPt(11.01, -1.08))
	})
	t.Run("rounds the product before adding", func(t *testing.T) {
		lerp := FracPt(tenth, tenth).Lerp(FracPt(0.2, 0.2), 0.1)
		assert.Equal(t, lerp.Q, 0.11000000000000001)
		assert.Equal(t, lerp.R, 0.11000000000000001)
	})
	t.Run("starts at the receiver and ends at the argument", func(t *testing.T) {
		hextest.AssertFractionalHex(t, testFracHex.Lerp(FracPt(12, 0), 0), testFracHex)
		hextest.AssertFractionalHex(t, testFracHex.Lerp(FracPt(12, 0), 1), FracPt(12, 0))
	})
	t.Run("the distances to both ends add up", func(t *testing.T) {
		end := FracPt(12, 0)
		for i := range 11 {
			h := testFracHex.Lerp(end, float64(i)/10)
			assert.EqualDelta(t, testFracHex.DistanceTo(h)+h.DistanceTo(end), testFracHex.DistanceTo(end), geom.Delta, strconv.Itoa(i)+": ")
		}
	})
}

func BenchmarkFractionalHex_Lerp(b *testing.B) {
	end := FracPt(12, 0)

	for b.Loop() {
		sinkFracHex = testFracHex.Lerp(end, 0.3)
	}
}

func TestFractionalHex_Turn(t *testing.T) {
	h := FracPt(1.5, -0.5)

	t.Run("zero and six steps are the identity", func(t *testing.T) {
		hextest.AssertFractionalHex(t, h.Turn(0), h)
		hextest.AssertFractionalHex(t, h.Turn(6), h)
	})
	t.Run("one step forward", func(t *testing.T) {
		hextest.AssertFractionalHex(t, h.Turn(1), FracPt(0.5, 1))
	})
	t.Run("one step back", func(t *testing.T) {
		hextest.AssertFractionalHex(t, h.Turn(-1), FracPt(1, -1.5))
	})
	t.Run("there and back returns to the start", func(t *testing.T) {
		hextest.AssertFractionalHex(t, testFracHex.Turn(1).Turn(-1), testFracHex)
	})
	t.Run("keeps the length", func(t *testing.T) {
		for steps := -6; steps <= 6; steps++ {
			assert.EqualDelta(t, testFracHex.Turn(steps).Length(), testFracHex.Length(), geom.Delta, strconv.Itoa(steps)+": ")
		}
	})
	t.Run("agrees with Hex.Turn on whole coordinates", func(t *testing.T) {
		for _, whole := range Pt(0, 0).Spiral(2) {
			for steps := -6; steps <= 6; steps++ {
				hextest.AssertFractionalHex(t, whole.Float().Turn(steps), whole.Turn(steps).Float(), whole.String()+" by "+strconv.Itoa(steps)+": ")
			}
		}
	})
	t.Run("steps add up", func(t *testing.T) {
		for first := -6; first <= 6; first++ {
			for second := -6; second <= 6; second++ {
				hextest.AssertFractionalHex(t, testFracHex.Turn(first).Turn(second), testFracHex.Turn(first+second), strconv.Itoa(first)+" then "+strconv.Itoa(second)+": ")
			}
		}
	})
	t.Run("a full turn returns the same bits", func(t *testing.T) {
		assert.Equal(t, testFracHex.Turn(6), testFracHex)
		assert.Equal(t, testFracHex.Turn(-12), testFracHex)
	})
	t.Run("rounds to the turned hex", func(t *testing.T) {
		for steps := -6; steps <= 6; steps++ {
			hextest.AssertHex(t, testFracHex.Turn(steps).Round(), testFracHex.Round().Turn(steps), strconv.Itoa(steps)+": ")
		}
	})
}

func BenchmarkFractionalHex_Turn(b *testing.B) {
	for b.Loop() {
		sinkFracHex = testFracHex.Turn(5)
	}
}

func TestFractionalHex_TurnAround(t *testing.T) {
	center := FracPt(1.25, 0.75)
	h := FracPt(3.5, -0.5)

	t.Run("zero and six steps are the identity", func(t *testing.T) {
		hextest.AssertFractionalHex(t, h.TurnAround(center, 0), h)
		hextest.AssertFractionalHex(t, h.TurnAround(center, 6), h)
	})
	t.Run("keeps the distance from the center", func(t *testing.T) {
		for steps := 1; steps <= 5; steps++ {
			assert.EqualDelta(t, center.DistanceTo(h.TurnAround(center, steps)), center.DistanceTo(h), geom.Delta, strconv.Itoa(steps)+": ")
		}
	})
	t.Run("there and back returns to the start", func(t *testing.T) {
		hextest.AssertFractionalHex(t, h.TurnAround(center, 1).TurnAround(center, -1), h)
	})
	t.Run("around the origin is Turn", func(t *testing.T) {
		hextest.AssertFractionalHex(t, h.TurnAround(FracPt(0, 0), 1), h.Turn(1))
		hextest.AssertFractionalHex(t, h.TurnAround(FracPt(0, 0), -1), h.Turn(-1))
	})
	t.Run("agrees with Hex.TurnAround on whole coordinates", func(t *testing.T) {
		around := Pt(1, 1)
		for _, whole := range Pt(0, 0).Spiral(2) {
			for steps := -6; steps <= 6; steps++ {
				hextest.AssertFractionalHex(t, whole.Float().TurnAround(around.Float(), steps), whole.TurnAround(around, steps).Float(), whole.String()+" by "+strconv.Itoa(steps)+": ")
			}
		}
	})
}

func TestFractionalHex_ReflectQ(t *testing.T) {
	t.Run("swaps r and s", func(t *testing.T) {
		hextest.AssertFractionalHex(t, testFracHex.ReflectQ(), FracPt(10.9, -9.7))
	})
	t.Run("keeps q", func(t *testing.T) {
		assert.Equal(t, testFracHex.ReflectQ().Q, testFracHex.Q)
	})
	t.Run("twice is the identity", func(t *testing.T) {
		hextest.AssertFractionalHex(t, testFracHex.ReflectQ().ReflectQ(), testFracHex)
	})
	t.Run("agrees with Hex.ReflectQ on whole coordinates", func(t *testing.T) {
		for _, whole := range Pt(0, 0).Spiral(2) {
			hextest.AssertFractionalHex(t, whole.Float().ReflectQ(), whole.ReflectQ().Float(), whole.String()+": ")
		}
	})
}

func TestFractionalHex_ReflectR(t *testing.T) {
	t.Run("swaps q and s", func(t *testing.T) {
		hextest.AssertFractionalHex(t, testFracHex.ReflectR(), FracPt(-9.7, -1.2))
	})
	t.Run("keeps r", func(t *testing.T) {
		assert.Equal(t, testFracHex.ReflectR().R, testFracHex.R)
	})
	t.Run("twice is the identity", func(t *testing.T) {
		hextest.AssertFractionalHex(t, testFracHex.ReflectR().ReflectR(), testFracHex)
	})
	t.Run("agrees with Hex.ReflectR on whole coordinates", func(t *testing.T) {
		for _, whole := range Pt(0, 0).Spiral(2) {
			hextest.AssertFractionalHex(t, whole.Float().ReflectR(), whole.ReflectR().Float(), whole.String()+": ")
		}
	})
}

func TestFractionalHex_ReflectS(t *testing.T) {
	t.Run("swaps q and r", func(t *testing.T) {
		hextest.AssertFractionalHex(t, testFracHex.ReflectS(), FracPt(-1.2, 10.9))
	})
	t.Run("keeps s", func(t *testing.T) {
		assert.EqualDelta(t, testFracHex.ReflectS().S(), testFracHex.S(), geom.Delta)
	})
	t.Run("twice is the identity", func(t *testing.T) {
		hextest.AssertFractionalHex(t, testFracHex.ReflectS().ReflectS(), testFracHex)
	})
	t.Run("agrees with Hex.ReflectS on whole coordinates", func(t *testing.T) {
		for _, whole := range Pt(0, 0).Spiral(2) {
			hextest.AssertFractionalHex(t, whole.Float().ReflectS(), whole.ReflectS().Float(), whole.String()+": ")
		}
	})
}

func TestFractionalHex_DistanceTo(t *testing.T) {
	a := FracPt(1.5, -0.5)
	b := FracPt(3.5, -0.5)

	t.Run("steps between hexes", func(t *testing.T) {
		assert.EqualDelta(t, a.DistanceTo(b), 2, geom.Delta)
	})
	t.Run("symmetric", func(t *testing.T) {
		assert.EqualDelta(t, a.DistanceTo(b), b.DistanceTo(a), geom.Delta)
	})
	t.Run("zero to itself", func(t *testing.T) {
		assert.EqualDelta(t, a.DistanceTo(a), 0, geom.Delta)
	})
}

func TestFractionalHex_Equal(t *testing.T) {
	t.Run("same coordinates", func(t *testing.T) {
		assert.True(t, testFracHex.Equal(FracPt(10.9, -1.2)))
	})
	t.Run("within the tolerance", func(t *testing.T) {
		assert.True(t, testFracHex.Equal(FracPt(10.9+geom.Delta/10, -1.2)))
	})
	t.Run("one coordinate differs", func(t *testing.T) {
		assert.False(t, testFracHex.Equal(FracPt(10.91, -1.2)))
		assert.False(t, testFracHex.Equal(FracPt(10.9, -1.21)))
	})
}

func TestFractionalHex_IsZero(t *testing.T) {
	t.Run("origin", func(t *testing.T) {
		assert.True(t, FracPt(0, 0).IsZero())
	})
	t.Run("within the tolerance", func(t *testing.T) {
		assert.True(t, FracPt(geom.Delta/10, 0).IsZero())
	})
	t.Run("any other hex", func(t *testing.T) {
		assert.False(t, testFracHex.IsZero())
	})
}

func TestFractionalHex_Point(t *testing.T) {
	assert.Equal(t, testFracHex.Point(), geom.Pt(10.9, -1.2))
}

func TestFractionalHex_Round(t *testing.T) {
	t.Run("nearest hex", func(t *testing.T) {
		hextest.AssertHex(t, FracPt(10.9, 16.2).Round(), Pt(11, 16))
		hextest.AssertHex(t, FracPt(10.5001, 16.4999).Round(), Pt(11, 16))
	})
	t.Run("keeps the cube constraint near a corner", func(t *testing.T) {
		hextest.AssertHex(t, FracPt(10.50000001, 16.5000001).Round(), Pt(10, 17))
		hextest.AssertHex(t, FracPt(10.500001, 16.500000001).Round(), Pt(11, 16))
	})
	t.Run("panics for a non-finite coordinate", func(t *testing.T) {
		assert.Panics(t, func() {
			sinkHex = FracPt(math.NaN(), 0).Round()
		})
		assert.Panics(t, func() {
			sinkHex = FracPt(0, math.Inf(-1)).Round()
		})
	})
	t.Run("a whole hex rounds to itself", func(t *testing.T) {
		for _, h := range testHex.Spiral(3) {
			hextest.AssertHex(t, h.Float().Round(), h)
		}
	})
	t.Run("stays in the hex anywhere short of its edge", func(t *testing.T) {
		for _, h := range testHex.Spiral(2) {
			for _, direction := range Directions() {
				inside := h.Float().Lerp(h.Neighbor(direction).Float(), 0.49)
				hextest.AssertHex(t, inside.Round(), h, h.String()+" towards "+direction.String()+": ")
			}
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkHex = testFracHex.Round()
		}), 0.0)
	})
}

func BenchmarkFractionalHex_Round(b *testing.B) {
	for b.Loop() {
		sinkHex = testFracHex.Round()
	}
}

func FuzzFractionalHex_Round(f *testing.F) {
	f.Add(10.9, -1.2)
	f.Add(0.5, 0.5)
	f.Add(-0.5, 0.25)
	f.Add(1.0/3, 1.0/3)
	f.Add(-999999.5, 999999.25)

	f.Fuzz(func(t *testing.T, q, r float64) {
		if math.IsNaN(q) || math.IsNaN(r) || math.Abs(q) > 1e6 || math.Abs(r) > 1e6 {
			t.Skip()
		}

		h := FracPt(q, r)
		rounded := h.Round()
		message := fmt.Sprintf("(%v,%v) → %s: ", q, r, rounded)

		for _, neighbor := range rounded.Neighbors() {
			assert.True(t, squaredDistanceOf(h, rounded.Float()) <= squaredDistanceOf(h, neighbor.Float())+geom.Delta, message)
		}
		hextest.AssertHex(t, rounded.Float().Round(), rounded, message)
	})
}

func TestFractionalHex_String(t *testing.T) {
	t.Run("two decimals", func(t *testing.T) {
		assert.Equal(t, testFracHex.String(), "(10.90,-1.20)")
	})
	t.Run("no sign on a value that rounds to zero", func(t *testing.T) {
		assert.Equal(t, FracPt(-0.001, math.Copysign(0, -1)).String(), "(0.00,0.00)")
	})
	t.Run("prints each coordinate as geom does", func(t *testing.T) {
		for _, value := range []float64{0, -0.001, 0.004, -0.004, 0.005, -0.005, 1.005, -1.005, 2.675, 1e6, -1e15, math.MaxFloat64, math.Inf(1), math.Inf(-1), math.NaN()} {
			assert.Equal(t, FracPt(value, value).String(), "("+geom.String(value)+","+geom.String(value)+")", geom.String(value)+": ")
		}
	})
	t.Run("allocates the string alone", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkString = testFracHex.String()
		}), 1.0)
	})
}

func BenchmarkFractionalHex_String(b *testing.B) {
	for b.Loop() {
		sinkString = testFracHex.String()
	}
}

func TestFractionalHex_JSON(t *testing.T) {
	data, err := json.Marshal(testFracHex)
	assert.NoError(t, err)
	assert.Equal(t, string(data), `{"q":10.9,"r":-1.2}`)

	var decoded FractionalHex
	assert.NoError(t, json.Unmarshal(data, &decoded))
	hextest.AssertFractionalHex(t, decoded, testFracHex)
}

// squaredDistanceOf returns the squared Euclidean distance between two fractional hexes, in
// units of the distance between two neighboring hex centers.
func squaredDistanceOf(a, b FractionalHex) float64 {
	delta := a.Subtract(b)

	return delta.Q*delta.Q + delta.Q*delta.R + delta.R*delta.R
}

func ExampleFracPt() {
	fmt.Println(FracPt(1.4, -1.8))
	// Output: (1.40,-1.80)
}

func ExampleFractionalHex_Lerp() {
	fmt.Println(FracPt(0, 0).Lerp(FracPt(3, -1), 0.5))
	// Output: (1.50,-0.50)
}

func ExampleFractionalHex_Turn() {
	h := FracPt(1.5, -0.5)

	fmt.Println(h.Turn(1))
	fmt.Println(h.TurnAround(FracPt(1, 0), 3))
	// Output:
	// (0.50,1.00)
	// (0.50,0.50)
}

func ExampleFractionalHex_Round() {
	fmt.Println(FracPt(1.2, -1.9).Round())
	fmt.Println(Pt(2, -1).Float().Round())
	// Output:
	// (1,-2)
	// (2,-1)
}
