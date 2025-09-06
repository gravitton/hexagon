package hex

import (
	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/floats"
)

type Layout struct {
	orientation orientation
	Size        floats.Size  // multiplication factor relative to the canonical hexagon
	Origin      floats.Point // center Point for hexagon with coordinates (0,0)
}

func LayoutFlatTop(size floats.Size, origin floats.Point) Layout {
	//o := orientationFlat
	//o.toPixel = o.toPixel.PreScale(size.XY()).PreTranslate(origin.XY())
	//o.fromPixel = o.fromPixel.Unscale(size.XY()).Untranslate(origin.XY())
	return Layout{orientationFlat, size, origin}
}

func LayoutPointyTop(size floats.Size, origin floats.Point) Layout {
	//o := orientationPointy
	//o.toPixel = o.toPixel.PreScale(size.XY()).PreTranslate(origin.XY())
	//o.fromPixel = o.fromPixel.Unscale(size.XY()).Untranslate(origin.XY())
	return Layout{orientationPointy, size, origin}
}

func (l Layout) Bounds() floats.Size {
	return l.Size.ScaleXY(l.orientation.bounds.XY())
}

func (l Layout) Spacing() floats.Size {
	return l.Size.ScaleXY(l.orientation.spacing.XY())
}

// ToPoint converts a hex to a pixel point of its center in the layout.
func (l Layout) ToPoint(hex Hex) floats.Point {
	return l.Origin.Add(floats.V(hex.QR()).Transform(l.orientation.toPixel).MultiplyXY(l.Size.XY()))
	//return floats.P(hex.QR()).Transform(l.orientation.toPixel)
}

// FromPoint converts a pixel point to a fractional hex in the layout.
func (l Layout) FromPoint(point floats.Point) FractionalHex {
	return F(point.Subtract(l.Origin).DivideXY(l.Size.XY()).Transform(l.orientation.fromPixel).XY())
	//return F(point.Transform(l.orientation.fromPixel).XY())
}

func (l Layout) Hexagon(h Hex) floats.RegularPolygon {
	return geom.Hexagon(l.ToPoint(h), l.Size, l.orientation.orientation)
}

type orientation struct {
	toPixel     geom.Matrix
	fromPixel   geom.Matrix
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
