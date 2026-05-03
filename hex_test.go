package hex

import (
	"testing"

	"github.com/gravitton/assert"
	geom "github.com/gravitton/geometry"
)

var (
	testHex     = Hex{-1, 3}
	testHexZero = Hex{0, 0}
	testFracHex = FractionalHex{10.9, -1.2}
)

func TestHex_New(t *testing.T) {
	AssertHex(t, testHex, -1, 3)
	AssertHex(t, Pt(1, -2), 1, -2)
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

func TestHex_Add(t *testing.T) {
	AssertHex(t, testHex.Add(Pt(3, -2)), 2, 1)
}

func TestHex_Subtract(t *testing.T) {
	AssertHex(t, testHex.Subtract(Pt(3, -2)), -4, 5)
}

func TestHex_Multiply(t *testing.T) {
	AssertHex(t, testHex.Multiply(3), -3, 9)
}

func TestHex_Length(t *testing.T) {
	assert.Equal(t, testHex.Length(), 3)
}

func TestHex_DistanceTo(t *testing.T) {
	assert.Equal(t, testHex.DistanceTo(Pt(0, 3)), 1)
	assert.Equal(t, testHex.DistanceTo(Pt(1, 6)), 5)
}

func TestHex_Neighbors(t *testing.T) {
	assert.Equal(t, testHex.Neighbors(), []Hex{{0, 3}, {0, 2}, {-1, 2}, {-2, 3}, {-2, 4}, {-1, 4}})
}

func TestHex_Neighbor(t *testing.T) {
	AssertHex(t, testHex.Neighbor(SMinus), 0, 3)
	AssertHex(t, testHex.Neighbor(QPlus), 0, 2)
	AssertHex(t, testHex.Neighbor(RMinus), -1, 2)
	AssertHex(t, testHex.Neighbor(SPlus), -2, 3)
	AssertHex(t, testHex.Neighbor(QMinus), -2, 4)
	AssertHex(t, testHex.Neighbor(RPlus), -1, 4)
}

func TestHex_Range(t *testing.T) {
	assert.Equal(t, testHexZero.Range(-1), nil)
	assert.Equal(t, testHexZero.Range(0), []Hex{{-0, 0}})
	assert.Equal(t, testHexZero.Range(1), []Hex{{-1, 0}, {-1, 1}, {0, -1}, {0, 0}, {0, 1}, {1, -1}, {1, 0}})
	assert.Equal(t, len(testHexZero.Range(2)), 19)
	assert.Equal(t, len(testHexZero.Range(3)), 37)
}

func TestHex_Line(t *testing.T) {
	assert.Equal(t, testHexZero.Line(Pt(3, 0)), []Hex{{0, 0}, {1, 0}, {2, 0}, {3, 0}})
	assert.Equal(t, testHexZero.Line(Pt(2, -1)), []Hex{{0, 0}, {1, 0}, {2, -1}})
	assert.Equal(t, testHexZero.Line(Pt(-2, 1)), []Hex{{0, 0}, {-1, 1}, {-2, 1}})
	assert.Equal(t, testHexZero.Line(Pt(4, -2)), []Hex{{0, 0}, {1, 0}, {2, -1}, {3, -1}, {4, -2}})
	assert.Equal(t, testHexZero.Line(Pt(-4, 2)), []Hex{{0, 0}, {-1, 1}, {-2, 1}, {-3, 2}, {-4, 2}})
}

func TestHex_HasLineOfSight(t *testing.T) {
	// Clear line with no blockers
	assert.True(t, testHexZero.HasLineOfSight(Pt(3, 0), nil))
	assert.True(t, testHexZero.HasLineOfSight(Pt(3, 0), []Hex{}))

	// Same hex always has line of sight to itself
	assert.True(t, testHexZero.HasLineOfSight(testHexZero, nil))

	// Blocked in the middle
	assert.False(t, testHexZero.HasLineOfSight(Pt(3, 0), []Hex{Pt(1, 0)}))
	assert.False(t, testHexZero.HasLineOfSight(Pt(3, 0), []Hex{Pt(2, 0)}))

	// Blocker beyond the target does not affect visibility
	assert.True(t, testHexZero.HasLineOfSight(Pt(2, 0), []Hex{Pt(3, 0)}))

	// Target in the blocking list is still visible (can see into, not through)
	assert.True(t, testHexZero.HasLineOfSight(Pt(3, 0), []Hex{Pt(3, 0)}))
}

func TestHex_FieldOfView(t *testing.T) {
	candidates := testHexZero.Range(3)

	// No blocking: all candidates are visible
	assert.Equal(t, len(testHexZero.FieldOfView(candidates, nil)), len(candidates))
	assert.Equal(t, len(testHexZero.FieldOfView(candidates, []Hex{})), len(candidates))

	// Immediate neighbors (distance <= 1) are always visible regardless of blocking
	blocking := testHexZero.Neighbors()
	visible := testHexZero.FieldOfView(candidates, blocking)
	for _, v := range visible {
		assert.True(t, testHexZero.DistanceTo(v) <= 1)
	}

	// A gap in the blocking ring allows seeing through it
	// H(1,0) is left unblocked, making H(2,0) and H(3,0) visible in that direction
	partialBlocking := []Hex{Pt(0, -1), Pt(1, -1), Pt(-1, 0), Pt(-1, 1), Pt(0, 1)}
	visible = testHexZero.FieldOfView([]Hex{Pt(2, 0), Pt(3, 0), Pt(0, -2)}, partialBlocking)
	assert.Equal(t, len(visible), 2)
	assert.Equal(t, visible[0], Pt(2, 0))
	assert.Equal(t, visible[1], Pt(3, 0))
}

func TestHex_Point(t *testing.T) {
	assert.Equal(t, testHex.Point(), geom.Pt(-1, 3))
	assert.Equal(t, testHexZero.Point(), geom.Pt(0, 0))
}

func TestHex_IsZero(t *testing.T) {
	assert.True(t, testHexZero.IsZero())
	assert.False(t, testHex.IsZero())
	assert.False(t, Pt(0, 1).IsZero())
	assert.False(t, Pt(1, 0).IsZero())
}

func TestHex_String(t *testing.T) {
	assert.Equal(t, testHex.String(), "(-1,3)")
}

func TestFractionalHex_New(t *testing.T) {
	AssertFracHex(t, testFracHex, 10.9, -1.2)
	AssertFracHex(t, FracPt(-1.2, -2.5), -1.2, -2.5)
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

func TestFractionalHex_Point(t *testing.T) {
	assert.Equal(t, testFracHex.Point(), geom.Pt(10.9, -1.2))
}

func TestFractionalHex_Lerp(t *testing.T) {
	AssertFracHex(t, testFracHex.Lerp(FracPt(12, 0), 0.1), 11.01, -1.08)
}

func TestFractionalHex_Round(t *testing.T) {
	AssertHex(t, FracPt(10.9, 16.2).Round(), 11, 16)
	AssertHex(t, FracPt(10.5001, 16.4999).Round(), 11, 16)
	AssertHex(t, FracPt(10.50000001, 16.5000001).Round(), 10, 17)
	AssertHex(t, FracPt(10.500001, 16.500000001).Round(), 11, 16)
}

func TestFractionalHex_String(t *testing.T) {
	assert.Equal(t, testFracHex.String(), "(10.90,-1.20)")
}

func TestHex_Multiply_Zero(t *testing.T) {
	AssertHex(t, testHex.Multiply(0), 0, 0)
}

func TestHex_Multiply_Negative(t *testing.T) {
	AssertHex(t, testHex.Multiply(-1), 1, -3)
	AssertHex(t, testHex.Multiply(-2), 2, -6)
}

func TestHex_DistanceTo_Symmetric(t *testing.T) {
	a := Pt(3, -2)
	b := Pt(-1, 5)
	assert.Equal(t, a.DistanceTo(b), b.DistanceTo(a))
}

func TestHex_Range_Large(t *testing.T) {
	for _, n := range []int{5, 10, 20} {
		expected := 1 + 3*n*(n+1)
		assert.Equal(t, len(testHexZero.Range(n)), expected)
	}
}

func TestHex_Line_SelfToSelf(t *testing.T) {
	assert.Equal(t, testHexZero.Line(testHexZero), []Hex{{0, 0}})
	assert.Equal(t, testHex.Line(testHex), []Hex{testHex})
}

func TestHex_FieldOfView_EmptyCandidates(t *testing.T) {
	assert.Equal(t, len(testHexZero.FieldOfView(nil, nil)), 0)
	assert.Equal(t, len(testHexZero.FieldOfView([]Hex{}, nil)), 0)
}

func TestHex_Ring(t *testing.T) {
	// negative radius returns nil
	assert.Equal(t, testHexZero.Ring(-1), nil)

	// radius 0 returns just the center
	assert.Equal(t, testHexZero.Ring(0), []Hex{{0, 0}})

	// radius 1: 6 hexes all at distance 1
	ring1 := testHexZero.Ring(1)
	assert.Equal(t, len(ring1), 6)
	// specific counterclockwise order starting from QMinus direction
	assert.Equal(t, ring1, []Hex{{-1, 1}, {0, 1}, {1, 0}, {1, -1}, {0, -1}, {-1, 0}})
	for _, h := range ring1 {
		assert.Equal(t, testHexZero.DistanceTo(h), 1)
	}

	// radius 2: 12 hexes all at distance 2
	ring2 := testHexZero.Ring(2)
	assert.Equal(t, len(ring2), 12)
	for _, h := range ring2 {
		assert.Equal(t, testHexZero.DistanceTo(h), 2)
	}

	// count is always 6 * radius for radius >= 1
	for _, radius := range []int{3, 5, 10} {
		assert.Equal(t, len(testHexZero.Ring(radius)), 6*radius)
	}

	// non-origin center
	center := Pt(2, -3)
	ring := center.Ring(1)
	assert.Equal(t, len(ring), 6)
	for _, h := range ring {
		assert.Equal(t, center.DistanceTo(h), 1)
	}
}

func TestHex_Spiral(t *testing.T) {
	assert.Equal(t, testHexZero.Spiral(-1), nil)
	assert.Equal(t, testHexZero.Spiral(0), []Hex{{0, 0}})

	// radius 1: center + ring1 = 7 hexes
	spiral1 := testHexZero.Spiral(1)
	assert.Equal(t, len(spiral1), 7)
	assert.Equal(t, spiral1[0], testHexZero) // center comes first

	// radius 2: 1 + 6 + 12 = 19 hexes, same count as Range(2)
	spiral2 := testHexZero.Spiral(2)
	assert.Equal(t, len(spiral2), len(testHexZero.Range(2)))

	// each hex in the spiral is within the radius
	for _, h := range spiral2 {
		assert.True(t, testHexZero.DistanceTo(h) <= 2)
	}

	// count matches Range for several radii
	for _, n := range []int{3, 5} {
		assert.Equal(t, len(testHexZero.Spiral(n)), len(testHexZero.Range(n)))
	}
}

func TestHex_Rotate(t *testing.T) {
	h := Pt(3, 0)

	// 0 steps is identity
	AssertHex(t, h.Rotate(0), h.Q, h.R)

	// 6 steps is identity
	AssertHex(t, h.Rotate(6), h.Q, h.R)

	// one CW step: (q,r) → (-r, q+r)
	AssertHex(t, h.Rotate(1), -h.R, h.Q+h.R)

	// one CCW step: (q,r) → (q+r, -q)
	AssertHex(t, h.Rotate(-1), h.Q+h.R, -h.Q)

	// rotating CW then CCW returns to start
	AssertHex(t, h.Rotate(1).Rotate(-1), h.Q, h.R)
}

func TestHex_RotateAround(t *testing.T) {
	center := Pt(1, 1)
	h := Pt(3, 0)

	// 0 steps is identity
	AssertHex(t, h.RotateAround(center, 0), h.Q, h.R)

	// 6 steps is identity
	AssertHex(t, h.RotateAround(center, 6), h.Q, h.R)

	// distance from center is preserved after rotation
	dist := center.DistanceTo(h)
	for steps := 1; steps <= 5; steps++ {
		assert.Equal(t, center.DistanceTo(h.RotateAround(center, steps)), dist)
	}

	// rotating CW then CCW (-1) returns to start
	AssertHex(t, h.RotateAround(center, 1).RotateAround(center, -1), h.Q, h.R)

	// RotateAround origin equals Rotate
	AssertHex(t, h.RotateAround(testHexZero, 1), h.Rotate(1).Q, h.Rotate(1).R)
	AssertHex(t, h.RotateAround(testHexZero, -1), h.Rotate(-1).Q, h.Rotate(-1).R)
}

func TestHex_ReflectQ(t *testing.T) {
	// double reflection is identity
	h := Pt(2, 1)
	AssertHex(t, h.ReflectQ().ReflectQ(), h.Q, h.R)

	// q coordinate is preserved
	assert.Equal(t, h.ReflectQ().Q, h.Q)

	// specific cases
	AssertHex(t, Pt(0, 0).ReflectQ(), 0, 0)
	AssertHex(t, Pt(1, 0).ReflectQ(), 1, -1) // r' = -q-r = -1
	AssertHex(t, Pt(2, 0).ReflectQ(), 2, -2)
	AssertHex(t, Pt(2, -2).ReflectQ(), 2, 0) // symmetric with Coord(2,0)
}

func TestHex_ReflectR(t *testing.T) {
	// double reflection is identity
	h := Pt(2, 1)
	AssertHex(t, h.ReflectR().ReflectR(), h.Q, h.R)

	// r coordinate is preserved
	assert.Equal(t, h.ReflectR().R, h.R)

	// specific cases
	AssertHex(t, Pt(0, 0).ReflectR(), 0, 0)
	AssertHex(t, Pt(1, 0).ReflectR(), -1, 0) // q' = -q-r = -1
	AssertHex(t, Pt(2, 0).ReflectR(), -2, 0)
	AssertHex(t, Pt(-2, 0).ReflectR(), 2, 0) // symmetric with Coord(2,0)
}

func TestHex_ReflectS(t *testing.T) {
	// double reflection is identity
	h := Pt(2, 1)
	AssertHex(t, h.ReflectS().ReflectS(), h.Q, h.R)

	// s coordinate is preserved
	assert.Equal(t, h.ReflectS().S(), h.S())

	// specific cases: q and r are swapped
	AssertHex(t, Pt(0, 0).ReflectS(), 0, 0)
	AssertHex(t, Pt(1, 0).ReflectS(), 0, 1)
	AssertHex(t, Pt(3, -1).ReflectS(), -1, 3)
	AssertHex(t, Pt(-1, 3).ReflectS(), 3, -1) // symmetric with Coord(3,-1)
}
