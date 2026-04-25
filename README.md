# Hexagon

[![Latest Stable Version][ico-release]][link-release]
[![Build Status][ico-workflow]][link-workflow]
[![Coverage Status][ico-coverage]][link-coverage]
[![Go Report Card][ico-go-report-card]][link-go-report-card]
[![Go Dev Reference][ico-go-dev-reference]][link-go-dev-reference]
[![Software License][ico-license]][link-licence]

Hexagon library for game development.

- Convert between coordinate systems
- Map hexes to pixels (and back)

## Installation

```bash
go get github.com/gravitton/hexagon
```


## Usage

```go
import (
	hex "github.com/gravitton/hexagon"
	geom "github.com/gravitton/geometry"
)

// create hexagon in axial coordinates (q,r) 
a := hex.Coord(1, -2)
b := hex.Coord(0, 3)

// use basic math on hexes
c := a.Add(b)
d := a.Subtract(b)
e := a.Multiply(2)
distance := a.DistanceTo(b)

// neighbors and range
neighbors := a.Neighbors()
neighbor := a.Neighbor(hex.QPlus)
ring := b.Range(2)

// conversions between coordinate systems
pOddR := a.To(hex.OffsetOddR)
pEvenQ := a.To(hex.OffsetEvenQ)
dw := a.To(hex.DoubleWidth)
dh := a.To(hex.DoubleHeight)

// hexagon map layout
layout := hex.LayoutFlatTop(geom.Sz(16, 16), geom.Pt(100, 80))
center := layout.ToPoint(a)
clicked := layout.FromPoint(geom.Pt(499.0, 123.4)).Round()
```


## Credits

- [Tomáš Novotný](https://github.com/tomas-novotny)
- [All Contributors][link-contributors]
- [Red Blob Games](https://www.redblobgames.com/grids/hexagons)


## License

The MIT License (MIT). Please see [License File][link-licence] for more information.


[ico-license]:              https://img.shields.io/github/license/gravitton/hexagon.svg?style=flat-square&colorB=blue
[ico-workflow]:             https://img.shields.io/github/actions/workflow/status/gravitton/hexagon/main.yml?branch=main&style=flat-square
[ico-release]:              https://img.shields.io/github/v/release/gravitton/hexagon?style=flat-square&colorB=blue
[ico-go-dev-reference]:     https://img.shields.io/badge/go.dev-reference-blue?style=flat-square
[ico-go-report-card]:       https://goreportcard.com/badge/github.com/gravitton/hexagon?style=flat-square
[ico-coverage]:             https://img.shields.io/coverallsCoverage/github/gravitton/hexagon?style=flat-square

[link-author]:              https://github.com/gravitton
[link-release]:             https://github.com/gravitton/hexagon/releases
[link-contributors]:        https://github.com/gravitton/hexagon/contributors
[link-licence]:             ./LICENSE.md
[link-changelog]:           ./CHANGELOG.md
[link-workflow]:            https://github.com/gravitton/hexagon/actions
[link-go-dev-reference]:    https://pkg.go.dev/github.com/gravitton/hexagon
[link-go-report-card]:      https://goreportcard.com/report/github.com/gravitton/hexagon
[link-coverage]:            https://coveralls.io/github/gravitton/hexagon