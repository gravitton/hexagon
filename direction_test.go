package hex

import (
	"testing"

	"github.com/gravitton/assert"
	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/ints"
)

func TestNeighbors(t *testing.T) {
	axialDirection := []ints.Vector{{X: 1, Y: 0}, {X: 1, Y: -1}, {X: 0, Y: -1}, {X: -1, Y: 0}, {X: -1, Y: 1}, {X: 0, Y: 1}}
	offsetOddRDirectionOddRow := []ints.Vector{{X: 1, Y: 0}, {X: 1, Y: -1}, {X: 0, Y: -1}, {X: -1, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}}
	offsetOddRDirectionEvenRow := []ints.Vector{{X: 1, Y: 0}, {X: 0, Y: -1}, {X: -1, Y: -1}, {X: -1, Y: 0}, {X: -1, Y: 1}, {X: 0, Y: 1}}
	offsetEvenRDirectionOddRow := []ints.Vector{{X: 1, Y: 0}, {X: 0, Y: -1}, {X: -1, Y: -1}, {X: -1, Y: 0}, {X: -1, Y: 1}, {X: 0, Y: 1}}
	offsetEvenRDirectionEvenRow := []ints.Vector{{X: 1, Y: 0}, {X: 1, Y: -1}, {X: 0, Y: -1}, {X: -1, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}}
	offsetOddQDirectionOddCol := []ints.Vector{{X: 1, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: -1}, {X: -1, Y: 0}, {X: -1, Y: 1}, {X: 0, Y: 1}}
	offsetOddQDirectionEvenCol := []ints.Vector{{X: 1, Y: 0}, {X: 1, Y: -1}, {X: 0, Y: -1}, {X: -1, Y: -1}, {X: -1, Y: 0}, {X: 0, Y: 1}}
	offsetEvenQDirectionOddCol := []ints.Vector{{X: 1, Y: 0}, {X: 1, Y: -1}, {X: 0, Y: -1}, {X: -1, Y: -1}, {X: -1, Y: 0}, {X: 0, Y: 1}}
	offsetEvenQDirectionEvenCol := []ints.Vector{{X: 1, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: -1}, {X: -1, Y: 0}, {X: -1, Y: 1}, {X: 0, Y: 1}}
	doubleWidthDirection := []ints.Vector{{X: 2, Y: 0}, {X: 1, Y: -1}, {X: -1, Y: -1}, {X: -2, Y: 0}, {X: -1, Y: 1}, {X: 1, Y: 1}}
	doubleHeightDirection := []ints.Vector{{X: 1, Y: 1}, {X: 1, Y: -1}, {X: 0, Y: -2}, {X: -1, Y: -1}, {X: -1, Y: 1}, {X: 0, Y: 2}}

	tests := []struct {
		index        ints.Point
		axial        []ints.Vector
		offsetOddR   []ints.Vector
		offsetEvenR  []ints.Vector
		offsetOddQ   []ints.Vector
		offsetEvenQ  []ints.Vector
		doubleWidth  []ints.Vector
		doubleHeight []ints.Vector
	}{
		{
			index:        geom.P(0, 0), // even col, even row
			axial:        axialDirection,
			offsetOddR:   offsetOddRDirectionEvenRow,
			offsetEvenR:  offsetEvenRDirectionEvenRow,
			offsetOddQ:   offsetOddQDirectionEvenCol,
			offsetEvenQ:  offsetEvenQDirectionEvenCol,
			doubleWidth:  doubleWidthDirection,
			doubleHeight: doubleHeightDirection,
		},
		{
			index:        geom.P(1, 1), // odd col, odd row
			axial:        axialDirection,
			offsetOddR:   offsetOddRDirectionOddRow,
			offsetEvenR:  offsetEvenRDirectionOddRow,
			offsetOddQ:   offsetOddQDirectionOddCol,
			offsetEvenQ:  offsetEvenQDirectionOddCol,
			doubleWidth:  doubleWidthDirection,
			doubleHeight: doubleHeightDirection,
		},
		{
			index:        geom.P(3, 2), // odd col, even row
			axial:        axialDirection,
			offsetOddR:   offsetOddRDirectionEvenRow,
			offsetEvenR:  offsetEvenRDirectionEvenRow,
			offsetOddQ:   offsetOddQDirectionOddCol,
			offsetEvenQ:  offsetEvenQDirectionOddCol,
			doubleWidth:  doubleWidthDirection,
			doubleHeight: doubleHeightDirection,
		},
		{
			index:        geom.P(-2, -3), // even col, odd row
			axial:        axialDirection,
			offsetOddR:   offsetOddRDirectionOddRow,
			offsetEvenR:  offsetEvenRDirectionOddRow,
			offsetOddQ:   offsetOddQDirectionEvenCol,
			offsetEvenQ:  offsetEvenQDirectionEvenCol,
			doubleWidth:  doubleWidthDirection,
			doubleHeight: doubleHeightDirection,
		},
	}

	for _, test := range tests {
		t.Run(test.index.String(), func(t *testing.T) {
			assert.Equal(t, NeighborOffsets(test.index, Axial), test.axial)
			assert.Equal(t, NeighborOffsets(test.index, OffsetOddR), test.offsetOddR)
			assert.Equal(t, NeighborOffsets(test.index, OffsetEvenR), test.offsetEvenR)
			assert.Equal(t, NeighborOffsets(test.index, OffsetOddQ), test.offsetOddQ)
			assert.Equal(t, NeighborOffsets(test.index, OffsetEvenQ), test.offsetEvenQ)
			assert.Equal(t, NeighborOffsets(test.index, DoubleWidth), test.doubleWidth)
			assert.Equal(t, NeighborOffsets(test.index, DoubleHeight), test.doubleHeight)
		})
	}
}
