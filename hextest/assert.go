// Package hextest provides test assertions for code built on package hex. Each compares the
// values field by field, exactly for a Hex and within geom.EpsilonRelative for a FractionalHex,
// reports every field that differs by its name, and returns whether the values are equal.
//
// The tolerance for a FractionalHex grows with its magnitude, where FractionalHex.Equal keeps
// geom.Epsilon, so far from the origin two positions this accepts as equal can round to
// different hexes. A test of what a position rounds to asserts the Hex Round returns.
package hextest

import (
	"slices"

	"github.com/gravitton/assert"
	"github.com/gravitton/geometry/geomtest"
	hex "github.com/gravitton/hexagon"
)

// Testing is the part of testing.TB the assertions report through, so a *testing.T, a *testing.B or
// a recorder of the caller's own can be passed.
type Testing interface {
	Helper()
	Errorf(format string, args ...any)
}

// AssertHex asserts that actual equals expected.
func AssertHex(t Testing, actual, expected hex.Hex, messages ...string) bool {
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

// AssertFractionalHex asserts that actual equals expected.
func AssertFractionalHex(t Testing, actual, expected hex.FractionalHex, messages ...string) bool {
	t.Helper()

	ok := true

	if !geomtest.AssertNumber(t, actual.Q, expected.Q, prefixed(messages, "Q: ")...) {
		ok = false
	}
	if !geomtest.AssertNumber(t, actual.R, expected.R, prefixed(messages, "R: ")...) {
		ok = false
	}

	return ok
}

// prefixed appends the field prefix to the caller's messages.
func prefixed(messages []string, prefix string) []string {
	return slices.Concat(messages, []string{prefix})
}
