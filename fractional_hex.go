package hex

import (
	"fmt"
	"math"

	geom "github.com/gravitton/geometry"
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

// ParseFractionalHex parses a fractional hex in the form "(q,r)", the form String prints, each
// coordinate a float [geom.Parse] accepts. String keeps two decimals, so a parsed value
// reprints the same string rather than restoring every bit. NaN and the infinities parse, as
// String prints them, to a fractional hex [FractionalHex.Round] panics for.
func ParseFractionalHex(s string) (FractionalHex, error) {
	q, r, err := parseCoordinates[float64](s, "fractional hex")
	if err != nil {
		return FractionalHex{}, err
	}

	return FractionalHex{q, r}, nil
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

// Subtract returns a new FractionalHex that is the vector difference h - hex.
func (h FractionalHex) Subtract(hex FractionalHex) FractionalHex {
	return FractionalHex{h.Q - hex.Q, h.R - hex.R}
}

// Multiply returns a new FractionalHex scaled by the given factor. Each product is rounded
// before it is returned, so a sum it feeds is never fused into a multiply-add.
func (h FractionalHex) Multiply(factor float64) FractionalHex {
	return FractionalHex{float64(h.Q * factor), float64(h.R * factor)}
}

// Lerp returns a new FractionalHex in linear interpolation towards given hex.
func (h FractionalHex) Lerp(hex FractionalHex, t float64) FractionalHex {
	return FractionalHex{geom.Lerp(h.Q, hex.Q, t), geom.Lerp(h.R, hex.R, t)}
}

// Turn returns the fractional hex rotated by steps×60° around the origin, in the same sense as
// [Hex.Turn]: clockwise as drawn on a screen with Y pointing down. Negative steps rotate the
// other way. A sixth-turn permutes and negates the cube coordinates, so the result carries the
// one rounding of s at most, however many steps it takes.
func (h FractionalHex) Turn(steps int) FractionalHex {
	switch geom.Mod(steps, 6) {
	case 1:
		return FractionalHex{-h.R, -h.S()}
	case 2:
		return FractionalHex{h.S(), h.Q}
	case 3:
		return FractionalHex{-h.Q, -h.R}
	case 4:
		return FractionalHex{h.R, h.S()}
	case 5:
		return FractionalHex{-h.S(), -h.Q}
	default:
		return h
	}
}

// TurnAround returns the fractional hex rotated by steps×60° around center, in the same sense
// as [Hex.Turn]: clockwise as drawn on a screen with Y pointing down. Negative steps rotate the
// other way.
func (h FractionalHex) TurnAround(center FractionalHex, steps int) FractionalHex {
	return center.Add(h.Subtract(center).Turn(steps))
}

// ReflectQ returns the fractional hex reflected across the q-axis (q unchanged, r and s swapped).
func (h FractionalHex) ReflectQ() FractionalHex {
	return FractionalHex{h.Q, h.S()}
}

// ReflectR returns the fractional hex reflected across the r-axis (r unchanged, q and s swapped).
func (h FractionalHex) ReflectR() FractionalHex {
	return FractionalHex{h.S(), h.R}
}

// ReflectS returns the fractional hex reflected across the s-axis (s unchanged, q and r swapped).
func (h FractionalHex) ReflectS() FractionalHex {
	return FractionalHex{h.R, h.Q}
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

// Point returns the axial (q,r) as a geom.Point[float64].
func (h FractionalHex) Point() geom.Point[float64] {
	return geom.Pt(h.Q, h.R)
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

// String returns a compact representation of the fractional hex as (q,r), each coordinate
// formatted by [geom.String] with two decimals.
func (h FractionalHex) String() string {
	return fmt.Sprintf("(%s,%s)", geom.String(h.Q), geom.String(h.R))
}
