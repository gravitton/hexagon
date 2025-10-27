package hex

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gravitton/assert"
	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/floats"
)

var (
	// flat-top layout
	// - width             = 2 * size.Width
	// - height            = sqrt(3) * size.Height
	// - horizontalSpacing = 3/2 * size.Width (= 3/4 * width)
	// - verticalSpacing   = sqrt(3) * size.Height (= height)
	// x = q * horizontalSpacing (= q * 3/2 * size.Width)
	// y = q * verticalSpacing/2 + r * verticalSpacing (= q * 1/2 * sqrt(3) * size.Height + r * sqrt(3) * size.Height)
	layoutFlatTop               = LayoutFlatTop(geom.Sz(10.0, 10.0), geom.Pt(0.0, 0.0))
	layoutFlatTopWithOrigin     = LayoutFlatTop(geom.Sz(10.0, 10.0), geom.Pt(-100.0, 50.0))
	layoutFlatTopWithOriginSkew = LayoutFlatTop(geom.Sz(10.0, 8.0), geom.Pt(-100.0, 50.0))

	// pointy-top layout
	// - width             = sqrt(3) * size.Width
	// - height            = 2 * size.Height
	// - horizontalSpacing = sqrt(3) * size.Width (= width)
	// - verticalSpacing   = 3/2 * size.Height (= 3/4 * height)
	// x = q * horizontalSpacing + r * horizontalSpacing/2  (= q * sqrt(3) * size.Width + r * 1/2 * sqrt(3) * size.Height)
	// y = r * verticalSpacing (= r * 3/2 * size.Height)
	layoutPointyTop               = LayoutPointyTop(geom.Sz(10.0, 10.0), geom.Pt(0.0, 0.0))
	layoutPointyTopWithOrigin     = LayoutPointyTop(geom.Sz(10.0, 10.0), geom.Pt(-100.0, 50.0))
	layoutPointyTopWithOriginSkew = LayoutPointyTop(geom.Sz(10.0, 8.0), geom.Pt(-100.0, 50.0))
)

func TestLayout_Point(t *testing.T) {
	tests := []struct {
		layout   Layout
		hex      Hex
		expected floats.Point
	}{
		{
			layout:   layoutFlatTop,
			hex:      H(0, 0),
			expected: geom.Pt(0.0, 0.0),
		},
		{
			layout:   layoutFlatTop,
			hex:      H(1, 1),
			expected: geom.Pt(15.0, 25.98076211353316),
		},
		{
			layout:   layoutFlatTop,
			hex:      H(-2, 3),
			expected: geom.Pt(-30.0, 34.64101615137755),
		},
		{
			layout:   layoutFlatTop,
			hex:      H(50, 89),
			expected: geom.Pt(750.0, 1974.53792062852),
		},

		{
			layout:   layoutFlatTopWithOrigin,
			hex:      H(0, 0),
			expected: geom.Pt(-100.0, 50.0),
		},
		{
			layout:   layoutFlatTopWithOrigin,
			hex:      H(1, 1),
			expected: geom.Pt(-85.0, 75.98076211353316),
		},
		{
			layout:   layoutFlatTopWithOrigin,
			hex:      H(-2, 3),
			expected: geom.Pt(-130.0, 84.64101615137756),
		},
		{
			layout:   layoutFlatTopWithOrigin,
			hex:      H(50, 89),
			expected: geom.Pt(650.0, 2024.53792062852),
		},

		{
			layout:   layoutFlatTopWithOriginSkew,
			hex:      H(0, 0),
			expected: geom.Pt(-100.0, 50.0),
		},
		{

			layout:   layoutFlatTopWithOriginSkew,
			hex:      H(1, 1),
			expected: geom.Pt(-85.0, 70.78460969082653),
		},
		{
			layout:   layoutFlatTopWithOriginSkew,
			hex:      H(-2, 3),
			expected: geom.Pt(-130.0, 77.71281292110204),
		},
		{
			layout:   layoutFlatTopWithOriginSkew,
			hex:      H(50, 89),
			expected: geom.Pt(650.0, 1629.630336502816),
		},

		{
			layout:   layoutPointyTop,
			hex:      H(0, 0),
			expected: geom.Pt(0.0, 0.0),
		},
		{
			layout:   layoutPointyTop,
			hex:      H(1, 1),
			expected: geom.Pt(25.98076211353316, 15.0),
		},
		{
			layout:   layoutPointyTop,
			hex:      H(-2, 3),
			expected: geom.Pt(-8.660254037844384, 45.0),
		},
		{
			layout:   layoutPointyTop,
			hex:      H(50, 89),
			expected: geom.Pt(1636.7880131525887, 1335.0),
		},

		{
			layout:   layoutPointyTopWithOrigin,
			hex:      H(0, 0),
			expected: geom.Pt(-100.0, 50.0),
		},
		{
			layout:   layoutPointyTopWithOrigin,
			hex:      H(1, 1),
			expected: geom.Pt(-74.01923788646684, 65.0),
		},
		{
			layout:   layoutPointyTopWithOrigin,
			hex:      H(-2, 3),
			expected: geom.Pt(-108.66025403784438, 95.0),
		},
		{
			layout:   layoutPointyTopWithOrigin,
			hex:      H(50, 89),
			expected: geom.Pt(1536.7880131525887, 1385.0),
		},

		{
			layout:   layoutPointyTopWithOriginSkew,
			hex:      H(0, 0),
			expected: geom.Pt(-100, 50.0),
		},
		{
			layout:   layoutPointyTopWithOriginSkew,
			hex:      H(1, 1),
			expected: geom.Pt(-74.01923788646684, 62.0),
		},
		{
			layout:   layoutPointyTopWithOriginSkew,
			hex:      H(-2, 3),
			expected: geom.Pt(-108.66025403784438, 86.0),
		},
		{
			layout:   layoutPointyTopWithOriginSkew,
			hex:      H(50, 89),
			expected: geom.Pt(1536.7880131525887, 1118.0),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s_H%s", testName(test.layout), test.hex), func(t *testing.T) {
			actual := test.layout.ToPoint(test.hex)

			assert.EqualDelta(t, actual.X, test.expected.X, geom.Delta)
			assert.EqualDelta(t, actual.Y, test.expected.Y, geom.Delta)

			actualHex := test.layout.FromPoint(test.expected).Round()

			assert.Equal(t, actualHex.Q, test.hex.Q)
			assert.Equal(t, actualHex.R, test.hex.R)
		})
	}
}

func TestLayout_Bounds(t *testing.T) {
	geom.AssertSize(t, layoutFlatTop.Bounds(), 20.0, 17.320508)
	geom.AssertSize(t, layoutFlatTopWithOrigin.Bounds(), 20.0, 17.320508)
	geom.AssertSize(t, layoutFlatTopWithOriginSkew.Bounds(), 20.0, 13.856406)
	geom.AssertSize(t, layoutPointyTop.Bounds(), 17.320508, 20.0)
	geom.AssertSize(t, layoutPointyTopWithOrigin.Bounds(), 17.320508, 20.0)
	geom.AssertSize(t, layoutPointyTopWithOriginSkew.Bounds(), 17.320508, 16.0)
}

func TestLayout_Spacing(t *testing.T) {
	geom.AssertSize(t, layoutFlatTop.Spacing(), 15.0, 17.320508)
	geom.AssertSize(t, layoutFlatTopWithOrigin.Spacing(), 15.0, 17.320508)
	geom.AssertSize(t, layoutFlatTopWithOriginSkew.Spacing(), 15.0, 13.856406)
	geom.AssertSize(t, layoutPointyTop.Spacing(), 17.320508, 15.0)
	geom.AssertSize(t, layoutPointyTopWithOrigin.Spacing(), 17.320508, 15.0)
	geom.AssertSize(t, layoutPointyTopWithOriginSkew.Spacing(), 17.320508, 12.0)
}

func TestLayout_Hexagon(t *testing.T) {
	geom.AssertRegularPolygon(t, layoutFlatTop.Hexagon(H(0, 0)), 0.0, 0.0, 10.0, 10.0, 6, 60*geom.DegToRad)
	geom.AssertRegularPolygon(t, layoutFlatTopWithOrigin.Hexagon(H(0, 0)), -100.0, 50.0, 10.0, 10.0, 6, 60*geom.DegToRad)
	geom.AssertRegularPolygon(t, layoutFlatTopWithOriginSkew.Hexagon(H(1, 0)), -85.0, 56.928203, 10.0, 8.0, 6, 60*geom.DegToRad)
	geom.AssertRegularPolygon(t, layoutPointyTop.Hexagon(H(0, 0)), 0.0, 0.0, 10.0, 10.0, 6, 90*geom.DegToRad)
	geom.AssertRegularPolygon(t, layoutPointyTopWithOrigin.Hexagon(H(0, 0)), -100.0, 50.0, 10.0, 10.0, 6, 90*geom.DegToRad)
	geom.AssertRegularPolygon(t, layoutPointyTopWithOriginSkew.Hexagon(H(0, 1)), -91.339745, 62.0, 10.0, 8.0, 6, 90*geom.DegToRad)
}

func testName(layout Layout) string {
	name := []string{}

	switch layout.orientation.orientation {
	case geom.FlatTop:
		name = append(name, "flat-top")
	case geom.PointTop:
		name = append(name, "point-top")
	}

	if !layout.Origin.IsZero() {
		name = append(name, "translated")
	}

	if layout.Size.Width != layout.Size.Height {
		name = append(name, "skewed")
	}

	return strings.Join(name, "-")
}
