package hex

import (
	"testing"

	"github.com/gravitton/assert"
	geom "github.com/gravitton/geometry"
)

func TestHex_New(t *testing.T) {
	testHex(t, H(10, 16), 10, 16)
	testHex(t, H(-1, -2), -1, -2)
}

func TestHex_S(t *testing.T) {
	assert.Equal(t, H(10, 16).S(), -26)
	assert.Equal(t, H(-1, 2).S(), -1)
}

//func TestHex_Add(t *testing.T) {
//	testHex(t, H(1, 2).Add(H(3, -2)), 4, 0)
//	testHex(t, P(1, 2).AddXY(3, -2), 4, 0)
//	testHex(t, P(0.4, -0.25).Add(V(100.1, -0.1)), 100.5, -0.35)
//	testHex(t, P(0.4, -0.25).AddXY(100.1, -0.1), 100.5, -0.35)
//}
//
//func TestHex_Subtract(t *testing.T) {
//	testVector(t, P(1, 2).Subtract(P(3, -3)), -2, 5)
//	testVector(t, P(0.4, -0.25).Subtract(P(100.1, -0.1)), -99.7, -0.15)
//}

func TestHex_Multiply(t *testing.T) {
	testHex(t, H(1, 2).Multiply(3), 3, 6)
}

func TestHex_Length(t *testing.T) {
	assert.Equal(t, H(1, -2).Length(), 2)
}

func TestHex_DistanceTo(t *testing.T) {
	assert.Equal(t, H(1, -2).DistanceTo(H(0, -2)), 1)
	assert.Equal(t, H(1, -2).DistanceTo(H(5, -1)), 5)
}

func TestHex_Range(t *testing.T) {
	assert.Equal(t, H(0, 0).Range(0), nil)
	assert.Equal(t, H(0, 0).Range(1), []Hex{{-1, 0}, {-1, 1}, {0, -1}, {0, 0}, {0, 1}, {1, -1}, {1, 0}})
	assert.Equal(t, len(H(0, 0).Range(2)), 19)
	assert.Equal(t, len(H(0, 0).Range(3)), 37)
}

func TestHex_String(t *testing.T) {
	assert.Equal(t, H(10, 16).String(), "(+10,+16)")
	assert.Equal(t, H(1, -2).String(), "(+1,-2)")
}

func TestFractionalHex_New(t *testing.T) {
	testFractionalHex(t, F(10.9, 16.2), 10.9, 16.2)
	testFractionalHex(t, F(-1.2, -2.5), -1.2, -2.5)
}

func TestFractionalHex_S(t *testing.T) {
	assert.Equal(t, F(10.9, 16.2).S(), -27.1)
	assert.Equal(t, F(-1.2, -2.5).S(), 3.7)
}

func TestFractionalHex_ToPoint(t *testing.T) {
	assert.Equal(t, F(10.9, 16.2).ToPoint(), geom.P(10.9, 16.2))
	assert.Equal(t, F(-1.2, -2.5).ToPoint(), geom.P(-1.2, -2.5))
}

func TestFractionalHex_Round(t *testing.T) {
	testHex(t, F(10.9, 16.2).Round(), 11, 16)
	testHex(t, F(10.5001, 16.4999).Round(), 11, 16)
	testHex(t, F(10.50000001, 16.5000001).Round(), 10, 17)
	testHex(t, F(10.500001, 16.500000001).Round(), 11, 16)
}

func TestFractionalHex_String(t *testing.T) {
	assert.Equal(t, F(10.9, 16.2).String(), "(+10.90,+16.20)")
	assert.Equal(t, F(-1.2, -2.5).String(), "(-1.20,-2.50)")
}

func testHex(t *testing.T, h Hex, q, r int) {
	t.Helper()

	assert.Equal(t, h.Q, q)
	assert.Equal(t, h.R, r)
}

func testFractionalHex(t *testing.T, h FractionalHex, q, r float64) {
	t.Helper()

	assert.Equal(t, h.Q, q)
	assert.Equal(t, h.R, r)
}
