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
			offsetOddR:   geom.P(0, 0),
			offsetEvenR:  geom.P(0, 0),
			offsetOddQ:   geom.P(0, 0),
			offsetEvenQ:  geom.P(0, 0),
			doubleWidth:  geom.P(0, 0),
			doubleHeight: geom.P(0, 0),
		},
		{
			hex:          H(1, 0),
			offsetOddR:   geom.P(1, 0),
			offsetEvenR:  geom.P(1, 0),
			offsetOddQ:   geom.P(1, 0),
			offsetEvenQ:  geom.P(1, 1),
			doubleWidth:  geom.P(2, 0),
			doubleHeight: geom.P(1, 1),
		},
		{
			hex:          H(1, -1),
			offsetOddR:   geom.P(0, -1),
			offsetEvenR:  geom.P(1, -1),
			offsetOddQ:   geom.P(1, -1),
			offsetEvenQ:  geom.P(1, 0),
			doubleWidth:  geom.P(1, -1),
			doubleHeight: geom.P(1, -1),
		},
		{
			hex:          H(0, -1),
			offsetOddR:   geom.P(-1, -1),
			offsetEvenR:  geom.P(0, -1),
			offsetOddQ:   geom.P(0, -1),
			offsetEvenQ:  geom.P(0, -1),
			doubleWidth:  geom.P(-1, -1),
			doubleHeight: geom.P(0, -2),
		},
		{
			hex:          H(-1, 0),
			offsetOddR:   geom.P(-1, 0),
			offsetEvenR:  geom.P(-1, 0),
			offsetOddQ:   geom.P(-1, -1),
			offsetEvenQ:  geom.P(-1, 0),
			doubleWidth:  geom.P(-2, 0),
			doubleHeight: geom.P(-1, -1),
		},
		{
			hex:          H(-1, 1),
			offsetOddR:   geom.P(-1, 1),
			offsetEvenR:  geom.P(-0, 1),
			offsetOddQ:   geom.P(-1, 0),
			offsetEvenQ:  geom.P(-1, 1),
			doubleWidth:  geom.P(-1, 1),
			doubleHeight: geom.P(-1, 1),
		},
		{
			hex:          H(0, 1),
			offsetOddR:   geom.P(0, 1),
			offsetEvenR:  geom.P(1, 1),
			offsetOddQ:   geom.P(0, 1),
			offsetEvenQ:  geom.P(0, 1),
			doubleWidth:  geom.P(1, 1),
			doubleHeight: geom.P(0, 2),
		},
	}

	for _, test := range tests {
		hexPoint := geom.P(test.hex.Q, test.hex.R)
		t.Run(test.hex.String(), func(t *testing.T) {
			assert.Equal(t, ToOffsetOddR(test.hex), test.offsetOddR)
			assert.Equal(t, ToOffsetEvenR(test.hex), test.offsetEvenR)
			assert.Equal(t, ToOffsetOddQ(test.hex), test.offsetOddQ)
			assert.Equal(t, ToOffsetEvenQ(test.hex), test.offsetEvenQ)
			assert.Equal(t, ToDoubleWidth(test.hex), test.doubleWidth)
			assert.Equal(t, ToDoubleHeight(test.hex), test.doubleHeight)

			assert.Equal(t, FromPoint(hexPoint), test.hex)
			assert.Equal(t, FromOffsetOddR(test.offsetOddR), test.hex)
			assert.Equal(t, FromOffsetEvenR(test.offsetEvenR), test.hex)
			assert.Equal(t, FromOffsetOddQ(test.offsetOddQ), test.hex)
			assert.Equal(t, FromOffsetEvenQ(test.offsetEvenQ), test.hex)
			assert.Equal(t, FromDoubleWidth(test.doubleWidth), test.hex)
			assert.Equal(t, FromDoubleHeight(test.doubleHeight), test.hex)

			assert.Equal(t, HexTo(test.hex, Axial), hexPoint)
			assert.Equal(t, HexTo(test.hex, OffsetOddR), test.offsetOddR)
			assert.Equal(t, HexTo(test.hex, OffsetEvenR), test.offsetEvenR)
			assert.Equal(t, HexTo(test.hex, OffsetOddQ), test.offsetOddQ)
			assert.Equal(t, HexTo(test.hex, OffsetEvenQ), test.offsetEvenQ)
			assert.Equal(t, HexTo(test.hex, DoubleWidth), test.doubleWidth)
			assert.Equal(t, HexTo(test.hex, DoubleHeight), test.doubleHeight)

			assert.Equal(t, HexFrom(hexPoint, Axial), test.hex)
			assert.Equal(t, HexFrom(test.offsetOddR, OffsetOddR), test.hex)
			assert.Equal(t, HexFrom(test.offsetEvenR, OffsetEvenR), test.hex)
			assert.Equal(t, HexFrom(test.offsetOddQ, OffsetOddQ), test.hex)
			assert.Equal(t, HexFrom(test.offsetEvenQ, OffsetEvenQ), test.hex)
			assert.Equal(t, HexFrom(test.doubleWidth, DoubleWidth), test.hex)
			assert.Equal(t, HexFrom(test.doubleHeight, DoubleHeight), test.hex)

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
