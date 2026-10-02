package hextest

import (
	"slices"

	"github.com/gravitton/assert"
	geom "github.com/gravitton/geometry"
	hex "github.com/gravitton/hexagon"
)

// AssertHex asserts that actual equals expected. Axial coordinates are integers,
// so the comparison is exact.
func AssertHex(t assert.Testing, actual, expected hex.Hex, messages ...string) bool {
	t.Helper()

	ok := true

	if !assert.Equal(t, actual.Q, expected.Q, prefixed(messages, "Q: ")...) {
		ok = false
	}
	if !assert.Equal(t, actual.R, expected.R, prefixed(messages, "R: ")...) {
		ok = false
	}

	return ok
}

// AssertFractionalHex asserts that actual equals expected within [geom.EpsilonRelative],
// so the tolerance holds at any magnitude.
func AssertFractionalHex(t assert.Testing, actual, expected hex.FractionalHex, messages ...string) bool {
	t.Helper()

	ok := true

	if !geom.AssertNumber(t, actual.Q, expected.Q, prefixed(messages, "Q: ")...) {
		ok = false
	}
	if !geom.AssertNumber(t, actual.R, expected.R, prefixed(messages, "R: ")...) {
		ok = false
	}

	return ok
}

// prefixed appends the field prefix to the caller's messages.
func prefixed(messages []string, prefix string) []string {
	return slices.Concat(messages, []string{prefix})
}
