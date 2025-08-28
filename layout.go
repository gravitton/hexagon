package hex

import (
	geom "github.com/gravitton/geometry"
	"github.com/gravitton/geometry/types/floats"
)

type Layout struct {
	Orientation orientation
	Size        floats.Size  // multiplication factor relative to the canonical hexagon
	Origin      floats.Point // center Point for hexagon with coordinates (0,0)
}

func LayoutFlatTop(size floats.Size, origin floats.Point) Layout {
	return Layout{OrientationFlat, size, origin}
}

func LayoutPointyTop(size floats.Size, origin floats.Point) Layout {
	return Layout{orientationPointy, size, origin}
}

func (l Layout) Bounds() floats.Size {
	return l.Size.ScaleXY(l.Orientation.bounds.XY())
}

func (l Layout) Spacing() floats.Size {
	return l.Size.ScaleXY(l.Orientation.spacing.XY())
}

func (l Layout) HexToPixel(h Hex) floats.Point {
	vector := geom.V(float64(h.Q), float64(h.R))

	x := l.Orientation.hexToX.Dot(vector)
	y := l.Orientation.hexToY.Dot(vector)

	return l.Origin.AddXY(x, y).MultiplyXY(l.Size.XY())
}

func (l Layout) PixelToHex(p floats.Point) FractionalHex {
	point := p.Subtract(l.Origin).DivideXY(l.Size.XY())

	q := point.Dot(l.Orientation.pixelToQ)
	r := point.Dot(l.Orientation.pixelToR)

	return FractionalHex{q, r}
}

func (l Layout) Hexagon(h Hex) floats.RegularPolygon {
	// TODO: use l.orientation.startAngle
	return geom.Hexagon(l.HexToPixel(h), l.Size)
}

type orientation struct {
	hexToX, hexToY     floats.Vector // TODO: use matrix
	pixelToQ, pixelToR floats.Vector // TODO: use matrix
	startAngle         float64
	bounds             floats.Vector
	spacing            floats.Vector
}

var orientationPointy = orientation{
	geom.V(geom.Sqrt3, geom.Sqrt3/2.0),
	geom.V(0.0, 3.0/2.0),
	geom.V(geom.Sqrt3/3.0, -1.0/3.0),
	geom.V(0.0, 2.0/3.0),
	0.5, // 90 degrees
	geom.V(geom.Sqrt3, 2.0),
	geom.V(geom.Sqrt3, 3.0/2.0),
}

var OrientationFlat = orientation{
	geom.V(3.0/2.0, 0.0),
	geom.V(geom.Sqrt3/2.0, geom.Sqrt3),
	geom.V(2.0/3.0, 0.0),
	geom.V(-1.0/3.0, geom.Sqrt3/3.0),
	0.0, // 0 degrees
	geom.V(2.0, geom.Sqrt3),
	geom.V(3.0/2.0, geom.Sqrt3),
}
