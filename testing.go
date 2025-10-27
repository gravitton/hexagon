package hex

import (
	"testing"

	"github.com/gravitton/assert"
	geom "github.com/gravitton/geometry"
)

func AssertHex(t *testing.T, h Hex, q, r int, messages ...string) {
	t.Helper()

	assert.Equal(t, h.Q, q, append(messages, "H.")...)
	assert.Equal(t, h.R, r, append(messages, "R.")...)
}

func AssertFracHex(t *testing.T, h FractionalHex, q, r float64, messages ...string) {
	t.Helper()

	assert.EqualDelta(t, h.Q, q, geom.Delta, append(messages, "H.")...)
	assert.EqualDelta(t, h.R, r, geom.Delta, append(messages, "R.")...)
}
