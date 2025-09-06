package hex

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gravitton/assert"
	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/floats"
)

func TestHexToPixel(t *testing.T) {
	// flat-top layout
	// - width             = 2 * size.Width
	// - height            = sqrt(3) * size.Height
	// - horizontalSpacing = 3/2 * size.Width (= 3/4 * width)
	// - verticalSpacing   = sqrt(3) * size.Height (= height)
	// x = q * horizontalSpacing (= q * 3/2 * size.Width)
	// y = q * verticalSpacing/2 + r * verticalSpacing (= q * 1/2 * sqrt(3) * size.Height + r * sqrt(3) * size.Height)
	layoutFlatTop := LayoutFlatTop(geom.S(10.0, 10.0), geom.P(0.0, 0.0))
	layoutFlatTopWithOrigin := LayoutFlatTop(geom.S(10.0, 10.0), geom.P(-100.0, 50.0))
	layoutFlatTopWithOriginSkew := LayoutFlatTop(geom.S(10.0, 8.0), geom.P(-100.0, 50.0))

	// pointy-top layout
	// - width             = sqrt(3) * size.Width
	// - height            = 2 * size.Height
	// - horizontalSpacing = sqrt(3) * size.Width (= width)
	// - verticalSpacing   = 3/2 * size.Height (= 3/4 * height)
	// x = q * horizontalSpacing + r * horizontalSpacing/2  (= q * sqrt(3) * size.Width + r * 1/2 * sqrt(3) * size.Height)
	// y = r * verticalSpacing (= r * 3/2 * size.Height)
	layoutPointyTop := LayoutPointyTop(geom.S(10.0, 10.0), geom.P(0.0, 0.0))
	layoutPointyTopWithOrigin := LayoutPointyTop(geom.S(10.0, 10.0), geom.P(-100.0, 50.0))
	layoutPointyTopWithOriginSkew := LayoutPointyTop(geom.S(10.0, 8.0), geom.P(-100.0, 50.0))

	tests := []struct {
		layout   Layout
		hex      Hex
		expected floats.Point
	}{
		{
			layout:   layoutFlatTop,
			hex:      H(0, 0),
			expected: geom.P(0.0, 0.0),
		},
		{
			layout:   layoutFlatTop,
			hex:      H(1, 1),
			expected: geom.P(15.0, 25.98076211353316),
		},
		{
			layout:   layoutFlatTop,
			hex:      H(-2, 3),
			expected: geom.P(-30.0, 34.64101615137755),
		},
		{
			layout:   layoutFlatTop,
			hex:      H(50, 89),
			expected: geom.P(750.0, 1974.53792062852),
		},

		{
			layout:   layoutFlatTopWithOrigin,
			hex:      H(0, 0),
			expected: geom.P(-100.0, 50.0),
		},
		{
			layout:   layoutFlatTopWithOrigin,
			hex:      H(1, 1),
			expected: geom.P(-85.0, 75.98076211353316),
		},
		{
			layout:   layoutFlatTopWithOrigin,
			hex:      H(-2, 3),
			expected: geom.P(-130.0, 84.64101615137756),
		},
		{
			layout:   layoutFlatTopWithOrigin,
			hex:      H(50, 89),
			expected: geom.P(650.0, 2024.53792062852),
		},

		{
			layout:   layoutFlatTopWithOriginSkew,
			hex:      H(0, 0),
			expected: geom.P(-100.0, 50.0),
		},
		{

			layout:   layoutFlatTopWithOriginSkew,
			hex:      H(1, 1),
			expected: geom.P(-85.0, 70.78460969082653),
		},
		{
			layout:   layoutFlatTopWithOriginSkew,
			hex:      H(-2, 3),
			expected: geom.P(-130.0, 77.71281292110204),
		},
		{
			layout:   layoutFlatTopWithOriginSkew,
			hex:      H(50, 89),
			expected: geom.P(650.0, 1629.630336502816),
		},

		{
			layout:   layoutPointyTop,
			hex:      H(0, 0),
			expected: geom.P(0.0, 0.0),
		},
		{
			layout:   layoutPointyTop,
			hex:      H(1, 1),
			expected: geom.P(25.98076211353316, 15.0),
		},
		{
			layout:   layoutPointyTop,
			hex:      H(-2, 3),
			expected: geom.P(-8.660254037844384, 45.0),
		},
		{
			layout:   layoutPointyTop,
			hex:      H(50, 89),
			expected: geom.P(1636.7880131525887, 1335.0),
		},

		{
			layout:   layoutPointyTopWithOrigin,
			hex:      H(0, 0),
			expected: geom.P(-100.0, 50.0),
		},
		{
			layout:   layoutPointyTopWithOrigin,
			hex:      H(1, 1),
			expected: geom.P(-74.01923788646684, 65.0),
		},
		{
			layout:   layoutPointyTopWithOrigin,
			hex:      H(-2, 3),
			expected: geom.P(-108.66025403784438, 95.0),
		},
		{
			layout:   layoutPointyTopWithOrigin,
			hex:      H(50, 89),
			expected: geom.P(1536.7880131525887, 1385.0),
		},

		{
			layout:   layoutPointyTopWithOriginSkew,
			hex:      H(0, 0),
			expected: geom.P(-100, 50.0),
		},
		{
			layout:   layoutPointyTopWithOriginSkew,
			hex:      H(1, 1),
			expected: geom.P(-74.01923788646684, 62.0),
		},
		{
			layout:   layoutPointyTopWithOriginSkew,
			hex:      H(-2, 3),
			expected: geom.P(-108.66025403784438, 86.0),
		},
		{
			layout:   layoutPointyTopWithOriginSkew,
			hex:      H(50, 89),
			expected: geom.P(1536.7880131525887, 1118.0),
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
