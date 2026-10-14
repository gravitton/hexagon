package hex

import (
	"math"

	geom "github.com/gravitton/geometry"
)

// FractionalHex is a position in axial coordinates between hex centers, for interpolation and for
// rounding a position onto the grid.
type FractionalHex struct {
	Q float64 `json:"q"`
	R float64 `json:"r"`
}

// FracPt is shorthand for FractionalHex{q, r}.
func FracPt(q, r float64) FractionalHex {
	return FractionalHex{q, r}
}

// ParseFractionalHex returns the fractional hex in the form String prints, and an error for any
// other string. Each coordinate is a float geom.Parse accepts, NaN and the infinities included.
// String keeps two decimals, so a parsed value reprints the same string rather than restoring every
// bit.
func ParseFractionalHex(s string) (FractionalHex, error) {
	q, r, err := parseCoordinates[float64](s, "fractional hex")
	if err != nil {
		return FractionalHex{}, err
	}

	return FractionalHex{q, r}, nil
}

// S returns the third cube coordinate (-q-r).
func (h FractionalHex) S() float64 {
	return -h.Q - h.R
}

// QR returns Q and R.
func (h FractionalHex) QR() (float64, float64) {
	return h.Q, h.R
}

// QRS returns Q, R and S.
func (h FractionalHex) QRS() (float64, float64, float64) {
	return h.Q, h.R, h.S()
}

// Length returns the distance from the origin in steps.
func (h FractionalHex) Length() float64 {
	return (math.Abs(h.Q) + math.Abs(h.R) + math.Abs(h.S())) / 2
}

// Add returns the sum of the fractional hexes.
func (h FractionalHex) Add(hex FractionalHex) FractionalHex {
	return FractionalHex{h.Q + hex.Q, h.R + hex.R}
}

// Subtract returns the fractional hex less the given one.
func (h FractionalHex) Subtract(hex FractionalHex) FractionalHex {
	return FractionalHex{h.Q - hex.Q, h.R - hex.R}
}

// Multiply returns the fractional hex with both coordinates multiplied by the factor.
func (h FractionalHex) Multiply(factor float64) FractionalHex {
	return FractionalHex{float64(h.Q * factor), float64(h.R * factor)}
}

// Lerp returns the fractional hex interpolated toward the given one, extrapolated outside [0, 1].
func (h FractionalHex) Lerp(hex FractionalHex, t float64) FractionalHex {
	return FractionalHex{geom.Lerp(h.Q, hex.Q, t), geom.Lerp(h.R, hex.R, t)}
}

// Turn returns the fractional hex turned about the origin by steps sixth-turns of increasing angle;
// negative steps go the other way. The result is exact but for one rounding of s, however many
// steps it takes.
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

// TurnAround returns the fractional hex turned about center by steps sixth-turns of increasing
// angle; negative steps go the other way.
func (h FractionalHex) TurnAround(center FractionalHex, steps int) FractionalHex {
	return center.Add(h.Subtract(center).Turn(steps))
}

// ReflectQ returns the fractional hex reflected across the q-axis, swapping r and s.
func (h FractionalHex) ReflectQ() FractionalHex {
	return FractionalHex{h.Q, h.S()}
}

// ReflectR returns the fractional hex reflected across the r-axis, swapping q and s.
func (h FractionalHex) ReflectR() FractionalHex {
	return FractionalHex{h.S(), h.R}
}

// ReflectS returns the fractional hex reflected across the s-axis, swapping q and r.
func (h FractionalHex) ReflectS() FractionalHex {
	return FractionalHex{h.R, h.Q}
}

// DistanceTo returns the distance to the given fractional hex in steps.
func (h FractionalHex) DistanceTo(hex FractionalHex) float64 {
	return h.Subtract(hex).Length()
}

// Equal reports whether the fractional hexes have equal coordinates.
func (h FractionalHex) Equal(hex FractionalHex) bool {
	return geom.Equal(h.Q, hex.Q) && geom.Equal(h.R, hex.R)
}

// IsZero reports whether the fractional hex is the origin.
func (h FractionalHex) IsZero() bool {
	return h.Equal(FractionalHex{})
}

// IsFinite reports whether both coordinates are neither NaN nor infinite, which Round needs.
func (h FractionalHex) IsFinite() bool {
	return geom.IsFinite(h.Q) && geom.IsFinite(h.R)
}

// Point returns Q and R as a point.
func (h FractionalHex) Point() geom.Point[float64] {
	return geom.Pt(h.Q, h.R)
}

// Round returns the hex the fractional hex lies in. A position on an edge or a corner goes to the
// same side on every hex, a tie of each coordinate going toward positive infinity, so rounding a
// fractional hex moved by a whole one gives the rounded hex moved by it. It panics for a NaN or
// infinite coordinate.
func (h FractionalHex) Round() Hex {
	q := geom.RoundHalfUp(h.Q)
	r := geom.RoundHalfUp(h.R)
	s := geom.RoundHalfUp(h.S())

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

// String returns the fractional hex as its coordinates, with two decimals: (q,r).
func (h FractionalHex) String() string {
	return h.Point().String()
}
