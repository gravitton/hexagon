package main

import (
	"math"

	geom "github.com/gravitton/geometry"
	"github.com/gravitton/grid"
	hex "github.com/gravitton/hexagon"
)

type orientation func(cellScale geom.Size[float64], system hex.CoordinateSystem) *grid.HexLayout

var (
	pointyTop orientation = grid.NewHexPointyTopLayout
	flatTop   orientation = grid.NewHexFlatTopLayout
)

type panel struct {
	canvas
	layout *grid.HexLayout
}

func newPanel(t theme, o orientation, size float64, caption string) *panel {
	return &panel{canvas: newCanvas(t, caption), layout: o(geom.SzU(size), hex.Axial)}
}

func backdrop(t theme, size float64, center hex.Hex, radius int, caption string) *panel {
	p := newPanel(t, pointyTop, size, caption)
	p.cells(center.Range(radius), t.cell)

	return p
}

func (p *panel) size() float64 {
	return p.layout.Polygon(hex.Pt(0, 0).Point()).BoundingCircle().Radius
}

func (p *panel) center(h hex.Hex) point {
	return p.layout.Center(h.Point())
}

func (p *panel) cell(h hex.Hex, fill string) {
	p.hexagon(h, style{fill: fill, stroke: p.theme.edge, width: 1})
}

func (p *panel) cells(hexes []hex.Hex, fill string) {
	for _, h := range hexes {
		p.cell(h, fill)
	}
}

func (p *panel) outline(h hex.Hex, color string, width float64) {
	p.hexagon(h, style{stroke: color, width: width})
}

func (p *panel) hexagon(h hex.Hex, s style) {
	p.polygon(p.layout.Polygon(h.Point()).Polygon().Points, s)
}

func (p *panel) text(h hex.Hex, dy float64, f font, spans ...span) {
	p.canvas.text(p.center(h).AddXY(0, dy), f, spans...)
}

func (p *panel) label(h hex.Hex, text, color string) {
	p.text(h, 0, font{size: math.Round(p.size() * 0.62)}, span{text, color})
}

func (p *panel) segment(a, b hex.Hex, color string, dashed bool) {
	p.line(p.center(a), p.center(b), style{stroke: color, width: 2, dashed: dashed})
}

func (p *panel) path(hexes []hex.Hex, color string) {
	points := make([]point, len(hexes))
	for i, h := range hexes {
		points[i] = p.center(h)
	}

	p.polyline(points, style{stroke: color, width: 2})
}

func (p *panel) pointer(from, to hex.Hex, color string, dashed bool) {
	a, b := p.center(from), p.center(to)
	along := b.Subtract(a).Normalize().Multiply(p.size() * 0.6)

	p.arrow(a.Add(along), b.Add(along.Negate()), style{stroke: color, width: 2, dashed: dashed})
}

func (p *panel) axis(h hex.Hex, reach float64, name string) {
	origin := p.center(hex.Pt(0, 0))
	along := p.center(h).Subtract(origin)
	p.line(origin.Add(along.Multiply(-reach)), origin.Add(along.Multiply(reach)), style{stroke: p.theme.muted, width: 1.2, dashed: true})
	p.canvas.text(origin.Add(along.Multiply(reach*1.1)), font{size: 11}, span{name, p.theme.muted})
}
