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
	AssertHex(t, H(1, -2), 1, -2)
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
	AssertHex(t, testHex.Add(H(3, -2)), 2, 1)
}

func TestHex_Subtract(t *testing.T) {
	AssertHex(t, testHex.Subtract(H(3, -2)), -4, 5)
}

func TestHex_Multiply(t *testing.T) {
	AssertHex(t, testHex.Multiply(3), -3, 9)
}

func TestHex_Length(t *testing.T) {
	assert.Equal(t, testHex.Length(), 3)
}

func TestHex_DistanceTo(t *testing.T) {
	assert.Equal(t, testHex.DistanceTo(H(0, 3)), 1)
	assert.Equal(t, testHex.DistanceTo(H(1, 6)), 5)
}

func TestHex_Neighbors(t *testing.T) {
	assert.Equal(t, testHex.Neighbors(), []Hex{{0, 3}, {0, 2}, {-1, 2}, {-2, 3}, {-2, 4}, {-1, 4}})
}

func TestHex_Neighbor(t *testing.T) {
	AssertHex(t, testHex.Neighbor(DirectionSMinus), 0, 3)
	AssertHex(t, testHex.Neighbor(DirectionQPlus), 0, 2)
	AssertHex(t, testHex.Neighbor(DirectionRMinus), -1, 2)
	AssertHex(t, testHex.Neighbor(DirectionSPlus), -2, 3)
	AssertHex(t, testHex.Neighbor(DirectionQMinus), -2, 4)
	AssertHex(t, testHex.Neighbor(DirectionRPlus), -1, 4)
}

func TestHex_Range(t *testing.T) {
	assert.Equal(t, testHexZero.Range(-1), nil)
	assert.Equal(t, testHexZero.Range(0), []Hex{{-0, 0}})
	assert.Equal(t, testHexZero.Range(1), []Hex{{-1, 0}, {-1, 1}, {0, -1}, {0, 0}, {0, 1}, {1, -1}, {1, 0}})
	assert.Equal(t, len(testHexZero.Range(2)), 19)
	assert.Equal(t, len(testHexZero.Range(3)), 37)
}

func TestHex_Line(t *testing.T) {
	assert.Equal(t, testHexZero.Line(H(3, 0)), []Hex{{0, 0}, {1, 0}, {2, 0}, {3, 0}})
	assert.Equal(t, testHexZero.Line(H(2, -1)), []Hex{{0, 0}, {1, 0}, {2, -1}})
	assert.Equal(t, testHexZero.Line(H(-2, 1)), []Hex{{0, 0}, {-1, 1}, {-2, 1}})
	assert.Equal(t, testHexZero.Line(H(4, -2)), []Hex{{0, 0}, {1, 0}, {2, -1}, {3, -1}, {4, -2}})
	assert.Equal(t, testHexZero.Line(H(-4, 2)), []Hex{{0, 0}, {-1, 1}, {-2, 1}, {-3, 2}, {-4, 2}})
}

func TestHex_HasLineOfSight(t *testing.T) {
	// TODO: add test
	//assert.True(t, testHexZero.HasLineOfSight())
}

func TestHex_FieldOfView(t *testing.T) {
	// TODO: add test
	//assert.Equal(t, testHexZero.FieldOfView())
}

func TestHex_String(t *testing.T) {
	assert.Equal(t, testHex.String(), "(-1,3)")
}

func TestFractionalHex_New(t *testing.T) {
	AssertFracHex(t, testFracHex, 10.9, -1.2)
	AssertFracHex(t, F(-1.2, -2.5), -1.2, -2.5)
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
	AssertFracHex(t, testFracHex.Lerp(F(12, 0), 0.1), 11.01, -1.08)
}

func TestFractionalHex_Round(t *testing.T) {
	AssertHex(t, F(10.9, 16.2).Round(), 11, 16)
	AssertHex(t, F(10.5001, 16.4999).Round(), 11, 16)
	AssertHex(t, F(10.50000001, 16.5000001).Round(), 10, 17)
	AssertHex(t, F(10.500001, 16.500000001).Round(), 11, 16)
}

func TestFractionalHex_String(t *testing.T) {
	assert.Equal(t, testFracHex.String(), "(10.90,-1.20)")
}
