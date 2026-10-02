package hex

import (
	"fmt"
	"math"

	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/floats"
)

// FractionalHex represents a hex with floating-point axial coordinates.
// Useful for interpolation and conversions from pixel space before rounding.
type FractionalHex struct {
	Q float64 `json:"q"`
	R float64 `json:"r"`
}

// FracPt is shorthand for FractionalHex{q, r}.
func FracPt(q, r float64) FractionalHex {
	return FractionalHex{q, r}
}

// S returns the implied s coordinate (-q - r).
func (h FractionalHex) S() float64 {
	return -h.Q - h.R
}

// QR returns the (q, r) coordinates.
func (h FractionalHex) QR() (float64, float64) {
	return h.Q, h.R
}

// QRS returns the (q, r, s) coordinates where s is implied.
func (h FractionalHex) QRS() (float64, float64, float64) {
	return h.Q, h.R, h.S()
}

// Length returns the distance from the origin (0,0) in hex steps, without rounding to a hex.
func (h FractionalHex) Length() float64 {
	return (math.Abs(h.Q) + math.Abs(h.R) + math.Abs(h.S())) / 2
}

// Add returns a new FractionalHex that is the vector sum of h and hex.
func (h FractionalHex) Add(hex FractionalHex) FractionalHex {
	return FractionalHex{h.Q + hex.Q, h.R + hex.R}
}

// Subtract creates a new FractionalHex that is the vector difference h - hex.
func (h FractionalHex) Subtract(hex FractionalHex) FractionalHex {
	return FractionalHex{h.Q - hex.Q, h.R - hex.R}
}

// Multiply creates a new FractionalHex scaled by the given factor.
func (h FractionalHex) Multiply(factor float64) FractionalHex {
	return FractionalHex{h.Q * factor, h.R * factor}
}

// Lerp creates a new FractionalHex in linear interpolation towards given hex.
func (h FractionalHex) Lerp(hex FractionalHex, t float64) FractionalHex {
	return FractionalHex{geom.Lerp(h.Q, hex.Q, t), geom.Lerp(h.R, hex.R, t)}
}

// DistanceTo returns the hex distance between h and the given hex, without rounding to a hex.
func (h FractionalHex) DistanceTo(hex FractionalHex) float64 {
	return h.Subtract(hex).Length()
}

// Equal reports whether h and hex are the same hex within geom.Epsilon for float64.
func (h FractionalHex) Equal(hex FractionalHex) bool {
	return geom.Equal(h.Q, hex.Q) && geom.Equal(h.R, hex.R)
}

// IsZero reports whether the hex is at the origin (0, 0), within the tolerance Equal applies.
func (h FractionalHex) IsZero() bool {
	return h.Equal(FractionalHex{})
}

// Round converts a FractionalHex to the nearest Hex while preserving q+r+s=0.
// It is the [Hex] conversion of a fractional hex, the counterpart of [Hex.Float].
// It converts through [geom.Cast], so it panics for a NaN or infinite coordinate, which has
// no hex, and a finite coordinate beyond the range of int gives a platform-dependent hex.
func (h FractionalHex) Round() Hex {
	q := math.Round(h.Q)
	r := math.Round(h.R)
	s := math.Round(h.S())

	qDiff := math.Abs(q - h.Q)
	rDiff := math.Abs(r - h.R)
	sDiff := math.Abs(s - h.S())

	if qDiff > rDiff && qDiff > sDiff {
		q = -r - s
	} else if rDiff > sDiff {
		r = -q - s
	}

	return Hex{geom.Cast[int](q), geom.Cast[int](r)}
}

// Point returns the axial (q,r) as a floats.Point.
func (h FractionalHex) Point() floats.Point {
	return geom.Pt(h.Q, h.R)
}

// String returns a compact representation of the fractional hex as (q,r) with 2 decimals.
func (h FractionalHex) String() string {
	return fmt.Sprintf("(%.2f,%.2f)", h.Q, h.R)
}
