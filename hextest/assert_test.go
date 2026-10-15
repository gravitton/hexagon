package hextest

import (
	"fmt"
	"strings"
	"testing"

	hex "github.com/gravitton/hexagon"
)

func TestAssertHex(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertHex, hex.Pt(1, 2), hex.Pt(1, 2), nil)
	})
	t.Run("one coordinate differs", func(t *testing.T) {
		assertHelper(t, AssertHex, hex.Pt(1, 2), hex.Pt(9, 2), nil, "Q: ")
		assertHelper(t, AssertHex, hex.Pt(1, 2), hex.Pt(1, 9), nil, "R: ")
	})
	t.Run("reports every coordinate that differs", func(t *testing.T) {
		assertHelper(t, AssertHex, hex.Pt(1, 2), hex.Pt(9, 9), nil, "Q: ", "R: ")
	})
	t.Run("the messages come before the field", func(t *testing.T) {
		assertHelper(t, AssertHex, hex.Pt(1, 2), hex.Pt(9, 2), []string{"after ", "lerp: "}, "after lerp: Q: ")
	})
}

func TestAssertFractionalHex(t *testing.T) {
	t.Run("equal", func(t *testing.T) {
		assertHelper(t, AssertFractionalHex, hex.FracPt(1, 2), hex.FracPt(1, 2), nil)
	})
	t.Run("one coordinate differs", func(t *testing.T) {
		assertHelper(t, AssertFractionalHex, hex.FracPt(1, 2), hex.FracPt(9, 2), nil, "Q: ")
		assertHelper(t, AssertFractionalHex, hex.FracPt(1, 2), hex.FracPt(1, 9), nil, "R: ")
	})
	t.Run("reports every coordinate that differs", func(t *testing.T) {
		assertHelper(t, AssertFractionalHex, hex.FracPt(1, 2), hex.FracPt(9, 9), nil, "Q: ", "R: ")
	})
	t.Run("the messages come before the field", func(t *testing.T) {
		assertHelper(t, AssertFractionalHex, hex.FracPt(1, 2), hex.FracPt(1, 9), []string{"after ", "lerp: "}, "after lerp: R: ")
	})
	t.Run("the tolerance is relative to the magnitude", func(t *testing.T) {
		assertHelper(t, AssertFractionalHex, hex.FracPt(1, 2), hex.FracPt(1.0000001, 2), nil)
		assertHelper(t, AssertFractionalHex, hex.FracPt(1e6, 2), hex.FracPt(1e6+0.5, 2), nil)
		assertHelper(t, AssertFractionalHex, hex.FracPt(1e6, 2), hex.FracPt(1e6+2, 2), nil, "Q: ")
	})
}

// assertHelper runs one of the Assert helpers with the messages against a fresh recorder and
// checks it returned whether the values are equal and reported one failure for each of the
// reports, starting with it, in order. Every helper takes (actual, expected) of the same type, so
// one signature covers them all.
func assertHelper[V any](t *testing.T, assertion func(Testing, V, V, ...string) bool, actual, expected V, messages []string, reports ...string) {
	t.Helper()

	recorder := &logger{}
	if assertion(recorder, actual, expected, messages...) != (len(reports) == 0) {
		t.Errorf("assert(%#v, %#v) should return %#v: %q", actual, expected, len(reports) == 0, recorder.errors)
	}
	if len(recorder.errors) != len(reports) {
		t.Errorf("assert(%#v, %#v) should report %d failures: %q", actual, expected, len(reports), recorder.errors)

		return
	}
	for i, report := range reports {
		if !strings.HasPrefix(recorder.errors[i], report) {
			t.Errorf("assert(%#v, %#v) should report %q: %q", actual, expected, report, recorder.errors[i])
		}
	}
}

type logger struct {
	errors []string
}

func (l *logger) Helper() {
}

func (l *logger) Errorf(format string, args ...any) {
	l.errors = append(l.errors, fmt.Sprintf(format, args...))
}
