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
	//o := orientationFlat
	//o.toPoint = o.toPoint.PreScale(size.XY()).PreTranslate(origin.XY())
	//o.fromPoint = o.fromPoint.Unscale(size.XY()).Untranslate(origin.XY())
	return Layout{orientationFlat, size, origin}
}

// LayoutPointyTop constructs a Layout for pointy-top hexes with the given size and origin.
func LayoutPointyTop(size floats.Size, origin floats.Point) Layout {
	//o := orientationPointy
	//o.toPoint = o.toPoint.PreScale(size.XY()).PreTranslate(origin.XY())
	//o.fromPoint = o.fromPoint.Unscale(size.XY()).Untranslate(origin.XY())
	return Layout{orientationPointy, size, origin}
}

// Bounds returns the pixel width/height of a single hex cell in this layout.
func (l Layout) Bounds() floats.Size {
	return l.Size.ScaleXY(l.orientation.bounds.XY())
}

// Spacing returns the horizontal and vertical distances between adjacent hex centers.
func (l Layout) Spacing() floats.Size {
	return l.Size.ScaleXY(l.orientation.spacing.XY())
}

// ToPoint converts a hex to a pixel point of its center in the layout.
func (l Layout) ToPoint(hex Hex) floats.Point {
	return l.Origin.Add(floats.V(hex.QR()).Transform(l.orientation.toPoint).MultiplyXY(l.Size.XY()))
	//return floats.P(hex.QR()).Transform(l.orientation.toPoint)
}

// FromPoint converts a pixel point to a fractional hex in the layout.
func (l Layout) FromPoint(point floats.Point) FractionalHex {
	return F(point.Subtract(l.Origin).DivideXY(l.Size.XY()).Transform(l.orientation.fromPoint).XY())
	//return F(point.Transform(l.orientation.fromPoint).XY())
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

var orientationPointy = orientation{
	geom.M(geom.Sqrt3, geom.Sqrt3/2.0, 0, 0.0, 3.0/2.0, 0.0),
	geom.M(geom.Sqrt3/3.0, -1.0/3.0, 0, 0.0, 2.0/3.0, 0),
	geom.PointTop,
	geom.V(geom.Sqrt3, 2.0),
	geom.V(geom.Sqrt3, 3.0/2.0),
}

var orientationFlat = orientation{
	geom.M(3.0/2.0, 0.0, 0, geom.Sqrt3/2.0, geom.Sqrt3, 0),
	geom.M(2.0/3.0, 0.0, 0, -1.0/3.0, geom.Sqrt3/3.0, 0),
	geom.FlatTop,
	geom.V(2.0, geom.Sqrt3),
	geom.V(3.0/2.0, geom.Sqrt3),
}
