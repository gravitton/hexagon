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

func Test_String(t *testing.T) {
	assert.Equal(t, SMinus.String(), "SMinus")
	assert.Equal(t, QPlus.String(), "QPlus")
	assert.Equal(t, RMinus.String(), "RMinus")
	assert.Equal(t, SPlus.String(), "SPlus")
	assert.Equal(t, QMinus.String(), "QMinus")
	assert.Equal(t, RPlus.String(), "RPlus")
	// values beyond 5 wrap via %6
	assert.Equal(t, Direction(6).String(), "SMinus")
	assert.Equal(t, Direction(8).String(), "RMinus")
}

func TestDirections(t *testing.T) {
	assert.Equal(t, SMinus.NeighborOffset(), axialDirection[0])
	assert.Equal(t, QPlus.NeighborOffset(), axialDirection[1])
	assert.Equal(t, RMinus.NeighborOffset(), axialDirection[2])
	assert.Equal(t, SPlus.NeighborOffset(), axialDirection[3])
	assert.Equal(t, QMinus.NeighborOffset(), axialDirection[4])
	assert.Equal(t, RPlus.NeighborOffset(), axialDirection[5])

	// direction neighbor direction module
	assert.Equal(t, Direction(6).NeighborOffset(), axialDirection[0])
	assert.Equal(t, Direction(8).NeighborOffset(), axialDirection[2])
	assert.Equal(t, Direction(15).NeighborOffset(), axialDirection[3])

	assert.Equal(t, FlatTopSE, SMinus)
	assert.Equal(t, FlatTopNE, QPlus)
	assert.Equal(t, FlatTopN, RMinus)
	assert.Equal(t, FlatTopNW, SPlus)
	assert.Equal(t, FlatTopSW, QMinus)
	assert.Equal(t, FlatTopS, RPlus)

	assert.Equal(t, PointyTopE, SMinus)
	assert.Equal(t, PointyTopNE, QPlus)
	assert.Equal(t, PointyTopNW, RMinus)
	assert.Equal(t, PointyTopW, SPlus)
	assert.Equal(t, PointyTopSW, QMinus)
	assert.Equal(t, PointyTopSE, RPlus)
}

func TestDirection_Opposite(t *testing.T) {
	assert.Equal(t, SMinus.Opposite(), SPlus)
	assert.Equal(t, QPlus.Opposite(), QMinus)
	assert.Equal(t, RMinus.Opposite(), RPlus)
	assert.Equal(t, SPlus.Opposite(), SMinus)
	assert.Equal(t, QMinus.Opposite(), QPlus)
	assert.Equal(t, RPlus.Opposite(), RMinus)

	// double opposite returns the original direction
	for _, d := range []Direction{SMinus, QPlus, RMinus, SPlus, QMinus, RPlus} {
		assert.Equal(t, d.Opposite().Opposite(), d)
	}

	// verify neighbor offset vectors are truly opposite
	for _, d := range []Direction{SMinus, QPlus, RMinus} {
		v := d.NeighborOffset()
		opp := d.Opposite().NeighborOffset()
		assert.Equal(t, v.X+opp.X, 0)
		assert.Equal(t, v.Y+opp.Y, 0)
	}
}

func TestDirection_Rotate(t *testing.T) {
	// Rotate(0) is identity
	assert.Equal(t, SMinus.Rotate(0), SMinus)
	assert.Equal(t, RPlus.Rotate(0), RPlus)

	// Rotate(1) advances one step counterclockwise (increments index)
	assert.Equal(t, SMinus.Rotate(1), QPlus) // 0 → 1
	assert.Equal(t, QPlus.Rotate(1), RMinus) // 1 → 2
	assert.Equal(t, RPlus.Rotate(1), SMinus) // 5 → 0 (wrap)

	// Rotate(-1) advances one step clockwise (decrements index)
	assert.Equal(t, SMinus.Rotate(-1), RPlus) // 0 → 5 (wrap)
	assert.Equal(t, QPlus.Rotate(-1), SMinus) // 1 → 0
	assert.Equal(t, RMinus.Rotate(-1), QPlus) // 2 → 1

	// Rotate(3) == Opposite()
	for _, d := range []Direction{SMinus, QPlus, RMinus, SPlus, QMinus, RPlus} {
		assert.Equal(t, d.Rotate(3), d.Opposite())
	}

	// Rotate(6) is identity
	for _, d := range []Direction{SMinus, QPlus, RMinus, SPlus, QMinus, RPlus} {
		assert.Equal(t, d.Rotate(6), d)
	}

	// Rotate and its inverse cancel out
	for _, d := range []Direction{SMinus, QPlus, RMinus, SPlus, QMinus, RPlus} {
		assert.Equal(t, d.Rotate(2).Rotate(-2), d)
	}
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
			index:        geom.Pt(0, 0), // even col, even row
			axial:        axialDirection,
			offsetOddR:   offsetOddRDirectionEvenRow,
			offsetEvenR:  offsetEvenRDirectionEvenRow,
			offsetOddQ:   offsetOddQDirectionEvenCol,
			offsetEvenQ:  offsetEvenQDirectionEvenCol,
			doubleWidth:  doubleWidthDirection,
			doubleHeight: doubleHeightDirection,
		},
		{
			index:        geom.Pt(1, 1), // odd col, odd row
			axial:        axialDirection,
			offsetOddR:   offsetOddRDirectionOddRow,
			offsetEvenR:  offsetEvenRDirectionOddRow,
			offsetOddQ:   offsetOddQDirectionOddCol,
			offsetEvenQ:  offsetEvenQDirectionOddCol,
			doubleWidth:  doubleWidthDirection,
			doubleHeight: doubleHeightDirection,
		},
		{
			index:        geom.Pt(3, 2), // odd col, even row
			axial:        axialDirection,
			offsetOddR:   offsetOddRDirectionEvenRow,
			offsetEvenR:  offsetEvenRDirectionEvenRow,
			offsetOddQ:   offsetOddQDirectionOddCol,
			offsetEvenQ:  offsetEvenQDirectionOddCol,
			doubleWidth:  doubleWidthDirection,
			doubleHeight: doubleHeightDirection,
		},
		{
			index:        geom.Pt(-2, -3), // even col, odd row
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
				assert.Equal(t, NeighborOffset(test.index, system, SMinus), offsets[0])
				assert.Equal(t, NeighborOffset(test.index, system, QPlus), offsets[1])
				assert.Equal(t, NeighborOffset(test.index, system, RMinus), offsets[2])
				assert.Equal(t, NeighborOffset(test.index, system, SPlus), offsets[3])
				assert.Equal(t, NeighborOffset(test.index, system, QMinus), offsets[4])
				assert.Equal(t, NeighborOffset(test.index, system, RPlus), offsets[5])
			}
		})
	}
}
