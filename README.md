# Hexagon

[![Latest Stable Version][ico-release]][link-release]
[![Build Status][ico-workflow]][link-workflow]
[![Coverage Status][ico-coverage]][link-coverage]
[![Go Dev Reference][ico-go-dev-reference]][link-go-dev-reference]
[![Software License][ico-license]][link-licence]

Hexagonal grid math library for game development.

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

## Core concepts

The library uses **axial coordinates**: every hex is identified by two integers `q` and `r`. The third cube coordinate `s` is derived as `-q - r` and is never stored. All `Hex` operations return new values — the type is immutable.

A companion floating-point type, `FractionalHex`, is used for interpolation and for converting pixel coordinates to hex coordinates before rounding with `Round()`.

## Usage

```go
import hex "github.com/gravitton/hexagon"

// Construct hexes in axial coordinates (q, r)
a := hex.Pt(1, -2)
b := hex.Pt(0, 3)

// Arithmetic and distance
c := a.Add(b)
d := a.Subtract(b)
e := a.Multiply(2)
distance := a.DistanceTo(b) // 5

// Neighbors
neighbors := a.Neighbors()          // all 6 adjacent hexes
neighbor  := a.Neighbor(hex.QPlus)  // one specific neighbor

// Area traversal
area   := b.Range(2)   // all hexes within radius 2 (filled disc)
ring   := b.Ring(2)    // hexes at exactly distance 2 (perimeter)
spiral := b.Spiral(2)  // same set as Range, ordered center-outward

// Line drawing and visibility
line    := a.Line(b)
visible := a.HasLineOfSight(b, blocking)
fov     := a.FieldOfView(candidates, blocking)

// Rotation and reflection
rotated       := a.Rotate(2)          // 2×60° clockwise around origin
rotatedAround := a.RotateAround(b, -1) // 1×60° counterclockwise around b
reflectedQ    := a.ReflectQ()
reflectedR    := a.ReflectR()
reflectedS    := a.ReflectS()

// Directions
dir    := hex.QPlus
opp    := dir.Opposite()     // SPlus
next   := dir.Rotate(1)      // RMinus (one step counterclockwise)
offset := dir.NeighborOffset() // axial vector for this direction

// Coordinate system conversions
pOddR := a.To(hex.OffsetOddR)
pEvenQ := a.To(hex.OffsetEvenQ)
dw   := a.To(hex.DoubleWidth)
dh   := a.To(hex.DoubleHeight)
back := hex.From(pOddR, hex.OffsetOddR) // round-trips exactly

// FractionalHex — interpolation and pixel→hex rounding
frac    := hex.FracPt(1.4, -1.8)
rounded := frac.Round()               // nearest Hex
mid     := hex.FracPt(0, 0).Lerp(hex.FracPt(3, -1), 0.5)
```

## Directions

Six named directions in cube space, with flat-top and pointy-top aliases:

| Constant | Index | Flat-top alias | Pointy-top alias |
|---|---|---|---|
| `SMinus` | 0 | `FlatTopSE` | `PointyTopE`  |
| `QPlus`  | 1 | `FlatTopNE` | `PointyTopNE` |
| `RMinus` | 2 | `FlatTopN`  | `PointyTopNW` |
| `SPlus`  | 3 | `FlatTopNW` | `PointyTopW`  |
| `QMinus` | 4 | `FlatTopSW` | `PointyTopSW` |
| `RPlus`  | 5 | `FlatTopS`  | `PointyTopSE` |

Directions are ordered counterclockwise starting from SE (flat-top) / E (pointy-top). `Direction.Rotate(n)` advances by `n` steps in the same order; negative steps go clockwise. `Direction.Opposite()` returns the direction 180° away.

## Coordinate systems

Seven systems are supported. Use `Hex.To(system)` / `hex.From(point, system)` to convert. All conversions are lossless round-trips.

| Constant | Orientation | Description |
|---|---|---|
| `Axial`        | — | Native storage: `(q, r)` directly |
| `OffsetOddR`   | Pointy-top | Odd rows shifted right |
| `OffsetEvenR`  | Pointy-top | Even rows shifted right |
| `OffsetOddQ`   | Flat-top | Odd columns shifted down |
| `OffsetEvenQ`  | Flat-top | Even columns shifted down |
| `DoubleWidth`  | Pointy-top | Column axis doubled; no parity split needed |
| `DoubleHeight` | Flat-top | Row axis doubled; no parity split needed |

The offset systems require parity-aware neighbor offset tables — use `NeighborOffsets(index, system)` or the per-system `DirectionsOffset*` variables when you need neighbors in offset space.

## Line-of-sight and field-of-view

`HasLineOfSight(target, blocking)` traces a straight line from `h` to `target`. The source cell is never treated as a blocker; neither is the target cell itself — only cells strictly between them matter.

`FieldOfView(candidates, blocking)` returns the subset of `candidates` visible from `h`. Any hex within distance 1 is always considered visible regardless of blockers.

Both functions accept a `[]Hex` slice for blockers. For large grids, build a map keyed by `Hex` and pass a function instead — the slice-based API is O(n×m) per call.

## Testing helpers

`AssertHex` and `AssertFracHex` are exported from the package for use in downstream test code:

```go
hex.AssertHex(t, result, 2, -1)
hex.AssertFracHex(t, frac, 1.5, -0.5)
```

## Credits

- [Tomáš Novotný](https://github.com/tomas-novotny)
- [All Contributors][link-contributors]
- [Red Blob Games](https://www.redblobgames.com/grids/hexagons) — the canonical reference for hex grid math

## License

The MIT License (MIT). Please see [License File][link-licence] for more information.


[ico-license]:              https://img.shields.io/github/license/gravitton/hexagon.svg?style=flat-square&colorB=blue
[ico-workflow]:             https://img.shields.io/github/actions/workflow/status/gravitton/hexagon/main.yml?branch=main&style=flat-square
[ico-release]:              https://img.shields.io/github/v/release/gravitton/hexagon?style=flat-square&colorB=blue
[ico-go-dev-reference]:     https://img.shields.io/badge/go.dev-reference-blue?style=flat-square
[ico-coverage]:             https://img.shields.io/coverallsCoverage/github/gravitton/hexagon?style=flat-square

[link-author]:              https://github.com/gravitton
[link-release]:             https://github.com/gravitton/hexagon/releases
[link-contributors]:        https://github.com/gravitton/hexagon/contributors
[link-licence]:             ./LICENSE.md
[link-changelog]:           ./CHANGELOG.md
[link-workflow]:            https://github.com/gravitton/hexagon/actions
[link-go-dev-reference]:    https://pkg.go.dev/github.com/gravitton/hexagon
[link-coverage]:            https://coveralls.io/github/gravitton/hexagon
