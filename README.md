<div align="center" width="100%">

<a href="https://github.com/gravitton">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/logo-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/logo-light.svg">
  <img alt="Gravitton hexagon" src="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/logo-light.svg" width="300">
</picture>
</a>

[![Latest Stable Version][ico-release]][link-release]
[![Build Status][ico-workflow]][link-workflow]
[![Coverage Status][ico-coverage]][link-coverage]
[![Go Dev Reference][ico-go-dev-reference]][link-go-dev-reference]
[![Software License][ico-license]][link-licence]

Hexagonal grid math library for game development

<hr>

</div>


## Features

- **Axial coordinates** – two integers per hex, the third cube coordinate derived.
- **Immutable** – every method returns a new value.
- **Traversal** – neighbors, range, ring, spiral, lines, line of sight and field of view.
- **Cube space** – rotation by sixths of a turn and reflection across the q-, r- and s-axis.
- **Directions** – six, clockwise on screen, with flat-top and pointy-top aliases.
- **Seven coordinate systems** with lossless round-trip conversion.
- **Map keys** – a `Hex` compares exactly with `==` and orders with `Compare`.
- **Reproducible** – fractional hexes round the same on amd64 and arm64.
- **Extras** – `geometry` interop, JSON, name parsing and `hextest` assertions.

The rules behind them are in the [package documentation][link-go-dev-reference].

## Installation

```shell
go get github.com/gravitton/hexagon
```

## Usage

```go
import hex "github.com/gravitton/hexagon"
```

### Hexes

```go
a := hex.Pt(1, -2) // axial (q, r)
b := hex.Pt(0, 3)

a.S()           // 1, the derived cube coordinate
a.Length()      // 2, distance from the origin in hex steps
a.DistanceTo(b) // 5

a.Add(b)       // Hex{1, 1}
a.Subtract(b)  // Hex{1, -5}
a.Multiply(2)  // Hex{2, -4}
a.Lerp(b, 0.5) // the hex halfway along the line to b

a.String()             // "(1,-2)"
hex.ParseHex("(1,-2)") // the string back to the hex
```

### Neighbors and traversal

```go
h := hex.Pt(0, 0)

h.Neighbors()         // all 6 adjacent hexes as an array, by increasing angle
h.Neighbor(hex.QPlus) // one specific neighbor

h.Range(2)  // all hexes within radius 2 (filled disc)
h.Ring(2)   // hexes at exactly distance 2 (perimeter)
h.Spiral(2) // the same set as Range, ordered center-outward
```

`Ring(1)` is exactly `Neighbors()`. A negative radius gives `nil`, a zero radius the center alone.

Every slice result has an `Append` form that appends to a buffer you reuse, so a loop allocates
nothing once the buffer has room:

```go
buffer := make([]hex.Hex, 0, 37)
for _, unit := range units {
	buffer = unit.AppendRange(buffer[:0], 3)
	// ...
}
```

### Visibility

```go
a.Line(b)                     // the hexes connecting a to b in a straight line
a.HasLineOfSight(b, blocking) // only the hexes strictly between a and b block, no allocation
a.FieldOfView(candidates, blocking)
```

Both take a `[]Hex` of blockers and cost O(n×m) per call. For large grids, keep the blockers in a
map keyed by `Hex` and filter before calling.

### Rotation and reflection

```go
a.Turn(2)           // 2×60° around the origin
a.TurnAround(b, -1) // 1×60° the other way, around b
a.ReflectQ()        // mirrored across the q-axis; also ReflectR and ReflectS
```

### Directions

```go
dir := hex.QPlus
dir.Opposite() // QMinus
dir.Turn(1)    // SMinus, one step of increasing angle
dir.Offset()   // ints.Vector, the axial step
dir.Hex()      // Hex, the same step as a hex
dir.Angle()    // 5π/3

hex.Directions()                // all six, by increasing angle, as a fresh array
hex.DirectionFromAngle(math.Pi) // SPlus, the nearest of the six
hex.ParseDirection("QPlus")     // the name back to the constant, "None" to DirectionNone
```

Six named directions in cube space, with flat-top and pointy-top aliases:

| Index | Angle  | Constant | Flat-top alias     | Pointy-top alias     |
|-------|--------|----------|--------------------|----------------------|
| 0     | `0°`   | `SMinus` | `FlatTopSouthEast` | `PointyTopEast`      |
| 1     | `60°`  | `RPlus`  | `FlatTopSouth`     | `PointyTopSouthEast` |
| 2     | `120°` | `QMinus` | `FlatTopSouthWest` | `PointyTopSouthWest` |
| 3     | `180°` | `SPlus`  | `FlatTopNorthWest` | `PointyTopWest`      |
| 4     | `240°` | `RMinus` | `FlatTopNorth`     | `PointyTopNorthWest` |
| 5     | `300°` | `QPlus`  | `FlatTopNorthEast` | `PointyTopNorthEast` |

The angles, `Angle` and `DirectionFromAngle` follow the pointy-top layout. In a flat-top layout every
direction lies 30° further along, so take π/6 off a flat-top angle before `DirectionFromAngle`.

### Coordinate systems

```go
p := a.To(hex.OffsetOddR) // or hex.OffsetOddR.To(a)
hex.OffsetOddR.From(p)    // round-trips exactly

hex.CoordinateSystems()            // all seven, as a fresh array
hex.ParseCoordinateSystem("Axial") // the name back to the constant
```

| Constant       | Orientation | Description                                 |
|----------------|-------------|---------------------------------------------|
| `Axial`        | —           | Native storage: `(q, r)` directly           |
| `OffsetOddR`   | Pointy-top  | Odd rows shifted right                      |
| `OffsetEvenR`  | Pointy-top  | Even rows shifted right                     |
| `OffsetOddQ`   | Flat-top    | Odd columns shifted down                    |
| `OffsetEvenQ`  | Flat-top    | Even columns shifted down                   |
| `DoubleWidth`  | Pointy-top  | Column axis doubled; no parity split needed |
| `DoubleHeight` | Flat-top    | Row axis doubled; no parity split needed    |

The offset systems need parity-aware neighbor offsets. Ask the system for them — it handles the
row and column parity:

```go
hex.OffsetOddR.Offsets(index)                   // all 6, indexed by Direction
hex.OffsetOddR.Offset(index, hex.PointyTopEast) // just one
```

`Direction.Offset()` is the axial case of the same thing, equal to `hex.Axial.Offset(index, direction)`.

### Fractional hexes

`FractionalHex` is the float companion, for interpolation and for landing a position from pixel
space on the grid:

```go
f := hex.FracPt(1.4, -1.8)
f.Round()                                     // the nearest Hex, preserving q+r+s=0
f.Length()                                    // 1.8, without rounding to a hex
hex.FracPt(0, 0).Lerp(hex.FracPt(3, -1), 0.5) // FractionalHex{1.5, -0.5}

hex.Pt(2, -1).Float()                  // the other way, exactly
hex.ParseFractionalHex("(1.40,-1.80)") // the form String prints, two decimals


f.Turn(1)                          // a sub-hex position turns and reflects like a Hex
f.TurnAround(hex.FracPt(1, 0), -1) // also ReflectQ, ReflectR and ReflectS
```

It also has `Add`, `Subtract`, `Multiply`, `DistanceTo`, `Equal` and `IsZero`, with the tolerance
of `geom.Epsilon` on the float comparisons.

### Interop

Pixel layout lives in [`gravitton/grid`](https://github.com/gravitton/grid); the shapes and vectors
this package builds on live in [`gravitton/geometry`](https://github.com/gravitton/geometry).

```go
a.Point()         // ints.Point{1, -2}
a.Float().Point() // floats.Point{1, -2}

json.Marshal(hex.Pt(1, -2))  // {"q":1,"r":-2}
json.Marshal(hex.QPlus)      // "QPlus"
json.Marshal(hex.OffsetOddR) // "OffsetOddR"
```

### Testing

The `hextest` package has one assertion per type, comparing with the tolerance of the asserted
type, so the helpers stay out of the main package:

```go
import "github.com/gravitton/hexagon/hextest"

hextest.AssertHex(t, got, hex.Pt(2, -1))
hextest.AssertFractionalHex(t, got, hex.FracPt(1.5, -0.5), "after lerp: ")
```

Full reference: [pkg.go.dev][link-go-dev-reference].

## Not planned, by design

This package is the coordinate math alone: a hex knows its neighbors, its distances and its
shape on the grid, and nothing about where it is drawn or what it holds. That layer is
[`gravitton/grid`](https://github.com/gravitton/grid), which builds on this package and holds
everything it leaves out.

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
