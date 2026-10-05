// Package hextest provides assertions for tests of code built on package hex, comparing a
// value field by field with the tolerance of its type and naming the field that differs.
package hextest

import (
	"slices"

	"github.com/gravitton/assert"
	"github.com/gravitton/geometry/geomtest"
	hex "github.com/gravitton/hexagon"
)

// Testing is what the assertions report a failure through, the part of testing.TB they call, so a
// *testing.T, a *testing.B or a recorder of the caller's own is passed as it is. It is declared
// here rather than taken from the assertion library underneath, which no signature names.
type Testing interface {
	Helper()
	Errorf(format string, args ...any)
}

// AssertHex asserts that actual equals expected. Axial coordinates are integers,
// so the comparison is exact.
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

// AssertFractionalHex asserts that actual equals expected within the relative tolerance
// [geomtest.AssertNumber] applies, so the tolerance holds at any magnitude.
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
