package hex

import (
	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/floats"
)

// Layout describes the mapping between hex coordinates and pixel space.
// It holds the hex orientation (flat-top or pointy-top), the hex Size
// (used as a scaling factor relative to a canonical hex), and the pixel
// Origin that corresponds to the axial hex (0,0) center.
type Layout struct {
	orientation orientation
	Size        floats.Size  // multiplication factor relative to the canonical hexagon
	Origin      floats.Point // center Point for hexagon with coordinates (0,0)
}

// LayoutFlatTop constructs a Layout for flat-top hexes with the given size and origin.
func LayoutFlatTop(size floats.Size, origin floats.Point) Layout {
	return Layout{orientationFlatTop, size, origin}
}

// LayoutPointyTop constructs a Layout for pointy-top hexes with the given size and origin.
func LayoutPointyTop(size floats.Size, origin floats.Point) Layout {
	return Layout{orientationPointyTop, size, origin}
}

// Bounds returns the pixel width/height of a single hex cell in this layout.
func (l Layout) Bounds() floats.Size {
	return l.Size.ScaleXY(l.orientation.bounds.XY())
}

// Spacing returns the horizontal and vertical distances between adjacent hex centers.
func (l Layout) Spacing() floats.Size {
	return l.Size.ScaleXY(l.orientation.spacing.XY())
}

// ToPoint converts a hex to the pixel coordinates of its center in the layout.
func (l Layout) ToPoint(hex Hex) floats.Point {
	//return floats.Pt(hex.QR()).Transform(l.orientation.toPoint)
	return l.Origin.Add(floats.Vec(hex.QR()).Transform(l.orientation.toPoint).MultiplyXY(l.Size.XY()))
}

// FromPoint converts a pixel point to a fractional hex in the layout.
func (l Layout) FromPoint(point floats.Point) FractionalHex {
	//return F(point.Transform(l.orientation.fromPoint).XY())
	return F(point.Subtract(l.Origin).DivideXY(l.Size.XY()).Transform(l.orientation.fromPoint).XY())
}

// Hexagon creates a floats.RegularPolygon representing the hex cell in pixel space.
func (l Layout) Hexagon(h Hex) floats.RegularPolygon {
	return geom.Hexagon(l.ToPoint(h), l.Size, l.orientation.orientation)
}

type orientation struct {
	toPoint     geom.Matrix
	fromPoint   geom.Matrix
	orientation geom.Orientation
	bounds      floats.Vector
	spacing     floats.Vector
}

var orientationPointyTop = orientation{
	geom.Mat(geom.Sqrt3, geom.Sqrt3/2.0, 0, 0.0, 3.0/2.0, 0.0),
	geom.Mat(geom.Sqrt3/3.0, -1.0/3.0, 0, 0.0, 2.0/3.0, 0),
	geom.PointyTop,
	geom.Vec(geom.Sqrt3, 2.0),
	geom.Vec(geom.Sqrt3, 3.0/2.0),
}

var orientationFlatTop = orientation{
	geom.Mat(3.0/2.0, 0.0, 0, geom.Sqrt3/2.0, geom.Sqrt3, 0),
	geom.Mat(2.0/3.0, 0.0, 0, -1.0/3.0, geom.Sqrt3/3.0, 0),
	geom.FlatTop,
	geom.Vec(2.0, geom.Sqrt3),
	geom.Vec(3.0/2.0, geom.Sqrt3),
}
