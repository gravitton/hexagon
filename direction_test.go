package hex

import (
	"testing"

	"github.com/gravitton/assert"
	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/ints"
)

var axialDirection = []ints.Vector{{X: 1, Y: 0}, {X: 1, Y: -1}, {X: 0, Y: -1}, {X: -1, Y: 0}, {X: -1, Y: 1}, {X: 0, Y: 1}}
var offsetOddRDirectionOddRow = []ints.Vector{{X: 1, Y: 0}, {X: 1, Y: -1}, {X: 0, Y: -1}, {X: -1, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}}
var offsetOddRDirectionEvenRow = []ints.Vector{{X: 1, Y: 0}, {X: 0, Y: -1}, {X: -1, Y: -1}, {X: -1, Y: 0}, {X: -1, Y: 1}, {X: 0, Y: 1}}
var offsetEvenRDirectionOddRow = []ints.Vector{{X: 1, Y: 0}, {X: 0, Y: -1}, {X: -1, Y: -1}, {X: -1, Y: 0}, {X: -1, Y: 1}, {X: 0, Y: 1}}
var offsetEvenRDirectionEvenRow = []ints.Vector{{X: 1, Y: 0}, {X: 1, Y: -1}, {X: 0, Y: -1}, {X: -1, Y: 0}, {X: 0, Y: 1}, {X: 1, Y: 1}}
var offsetOddQDirectionOddCol = []ints.Vector{{X: 1, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: -1}, {X: -1, Y: 0}, {X: -1, Y: 1}, {X: 0, Y: 1}}
var offsetOddQDirectionEvenCol = []ints.Vector{{X: 1, Y: 0}, {X: 1, Y: -1}, {X: 0, Y: -1}, {X: -1, Y: -1}, {X: -1, Y: 0}, {X: 0, Y: 1}}
var offsetEvenQDirectionOddCol = []ints.Vector{{X: 1, Y: 0}, {X: 1, Y: -1}, {X: 0, Y: -1}, {X: -1, Y: -1}, {X: -1, Y: 0}, {X: 0, Y: 1}}
var offsetEvenQDirectionEvenCol = []ints.Vector{{X: 1, Y: 1}, {X: 1, Y: 0}, {X: 0, Y: -1}, {X: -1, Y: 0}, {X: -1, Y: 1}, {X: 0, Y: 1}}
var doubleWidthDirection = []ints.Vector{{X: 2, Y: 0}, {X: 1, Y: -1}, {X: -1, Y: -1}, {X: -2, Y: 0}, {X: -1, Y: 1}, {X: 1, Y: 1}}
var doubleHeightDirection = []ints.Vector{{X: 1, Y: 1}, {X: 1, Y: -1}, {X: 0, Y: -2}, {X: -1, Y: -1}, {X: -1, Y: 1}, {X: 0, Y: 2}}

func TestDirections(t *testing.T) {
	assert.Equal(t, DirectionSMinus.NeighborOffset(), axialDirection[0])
	assert.Equal(t, DirectionQPlus.NeighborOffset(), axialDirection[1])
	assert.Equal(t, DirectionRMinus.NeighborOffset(), axialDirection[2])
	assert.Equal(t, DirectionSPlus.NeighborOffset(), axialDirection[3])
	assert.Equal(t, DirectionQMinus.NeighborOffset(), axialDirection[4])
	assert.Equal(t, DirectionRPlus.NeighborOffset(), axialDirection[5])

	// direction neighbor direction module
	assert.Equal(t, Direction(6).NeighborOffset(), axialDirection[0])
	assert.Equal(t, Direction(8).NeighborOffset(), axialDirection[2])
	assert.Equal(t, Direction(15).NeighborOffset(), axialDirection[3])

	assert.Equal(t, DirectionFlatTopSE, DirectionSMinus)
	assert.Equal(t, DirectionFlatTopNE, DirectionQPlus)
	assert.Equal(t, DirectionFlatTopN, DirectionRMinus)
	assert.Equal(t, DirectionFlatTopNW, DirectionSPlus)
	assert.Equal(t, DirectionFlatTopSW, DirectionQMinus)
	assert.Equal(t, DirectionFlatTopS, DirectionRPlus)

	assert.Equal(t, DirectionPointyTopE, DirectionSMinus)
	assert.Equal(t, DirectionPointyTopNE, DirectionQPlus)
	assert.Equal(t, DirectionPointyTopNW, DirectionRMinus)
	assert.Equal(t, DirectionPointyTopW, DirectionSPlus)
	assert.Equal(t, DirectionPointyTopSW, DirectionQMinus)
	assert.Equal(t, DirectionPointyTopSE, DirectionRPlus)
}

func TestNeighbors(t *testing.T) {
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
			assert.Equal(t, NeighborOffsetsAxial(), test.axial)
			assert.Equal(t, NeighborOffsetsOffsetOddR(test.index), test.offsetOddR)
			assert.Equal(t, NeighborOffsetsOffsetEvenR(test.index), test.offsetEvenR)
			assert.Equal(t, NeighborOffsetsOffsetOddQ(test.index), test.offsetOddQ)
			assert.Equal(t, NeighborOffsetsOffsetEvenQ(test.index), test.offsetEvenQ)
			assert.Equal(t, NeighborOffsetsDoubleWidth(), test.doubleWidth)
			assert.Equal(t, NeighborOffsetsDoubleHeight(), test.doubleHeight)

			assert.Equal(t, NeighborOffsets(test.index, Axial), test.axial)
			assert.Equal(t, NeighborOffsets(test.index, OffsetOddR), test.offsetOddR)
			assert.Equal(t, NeighborOffsets(test.index, OffsetEvenR), test.offsetEvenR)
			assert.Equal(t, NeighborOffsets(test.index, OffsetOddQ), test.offsetOddQ)
			assert.Equal(t, NeighborOffsets(test.index, OffsetEvenQ), test.offsetEvenQ)
			assert.Equal(t, NeighborOffsets(test.index, DoubleWidth), test.doubleWidth)
			assert.Equal(t, NeighborOffsets(test.index, DoubleHeight), test.doubleHeight)

			for system, offsets := range map[CoordinateSystem][]ints.Vector{
				Axial:        test.axial,
				OffsetOddR:   test.offsetOddR,
				OffsetEvenR:  test.offsetEvenR,
				OffsetOddQ:   test.offsetOddQ,
				OffsetEvenQ:  test.offsetEvenQ,
				DoubleWidth:  test.doubleWidth,
				DoubleHeight: test.doubleHeight,
			} {
				assert.Equal(t, NeighborOffset(test.index, system, DirectionSMinus), offsets[0])
				assert.Equal(t, NeighborOffset(test.index, system, DirectionQPlus), offsets[1])
				assert.Equal(t, NeighborOffset(test.index, system, DirectionRMinus), offsets[2])
				assert.Equal(t, NeighborOffset(test.index, system, DirectionSPlus), offsets[3])
				assert.Equal(t, NeighborOffset(test.index, system, DirectionQMinus), offsets[4])
				assert.Equal(t, NeighborOffset(test.index, system, DirectionRPlus), offsets[5])
			}
		})
	}
}
