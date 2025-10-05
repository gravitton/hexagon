package hex

import (
	"testing"

	"github.com/gravitton/assert"
	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/ints"
)

func TestCoordinateConversion(t *testing.T) {
	tests := []struct {
		hex          Hex
		offsetOddR   ints.Point
		offsetEvenR  ints.Point
		offsetOddQ   ints.Point
		offsetEvenQ  ints.Point
		doubleWidth  ints.Point
		doubleHeight ints.Point
	}{
		{
			hex:          H(0, 0),
			offsetOddR:   geom.Pt(0, 0),
			offsetEvenR:  geom.Pt(0, 0),
			offsetOddQ:   geom.Pt(0, 0),
			offsetEvenQ:  geom.Pt(0, 0),
			doubleWidth:  geom.Pt(0, 0),
			doubleHeight: geom.Pt(0, 0),
		},
		{
			hex:          H(1, 0),
			offsetOddR:   geom.Pt(1, 0),
			offsetEvenR:  geom.Pt(1, 0),
			offsetOddQ:   geom.Pt(1, 0),
			offsetEvenQ:  geom.Pt(1, 1),
			doubleWidth:  geom.Pt(2, 0),
			doubleHeight: geom.Pt(1, 1),
		},
		{
			hex:          H(1, -1),
			offsetOddR:   geom.Pt(0, -1),
			offsetEvenR:  geom.Pt(1, -1),
			offsetOddQ:   geom.Pt(1, -1),
			offsetEvenQ:  geom.Pt(1, 0),
			doubleWidth:  geom.Pt(1, -1),
			doubleHeight: geom.Pt(1, -1),
		},
		{
			hex:          H(0, -1),
			offsetOddR:   geom.Pt(-1, -1),
			offsetEvenR:  geom.Pt(0, -1),
			offsetOddQ:   geom.Pt(0, -1),
			offsetEvenQ:  geom.Pt(0, -1),
			doubleWidth:  geom.Pt(-1, -1),
			doubleHeight: geom.Pt(0, -2),
		},
		{
			hex:          H(-1, 0),
			offsetOddR:   geom.Pt(-1, 0),
			offsetEvenR:  geom.Pt(-1, 0),
			offsetOddQ:   geom.Pt(-1, -1),
			offsetEvenQ:  geom.Pt(-1, 0),
			doubleWidth:  geom.Pt(-2, 0),
			doubleHeight: geom.Pt(-1, -1),
		},
		{
			hex:          H(-1, 1),
			offsetOddR:   geom.Pt(-1, 1),
			offsetEvenR:  geom.Pt(-0, 1),
			offsetOddQ:   geom.Pt(-1, 0),
			offsetEvenQ:  geom.Pt(-1, 1),
			doubleWidth:  geom.Pt(-1, 1),
			doubleHeight: geom.Pt(-1, 1),
		},
		{
			hex:          H(0, 1),
			offsetOddR:   geom.Pt(0, 1),
			offsetEvenR:  geom.Pt(1, 1),
			offsetOddQ:   geom.Pt(0, 1),
			offsetEvenQ:  geom.Pt(0, 1),
			doubleWidth:  geom.Pt(1, 1),
			doubleHeight: geom.Pt(0, 2),
		},
	}

	for _, test := range tests {
		hexPoint := geom.Pt(test.hex.Q, test.hex.R)

		t.Run(test.hex.String(), func(t *testing.T) {
			assert.Equal(t, ToAxial(test.hex), hexPoint)
			assert.Equal(t, ToOffsetOddR(test.hex), test.offsetOddR)
			assert.Equal(t, ToOffsetEvenR(test.hex), test.offsetEvenR)
			assert.Equal(t, ToOffsetOddQ(test.hex), test.offsetOddQ)
			assert.Equal(t, ToOffsetEvenQ(test.hex), test.offsetEvenQ)
			assert.Equal(t, ToDoubleWidth(test.hex), test.doubleWidth)
			assert.Equal(t, ToDoubleHeight(test.hex), test.doubleHeight)

			assert.Equal(t, FromAxial(hexPoint), test.hex)
			assert.Equal(t, FromOffsetOddR(test.offsetOddR), test.hex)
			assert.Equal(t, FromOffsetEvenR(test.offsetEvenR), test.hex)
			assert.Equal(t, FromOffsetOddQ(test.offsetOddQ), test.hex)
			assert.Equal(t, FromOffsetEvenQ(test.offsetEvenQ), test.hex)
			assert.Equal(t, FromDoubleWidth(test.doubleWidth), test.hex)
			assert.Equal(t, FromDoubleHeight(test.doubleHeight), test.hex)
			assert.Equal(t, FromPoint(hexPoint), test.hex)

			assert.Equal(t, To(test.hex, Axial), hexPoint)
			assert.Equal(t, To(test.hex, OffsetOddR), test.offsetOddR)
			assert.Equal(t, To(test.hex, OffsetEvenR), test.offsetEvenR)
			assert.Equal(t, To(test.hex, OffsetOddQ), test.offsetOddQ)
			assert.Equal(t, To(test.hex, OffsetEvenQ), test.offsetEvenQ)
			assert.Equal(t, To(test.hex, DoubleWidth), test.doubleWidth)
			assert.Equal(t, To(test.hex, DoubleHeight), test.doubleHeight)

			assert.Equal(t, From(hexPoint, Axial), test.hex)
			assert.Equal(t, From(test.offsetOddR, OffsetOddR), test.hex)
			assert.Equal(t, From(test.offsetEvenR, OffsetEvenR), test.hex)
			assert.Equal(t, From(test.offsetOddQ, OffsetOddQ), test.hex)
			assert.Equal(t, From(test.offsetEvenQ, OffsetEvenQ), test.hex)
			assert.Equal(t, From(test.doubleWidth, DoubleWidth), test.hex)
			assert.Equal(t, From(test.doubleHeight, DoubleHeight), test.hex)

			assert.Equal(t, test.hex.To(Axial), hexPoint)
			assert.Equal(t, test.hex.To(OffsetOddR), test.offsetOddR)
			assert.Equal(t, test.hex.To(OffsetEvenR), test.offsetEvenR)
			assert.Equal(t, test.hex.To(OffsetOddQ), test.offsetOddQ)
			assert.Equal(t, test.hex.To(OffsetEvenQ), test.offsetEvenQ)
			assert.Equal(t, test.hex.To(DoubleWidth), test.doubleWidth)
			assert.Equal(t, test.hex.To(DoubleHeight), test.doubleHeight)
			assert.Equal(t, test.hex.ToPoint(), hexPoint)
		})
	}
}
