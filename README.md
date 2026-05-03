# Hexagon

[![Latest Stable Version][ico-release]][link-release]
[![Build Status][ico-workflow]][link-workflow]
[![Coverage Status][ico-coverage]][link-coverage]
[![Go Report Card][ico-go-report-card]][link-go-report-card]
[![Go Dev Reference][ico-go-dev-reference]][link-go-dev-reference]
[![Software License][ico-license]][link-licence]

Hexagon library for game development.

- Axial (cube) coordinates with arithmetic, distance, and neighbor lookup
- Area and perimeter traversal: range, ring, and spiral
- Line drawing, line-of-sight, and field-of-view
- Rotation and reflection in cube space
- Direction enum with opposite and rotation helpers
- Seven coordinate systems with lossless round-trip conversion

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
a := hex.Pt(1, -2)
b := hex.Pt(0, 3)

// arithmetic and distance
c := a.Add(b)
d := a.Subtract(b)
e := a.Multiply(2)
distance := a.DistanceTo(b)

// neighbors
neighbors := a.Neighbors()
neighbor := a.Neighbor(hex.QPlus)

// area traversal
area := b.Range(2) // all hexes within radius 2 (filled)
ring := b.Ring(2) // hexes at exactly distance 2 (perimeter)
spiral := b.Spiral(2) // same as Range but ordered center-outward

// line drawing and visibility
line := a.Line(b)
visible := a.HasLineOfSight(b, blocking)
fov := a.FieldOfView(candidates, blocking)

// rotation and reflection (cube space)
rotated := a.Rotate(2) // 2×60° clockwise around origin
rotatedAround := a.RotateAround(b, -1) // 1×60° counterclockwise around b
reflectedQ := a.ReflectQ()
reflectedR := a.ReflectR()
reflectedS := a.ReflectS()

// directions
dir := hex.QPlus
opp := dir.Opposite() // SPlus
next := dir.Rotate(1) // RMinus (one step counterclockwise)
offset := dir.NeighborOffset() // axial vector for this direction

// coordinate system conversions
pOddR := a.To(hex.OffsetOddR)
pEvenQ := a.To(hex.OffsetEvenQ)
dw := a.To(hex.DoubleWidth)
dh := a.To(hex.DoubleHeight)
back := hex.From(pOddR, hex.OffsetOddR)
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