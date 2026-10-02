package hex_test

import (
	"encoding/json"
	"testing"

	"github.com/gravitton/assert"
	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/ints"
	. "github.com/gravitton/hexagon"
)

var testConversions = []struct {
	hex          Hex
	offsetOddR   ints.Point
	offsetEvenR  ints.Point
	offsetOddQ   ints.Point
	offsetEvenQ  ints.Point
	doubleWidth  ints.Point
	doubleHeight ints.Point
}{
	{
		hex:          Pt(0, 0),
		offsetOddR:   geom.Pt(0, 0),
		offsetEvenR:  geom.Pt(0, 0),
		offsetOddQ:   geom.Pt(0, 0),
		offsetEvenQ:  geom.Pt(0, 0),
		doubleWidth:  geom.Pt(0, 0),
		doubleHeight: geom.Pt(0, 0),
	},
	{
		hex:          Pt(1, 0),
		offsetOddR:   geom.Pt(1, 0),
		offsetEvenR:  geom.Pt(1, 0),
		offsetOddQ:   geom.Pt(1, 0),
		offsetEvenQ:  geom.Pt(1, 1),
		doubleWidth:  geom.Pt(2, 0),
		doubleHeight: geom.Pt(1, 1),
	},
	{
		hex:          Pt(1, -1),
		offsetOddR:   geom.Pt(0, -1),
		offsetEvenR:  geom.Pt(1, -1),
		offsetOddQ:   geom.Pt(1, -1),
		offsetEvenQ:  geom.Pt(1, 0),
		doubleWidth:  geom.Pt(1, -1),
		doubleHeight: geom.Pt(1, -1),
	},
	{
		hex:          Pt(0, -1),
		offsetOddR:   geom.Pt(-1, -1),
		offsetEvenR:  geom.Pt(0, -1),
		offsetOddQ:   geom.Pt(0, -1),
		offsetEvenQ:  geom.Pt(0, -1),
		doubleWidth:  geom.Pt(-1, -1),
		doubleHeight: geom.Pt(0, -2),
	},
	{
		hex:          Pt(-1, 0),
		offsetOddR:   geom.Pt(-1, 0),
		offsetEvenR:  geom.Pt(-1, 0),
		offsetOddQ:   geom.Pt(-1, -1),
		offsetEvenQ:  geom.Pt(-1, 0),
		doubleWidth:  geom.Pt(-2, 0),
		doubleHeight: geom.Pt(-1, -1),
	},
	{
		hex:          Pt(-1, 1),
		offsetOddR:   geom.Pt(-1, 1),
		offsetEvenR:  geom.Pt(-0, 1),
		offsetOddQ:   geom.Pt(-1, 0),
		offsetEvenQ:  geom.Pt(-1, 1),
		doubleWidth:  geom.Pt(-1, 1),
		doubleHeight: geom.Pt(-1, 1),
	},
	{
		hex:          Pt(0, 1),
		offsetOddR:   geom.Pt(0, 1),
		offsetEvenR:  geom.Pt(1, 1),
		offsetOddQ:   geom.Pt(0, 1),
		offsetEvenQ:  geom.Pt(0, 1),
		doubleWidth:  geom.Pt(1, 1),
		doubleHeight: geom.Pt(0, 2),
	},
}

func TestCoordinateSystems(t *testing.T) {
	t.Run("the seven in order", func(t *testing.T) {
		assert.Equal(t, CoordinateSystems(), [7]CoordinateSystem{Axial, OffsetOddR, OffsetEvenR, OffsetOddQ, OffsetEvenQ, DoubleWidth, DoubleHeight})
	})
	t.Run("returns a fresh array", func(t *testing.T) {
		list := CoordinateSystems()
		list[0] = DoubleWidth
		assert.Equal(t, CoordinateSystems()[0], Axial)
	})
}

func TestParseCoordinateSystem(t *testing.T) {
	t.Run("round-trips with String", func(t *testing.T) {
		for _, system := range CoordinateSystems() {
			parsed, err := ParseCoordinateSystem(system.String())
			assert.NoError(t, err)
			assert.Equal(t, parsed, system)
		}
	})
	t.Run("none", func(t *testing.T) {
		parsed, err := ParseCoordinateSystem("None")
		assert.NoError(t, err)
		assert.Equal(t, parsed, CoordinateSystemNone)
	})
	t.Run("unknown name", func(t *testing.T) {
		parsed, err := ParseCoordinateSystem("Cube")
		assert.Error(t, err)
		assert.Equal(t, parsed, CoordinateSystemNone)
	})
}

func TestCoordinateSystem_Offsets(t *testing.T) {
	t.Run("the tables at each parity", func(t *testing.T) {
		tests := []struct {
			index        ints.Point
			axial        [6]ints.Vector
			offsetOddR   [6]ints.Vector
			offsetEvenR  [6]ints.Vector
			offsetOddQ   [6]ints.Vector
			offsetEvenQ  [6]ints.Vector
			doubleWidth  [6]ints.Vector
			doubleHeight [6]ints.Vector
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
				assert.Equal(t, Axial.Offsets(test.index), test.axial)
				assert.Equal(t, OffsetOddR.Offsets(test.index), test.offsetOddR)
				assert.Equal(t, OffsetEvenR.Offsets(test.index), test.offsetEvenR)
				assert.Equal(t, OffsetOddQ.Offsets(test.index), test.offsetOddQ)
				assert.Equal(t, OffsetEvenQ.Offsets(test.index), test.offsetEvenQ)
				assert.Equal(t, DoubleWidth.Offsets(test.index), test.doubleWidth)
				assert.Equal(t, DoubleHeight.Offsets(test.index), test.doubleHeight)
			})
		}
	})
	t.Run("derived from Directions and the conversion", func(t *testing.T) {
		// reordering the directions without permuting every table in lockstep must fail here
		for _, system := range CoordinateSystems() {
			for _, h := range testHexZero.Spiral(8) {
				index := system.To(h)
				assert.Equal(t, system.Offsets(index), deriveOffsets(system, index), system.String(), index.String())
			}
		}
	})
	t.Run("panics for a system outside the seven", func(t *testing.T) {
		assert.PanicsWith(t, func() {
			CoordinateSystem(99).Offsets(geom.Pt(0, 0))
		}, "hex: unknown coordinate system 99")
	})
}

func TestCoordinateSystem_Offset(t *testing.T) {
	indexes := []ints.Point{geom.Pt(0, 0), geom.Pt(1, 1), geom.Pt(3, 2), geom.Pt(-2, -3)}

	t.Run("reads Offsets by direction", func(t *testing.T) {
		for _, system := range CoordinateSystems() {
			for _, index := range indexes {
				offsets := system.Offsets(index)
				for i, direction := range Directions() {
					assert.Equal(t, system.Offset(index, direction), offsets[i], system.String())
				}
			}
		}
	})
	t.Run("out-of-range directions wrap", func(t *testing.T) {
		for _, system := range CoordinateSystems() {
			for _, index := range indexes {
				assert.Equal(t, system.Offset(index, Direction(6)), system.Offset(index, SMinus), system.String())
				assert.Equal(t, system.Offset(index, Direction(-2)), system.Offset(index, RMinus), system.String())
			}
		}
	})
	t.Run("none steps nowhere", func(t *testing.T) {
		for _, system := range CoordinateSystems() {
			assert.Equal(t, system.Offset(geom.Pt(1, 1), DirectionNone), ints.Vector{}, system.String())
		}
	})
	t.Run("panics for a system outside the seven", func(t *testing.T) {
		assert.PanicsWith(t, func() {
			CoordinateSystemNone.Offset(geom.Pt(0, 0), SMinus)
		}, "hex: unknown coordinate system -1")
	})
}

func TestCoordinateSystem_To(t *testing.T) {
	t.Run("the seven systems", func(t *testing.T) {
		for _, test := range testConversions {
			t.Run(test.hex.String(), func(t *testing.T) {
				assert.Equal(t, Axial.To(test.hex), geom.Pt(test.hex.Q, test.hex.R))
				assert.Equal(t, OffsetOddR.To(test.hex), test.offsetOddR)
				assert.Equal(t, OffsetEvenR.To(test.hex), test.offsetEvenR)
				assert.Equal(t, OffsetOddQ.To(test.hex), test.offsetOddQ)
				assert.Equal(t, OffsetEvenQ.To(test.hex), test.offsetEvenQ)
				assert.Equal(t, DoubleWidth.To(test.hex), test.doubleWidth)
				assert.Equal(t, DoubleHeight.To(test.hex), test.doubleHeight)
			})
		}
	})
	t.Run("panics for a system outside the seven", func(t *testing.T) {
		assert.PanicsWith(t, func() {
			CoordinateSystem(99).To(testHexZero)
		}, "hex: unknown coordinate system 99")
	})
}

func TestCoordinateSystem_From(t *testing.T) {
	t.Run("the seven systems", func(t *testing.T) {
		for _, test := range testConversions {
			t.Run(test.hex.String(), func(t *testing.T) {
				assert.Equal(t, Axial.From(geom.Pt(test.hex.Q, test.hex.R)), test.hex)
				assert.Equal(t, OffsetOddR.From(test.offsetOddR), test.hex)
				assert.Equal(t, OffsetEvenR.From(test.offsetEvenR), test.hex)
				assert.Equal(t, OffsetOddQ.From(test.offsetOddQ), test.hex)
				assert.Equal(t, OffsetEvenQ.From(test.offsetEvenQ), test.hex)
				assert.Equal(t, DoubleWidth.From(test.doubleWidth), test.hex)
				assert.Equal(t, DoubleHeight.From(test.doubleHeight), test.hex)
			})
		}
	})
	t.Run("round-trips with To", func(t *testing.T) {
		for _, system := range CoordinateSystems() {
			for _, h := range testHex.Spiral(4) {
				assert.Equal(t, system.From(system.To(h)), h, system.String(), h.String())
			}
		}
	})
	t.Run("a double coordinate of mixed parity lands beside a cell", func(t *testing.T) {
		for _, system := range []CoordinateSystem{DoubleWidth, DoubleHeight} {
			for x := -3; x <= 3; x++ {
				for y := -3; y <= 3; y++ {
					index := geom.Pt(x, y)
					cell := system.To(system.From(index))
					assert.True(t, geom.Abs(cell.X-x)+geom.Abs(cell.Y-y) <= 1, system.String(), index.String())
				}
			}
		}
	})
	t.Run("panics for a system outside the seven", func(t *testing.T) {
		assert.PanicsWith(t, func() {
			CoordinateSystem(99).From(geom.Pt(0, 0))
		}, "hex: unknown coordinate system 99")
	})
}

func TestCoordinateSystem_IsNone(t *testing.T) {
	t.Run("none and every value outside the seven", func(t *testing.T) {
		assert.True(t, CoordinateSystemNone.IsNone())
		assert.True(t, CoordinateSystem(99).IsNone())
	})
	t.Run("the seven are not none", func(t *testing.T) {
		for _, system := range CoordinateSystems() {
			assert.False(t, system.IsNone(), system.String())
		}
	})
}

func TestCoordinateSystem_String(t *testing.T) {
	t.Run("the constant names", func(t *testing.T) {
		assert.Equal(t, Axial.String(), "Axial")
		assert.Equal(t, OffsetOddR.String(), "OffsetOddR")
		assert.Equal(t, OffsetEvenR.String(), "OffsetEvenR")
		assert.Equal(t, OffsetOddQ.String(), "OffsetOddQ")
		assert.Equal(t, OffsetEvenQ.String(), "OffsetEvenQ")
		assert.Equal(t, DoubleWidth.String(), "DoubleWidth")
		assert.Equal(t, DoubleHeight.String(), "DoubleHeight")
	})
	t.Run("none and every value outside the seven", func(t *testing.T) {
		assert.Equal(t, CoordinateSystemNone.String(), "None")
		assert.Equal(t, CoordinateSystem(99).String(), "None")
	})
}

func TestCoordinateSystem_Text(t *testing.T) {
	t.Run("round-trips through JSON as the name", func(t *testing.T) {
		systems := CoordinateSystems()
		for _, system := range append(systems[:], CoordinateSystemNone) {
			data, err := json.Marshal(system)
			assert.NoError(t, err)
			assert.Equal(t, string(data), `"`+system.String()+`"`)

			var decoded CoordinateSystem
			assert.NoError(t, json.Unmarshal(data, &decoded))
			assert.Equal(t, decoded, system)
		}
	})
	t.Run("unknown name fails", func(t *testing.T) {
		var decoded CoordinateSystem
		assert.Error(t, json.Unmarshal([]byte(`"Cube"`), &decoded))
	})
}

// deriveOffsets computes the neighbor offsets of a system straight from Directions and the
// coordinate conversion, which is the definition the literal tables cache.
func deriveOffsets(system CoordinateSystem, index ints.Point) [6]ints.Vector {
	h := system.From(index)

	var offsets [6]ints.Vector
	for i, direction := range Directions() {
		neighbor := system.To(h.Neighbor(direction))
		offsets[i] = geom.Vec(neighbor.X-index.X, neighbor.Y-index.Y)
	}

	return offsets
}
