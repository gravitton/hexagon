package hex_test

import (
	"encoding/json"
	"fmt"
	"slices"
	"testing"

	"github.com/gravitton/assert"
	. "github.com/gravitton/hexagon"
)

var (
	testRegion      = Reg(testHexZero.Range(2))
	testRegionOther = Reg(Pt(3, -1).Range(2))
	testRegionEmpty = Region{}
)

func TestRegion_Constructor(t *testing.T) {
	t.Run("orders the hexes and drops repeats", func(t *testing.T) {
		assert.Equal(t, Reg([]Hex{Pt(1, 2), Pt(-1, 0), Pt(1, 2), Pt(1, -3)}).Hexes(), []Hex{Pt(-1, 0), Pt(1, -3), Pt(1, 2)})
	})
	t.Run("no hexes make the empty region", func(t *testing.T) {
		assert.True(t, Reg(nil).Equal(testRegionEmpty))
		assert.True(t, Reg([]Hex{}).Equal(testRegionEmpty))
	})
	t.Run("copies the slice", func(t *testing.T) {
		hexes := []Hex{Pt(0, 0), Pt(1, 0)}
		region := Reg(hexes)
		hexes[0] = Pt(9, 9)

		assert.Equal(t, region.Hexes(), []Hex{Pt(0, 0), Pt(1, 0)})
	})
	t.Run("the same hexes in any order are the same region", func(t *testing.T) {
		assert.True(t, Reg(testHexZero.Spiral(2)).Equal(testRegion))
	})
	t.Run("allocates once", func(t *testing.T) {
		hexes := testHexZero.Spiral(2)
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkRegion = Reg(hexes)
		}), 1.0)
	})
}

func TestRegion_Hexes(t *testing.T) {
	t.Run("ordered like Compare", func(t *testing.T) {
		assert.Equal(t, testRegion.Hexes(), testHexZero.Range(2))
	})
	t.Run("returns a fresh slice", func(t *testing.T) {
		hexes := testRegion.Hexes()
		hexes[0] = Pt(9, 9)

		assert.False(t, testRegion.Contains(Pt(9, 9)))
		assert.Equal(t, testRegion.Hexes(), testHexZero.Range(2))
	})
	t.Run("empty is nil", func(t *testing.T) {
		assert.Equal(t, testRegionEmpty.Hexes(), nil)
	})
	t.Run("allocates once", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkHexes = testRegion.Hexes()
		}), 1.0)
	})
}

func TestRegion_AppendHexes(t *testing.T) {
	t.Run("appends after dst", func(t *testing.T) {
		assert.Equal(t, testRegion.AppendHexes([]Hex{testHex}), append([]Hex{testHex}, testRegion.Hexes()...))
	})
	t.Run("empty appends nothing", func(t *testing.T) {
		assert.Equal(t, testRegionEmpty.AppendHexes([]Hex{testHex}), []Hex{testHex})
	})
	t.Run("allocates nothing with room", func(t *testing.T) {
		buffer := make([]Hex, 0, 19)
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkHexes = testRegion.AppendHexes(buffer)
		}), 0.0)
	})
}

func TestRegion_All(t *testing.T) {
	t.Run("yields the hexes of Hexes in order", func(t *testing.T) {
		assert.Equal(t, slices.Collect(testRegion.All()), testRegion.Hexes())
	})
	t.Run("empty yields nothing", func(t *testing.T) {
		for range testRegionEmpty.All() {
			assert.Fail(t, "yielded a hex")
		}
	})
	t.Run("allocates nothing", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			for h := range testRegion.All() {
				sinkHex = h
			}
		}), 0.0)
	})
}

func TestRegion_Len(t *testing.T) {
	t.Run("the number of hexes", func(t *testing.T) {
		assert.Equal(t, testRegion.Len(), 19)
		assert.Equal(t, Reg([]Hex{testHex, testHex}).Len(), 1)
	})
	t.Run("empty", func(t *testing.T) {
		assert.Equal(t, testRegionEmpty.Len(), 0)
	})
}

func TestRegion_Union(t *testing.T) {
	union := testRegion.Union(testRegionOther)

	t.Run("holds exactly the hexes of either", func(t *testing.T) {
		for _, h := range testHexZero.Range(7) {
			assert.Equal(t, union.Contains(h), testRegion.Contains(h) || testRegionOther.Contains(h), h.String()+": ")
		}
	})
	t.Run("stays ordered without repeats", func(t *testing.T) {
		assert.True(t, Reg(union.Hexes()).Equal(union))
		assert.True(t, slices.IsSortedFunc(union.Hexes(), Hex.Compare))
	})
	t.Run("commutative", func(t *testing.T) {
		assert.True(t, testRegionOther.Union(testRegion).Equal(union))
	})
	t.Run("with itself or the empty region changes nothing", func(t *testing.T) {
		assert.True(t, testRegion.Union(testRegion).Equal(testRegion))
		assert.True(t, testRegion.Union(testRegionEmpty).Equal(testRegion))
		assert.True(t, testRegionEmpty.Union(testRegion).Equal(testRegion))
	})
	t.Run("of two empty regions is empty", func(t *testing.T) {
		assert.True(t, testRegionEmpty.Union(testRegionEmpty).IsEmpty())
	})
	t.Run("counts as the two less what they share", func(t *testing.T) {
		assert.Equal(t, union.Len(), testRegion.Len()+testRegionOther.Len()-testRegion.Intersection(testRegionOther).Len())
	})
	t.Run("allocates once at the exact size", func(t *testing.T) {
		assert.Equal(t, cap(union.AppendHexes(nil)), union.Len())
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkRegion = testRegion.Union(testRegionOther)
		}), 1.0)
	})
}

func BenchmarkRegion_Union(b *testing.B) {
	for b.Loop() {
		sinkRegion = testRegion.Union(testRegionOther)
	}
}

func TestRegion_Intersection(t *testing.T) {
	intersection := testRegion.Intersection(testRegionOther)

	t.Run("holds exactly the hexes of both", func(t *testing.T) {
		for _, h := range testHexZero.Range(7) {
			assert.Equal(t, intersection.Contains(h), testRegion.Contains(h) && testRegionOther.Contains(h), h.String()+": ")
		}
	})
	t.Run("of two ranges is their RangeIntersection", func(t *testing.T) {
		assert.Equal(t, intersection.Hexes(), testHexZero.RangeIntersection(2, Pt(3, -1), 2))
	})
	t.Run("commutative", func(t *testing.T) {
		assert.True(t, testRegionOther.Intersection(testRegion).Equal(intersection))
	})
	t.Run("with itself changes nothing", func(t *testing.T) {
		assert.True(t, testRegion.Intersection(testRegion).Equal(testRegion))
	})
	t.Run("with the empty region or one apart is empty", func(t *testing.T) {
		assert.True(t, testRegion.Intersection(testRegionEmpty).IsEmpty())
		assert.True(t, testRegionEmpty.Intersection(testRegion).IsEmpty())
		assert.True(t, testRegion.Intersection(Reg(Pt(9, 9).Range(1))).IsEmpty())
	})
	t.Run("allocates once", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkRegion = testRegion.Intersection(testRegionOther)
		}), 1.0)
	})
}

func BenchmarkRegion_Intersection(b *testing.B) {
	for b.Loop() {
		sinkRegion = testRegion.Intersection(testRegionOther)
	}
}

func TestRegion_Difference(t *testing.T) {
	difference := testRegion.Difference(testRegionOther)

	t.Run("holds exactly the hexes of the first alone", func(t *testing.T) {
		for _, h := range testHexZero.Range(7) {
			assert.Equal(t, difference.Contains(h), testRegion.Contains(h) && !testRegionOther.Contains(h), h.String()+": ")
		}
	})
	t.Run("with the intersection makes the region whole again", func(t *testing.T) {
		assert.True(t, difference.Union(testRegion.Intersection(testRegionOther)).Equal(testRegion))
	})
	t.Run("of itself is empty", func(t *testing.T) {
		assert.True(t, testRegion.Difference(testRegion).IsEmpty())
	})
	t.Run("of the empty region changes nothing", func(t *testing.T) {
		assert.True(t, testRegion.Difference(testRegionEmpty).Equal(testRegion))
		assert.True(t, testRegionEmpty.Difference(testRegion).IsEmpty())
	})
	t.Run("allocates once", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkRegion = testRegion.Difference(testRegionOther)
		}), 1.0)
	})
}

func TestRegion_Border(t *testing.T) {
	t.Run("of a range is its ring", func(t *testing.T) {
		for _, radius := range []int{0, 1, 2, 4} {
			assert.True(t, Reg(testHex.Range(radius)).Border().Equal(Reg(testHex.Ring(radius))))
		}
	})
	t.Run("of a ring is the ring itself", func(t *testing.T) {
		ring := Reg(testHex.Ring(3))
		assert.True(t, ring.Border().Equal(ring))
	})
	t.Run("runs around a hole as well", func(t *testing.T) {
		holed := Reg(testHexZero.Range(3)).Difference(Reg([]Hex{testHexZero}))
		assert.True(t, holed.Border().Equal(Reg(testHexZero.Ring(3)).Union(Reg(testHexZero.Ring(1)))))
	})
	t.Run("lies inside the region", func(t *testing.T) {
		union := testRegion.Union(testRegionOther)
		assert.True(t, union.Border().Difference(union).IsEmpty())
	})
	t.Run("holds exactly the hexes with a neighbor outside", func(t *testing.T) {
		holed := testRegion.Union(testRegionOther).Difference(Reg([]Hex{testHex}))
		border := holed.Border()
		for h := range holed.All() {
			neighbors := h.Neighbors()
			outside := slices.ContainsFunc(neighbors[:], func(neighbor Hex) bool {
				return !holed.Contains(neighbor)
			})
			assert.Equal(t, border.Contains(h), outside, h.String()+": ")
		}
	})
	t.Run("empty has none", func(t *testing.T) {
		assert.True(t, testRegionEmpty.Border().IsEmpty())
	})
	t.Run("allocates once", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkRegion = testRegion.Border()
		}), 1.0)
	})
}

func BenchmarkRegion_Border(b *testing.B) {
	for b.Loop() {
		sinkRegion = testRegion.Border()
	}
}

func TestRegion_Contains(t *testing.T) {
	t.Run("exactly the hexes it was made of", func(t *testing.T) {
		hexes := testHexZero.Range(2)
		for _, h := range testHexZero.Range(4) {
			assert.Equal(t, testRegion.Contains(h), slices.Contains(hexes, h), h.String()+": ")
		}
	})
	t.Run("empty holds nothing", func(t *testing.T) {
		assert.False(t, testRegionEmpty.Contains(testHexZero))
	})
	t.Run("allocates nothing", func(t *testing.T) {
		assert.Equal(t, testing.AllocsPerRun(100, func() {
			sinkBool = testRegion.Contains(testHex)
		}), 0.0)
	})
}

func BenchmarkRegion_Contains(b *testing.B) {
	for b.Loop() {
		sinkBool = testRegion.Contains(testHex)
	}
}

func TestRegion_Equal(t *testing.T) {
	t.Run("same hexes", func(t *testing.T) {
		assert.True(t, testRegion.Equal(Reg(testHexZero.Range(2))))
		assert.True(t, testRegionEmpty.Equal(Reg(nil)))
	})
	t.Run("one hex differs", func(t *testing.T) {
		assert.False(t, testRegion.Equal(testRegion.Difference(Reg([]Hex{testHexZero}))))
		assert.False(t, testRegion.Equal(testRegionOther))
		assert.False(t, testRegion.Equal(testRegionEmpty))
	})
}

func TestRegion_IsEmpty(t *testing.T) {
	t.Run("the zero value", func(t *testing.T) {
		assert.True(t, testRegionEmpty.IsEmpty())
	})
	t.Run("any hex makes it not empty", func(t *testing.T) {
		assert.False(t, Reg([]Hex{testHex}).IsEmpty())
	})
}

func TestRegion_String(t *testing.T) {
	t.Run("the hexes in order", func(t *testing.T) {
		assert.Equal(t, Reg([]Hex{Pt(1, 0), Pt(0, 0)}).String(), "Reg((0,0);(1,0))")
	})
	t.Run("empty", func(t *testing.T) {
		assert.Equal(t, testRegionEmpty.String(), "Reg()")
	})
}

func TestRegion_JSON(t *testing.T) {
	t.Run("round-trips as an array of hexes", func(t *testing.T) {
		data, err := json.Marshal(Reg([]Hex{Pt(1, 0), Pt(0, 0)}))
		assert.NoError(t, err)
		assert.Equal(t, string(data), `[{"q":0,"r":0},{"q":1,"r":0}]`)

		var decoded Region
		assert.NoError(t, json.Unmarshal(data, &decoded))
		assert.True(t, decoded.Equal(Reg([]Hex{Pt(0, 0), Pt(1, 0)})))
	})
	t.Run("empty is an empty array", func(t *testing.T) {
		data, err := json.Marshal(testRegionEmpty)
		assert.NoError(t, err)
		assert.Equal(t, string(data), "[]")

		decoded := testRegion
		assert.NoError(t, json.Unmarshal(data, &decoded))
		assert.True(t, decoded.IsEmpty())
	})
	t.Run("orders and drops repeats on the way in", func(t *testing.T) {
		var decoded Region
		assert.NoError(t, json.Unmarshal([]byte(`[{"q":1,"r":0},{"q":0,"r":0},{"q":1,"r":0}]`), &decoded))
		assert.Equal(t, decoded.Hexes(), []Hex{Pt(0, 0), Pt(1, 0)})
	})
	t.Run("anything but an array of hexes fails", func(t *testing.T) {
		var decoded Region
		assert.Error(t, json.Unmarshal([]byte(`{"q":1}`), &decoded))
	})
}

func ExampleReg() {
	region := Reg([]Hex{Pt(1, 0), Pt(0, 0), Pt(1, 0)})

	fmt.Println(region)
	fmt.Println(region.Len(), region.Contains(Pt(1, 0)))
	// Output:
	// Reg((0,0);(1,0))
	// 2 true
}

func ExampleRegion_All() {
	for h := range Reg([]Hex{Pt(1, 0), Pt(0, 0)}).All() {
		fmt.Println(h)
	}
	// Output:
	// (0,0)
	// (1,0)
}

func ExampleRegion_Intersection() {
	reach := Reg(Pt(0, 0).Range(1))
	blast := Reg(Pt(1, 0).Range(1))

	fmt.Println(reach.Intersection(blast))
	fmt.Println(reach.Difference(blast).Len())
	fmt.Println(reach.Union(blast).Len())
	// Output:
	// Reg((0,0);(0,1);(1,-1);(1,0))
	// 3
	// 10
}

func ExampleRegion_Border() {
	fmt.Println(Reg(Pt(0, 0).Range(2)).Border().Len())
	// Output: 12
}
