package hex

import (
	"encoding/json"
	"testing"

	"github.com/gravitton/assert"
	geom "github.com/gravitton/geometry"
)

var testFracHex = FractionalHex{10.9, -1.2}

// tenth is 0.1 read at run time: a test pinning exact bits takes an operand from it, since the
// compiler folds literal arithmetic into a constant without fusing.
var tenth = 0.1

func TestFractionalHex_Constructor(t *testing.T) {
	AssertFractionalHex(t, FracPt(10.9, -1.2), FractionalHex{Q: 10.9, R: -1.2})
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
	AssertFractionalHex(t, testFracHex.Add(FracPt(0.1, 1.2)), FracPt(11, 0))
}

func TestFractionalHex_Subtract(t *testing.T) {
	AssertFractionalHex(t, testFracHex.Subtract(FracPt(0.9, -0.2)), FracPt(10, -1))
}

func TestFractionalHex_Multiply(t *testing.T) {
	t.Run("positive", func(t *testing.T) {
		AssertFractionalHex(t, testFracHex.Multiply(2), FracPt(21.8, -2.4))
	})
	t.Run("zero", func(t *testing.T) {
		AssertFractionalHex(t, testFracHex.Multiply(0), FracPt(0, 0))
	})
	t.Run("negative", func(t *testing.T) {
		AssertFractionalHex(t, testFracHex.Multiply(-1), FracPt(-10.9, 1.2))
	})
}

func TestFractionalHex_Lerp(t *testing.T) {
	t.Run("interpolates both coordinates", func(t *testing.T) {
		AssertFractionalHex(t, testFracHex.Lerp(FracPt(12, 0), 0.1), FracPt(11.01, -1.08))
	})
	t.Run("rounds the product before adding", func(t *testing.T) {
		lerp := FracPt(tenth, tenth).Lerp(FracPt(0.2, 0.2), 0.1)
		assert.Equal(t, lerp.Q, 0.11000000000000001)
		assert.Equal(t, lerp.R, 0.11000000000000001)
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

func TestFractionalHex_Round(t *testing.T) {
	t.Run("nearest hex", func(t *testing.T) {
		AssertHex(t, FracPt(10.9, 16.2).Round(), Pt(11, 16))
		AssertHex(t, FracPt(10.5001, 16.4999).Round(), Pt(11, 16))
	})
	t.Run("keeps the cube constraint near a corner", func(t *testing.T) {
		AssertHex(t, FracPt(10.50000001, 16.5000001).Round(), Pt(10, 17))
		AssertHex(t, FracPt(10.500001, 16.500000001).Round(), Pt(11, 16))
	})
}

func TestFractionalHex_Point(t *testing.T) {
	assert.Equal(t, testFracHex.Point(), geom.Pt(10.9, -1.2))
}

func TestFractionalHex_String(t *testing.T) {
	assert.Equal(t, testFracHex.String(), "(10.90,-1.20)")
}

func TestFractionalHex_JSON(t *testing.T) {
	data, err := json.Marshal(testFracHex)
	assert.NoError(t, err)
	assert.Equal(t, string(data), `{"q":10.9,"r":-1.2}`)

	var decoded FractionalHex
	assert.NoError(t, json.Unmarshal(data, &decoded))
	AssertFractionalHex(t, decoded, testFracHex)
}
