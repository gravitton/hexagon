package hextest

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gravitton/assert"
	hex "github.com/gravitton/hexagon"
)

func TestAssertHex(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertHex, hex.Pt(1, 2), hex.Pt(1, 2), true)
	})
	t.Run("one coordinate differs", func(t *testing.T) {
		assertHelper(t, AssertHex, hex.Pt(1, 2), hex.Pt(9, 2), false, "Q")
		assertHelper(t, AssertHex, hex.Pt(1, 2), hex.Pt(1, 9), false, "R")
	})
}

func TestAssertFractionalHex(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertFractionalHex, hex.FracPt(1, 2), hex.FracPt(1, 2), true)
	})
	t.Run("one coordinate differs", func(t *testing.T) {
		assertHelper(t, AssertFractionalHex, hex.FracPt(1, 2), hex.FracPt(9, 2), false, "Q")
		assertHelper(t, AssertFractionalHex, hex.FracPt(1, 2), hex.FracPt(1, 9), false, "R")
	})
	t.Run("the tolerance is relative to the magnitude", func(t *testing.T) {
		assertHelper(t, AssertFractionalHex, hex.FracPt(1, 2), hex.FracPt(1.0000001, 2), true)
		assertHelper(t, AssertFractionalHex, hex.FracPt(1e6, 2), hex.FracPt(1e6+0.5, 2), true)
		assertHelper(t, AssertFractionalHex, hex.FracPt(1e6, 2), hex.FracPt(1e6+2, 2), false)
	})
}

// assertHelper runs one of the Assert helpers against a fresh recorder and checks it
// reported the expected outcome. Every helper takes (actual, expected) of the same
// type, so one signature covers them all.
func assertHelper[V any](t *testing.T, assertion func(assert.Testing, V, V, ...string) bool, actual, expected V, result bool, labels ...string) {
	t.Helper()

	recorder := &logger{}
	if assertion(recorder, actual, expected) != result {
		t.Errorf("%sassert(%#v, %#v) should return %#v: %s", strings.Join(labels, ""), actual, expected, result, recorder.LastError)
	}
}

type logger struct {
	LastError string
}

func (m *logger) Helper() {
}

func (m *logger) Errorf(format string, args ...any) {
	m.LastError = fmt.Sprintf(format, args...)
}
