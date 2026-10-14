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

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/hero-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/hero-light.svg">
  <img alt="A hex map: the field of view of a unit among walls, a line to the farthest hex it sees, and the range of a second unit" src="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/hero-light.svg">
</picture>
</p>

Every picture in this README is drawn from the package's own results: the hexes coloured are
the ones the code beside them returns.

## Features

- **Axial coordinates**: two integers per hex, the third cube coordinate derived.
- **Immutable**: every method returns a new value, and a `Hex` is a map key that compares with `==`.
- **Traversal**: neighbors, range, ring, spiral and lines, each with an allocation-free `Seq` form.
- **Visibility**: exact, symmetric line of sight and field of view, with a single-sweep field over a whole range.
- **Cube space**: rotation by sixths of a turn, around any hex, and reflection across the q-, r- and s-axis.
- **Directions**: six of them, by increasing angle, with flat-top and pointy-top compass aliases.
- **Seven coordinate systems**: offset and doubled, converting losslessly both ways.
- **Reproducible**: fractional hexes round the same on amd64 and arm64.
- **Extras**: `geometry` interop, JSON, name parsing and `hextest` assertions.

The rules behind them are in the [package documentation][link-go-dev-reference].

## Installation

```shell
go get github.com/gravitton/hexagon
```

```go
import hex "github.com/gravitton/hexagon"
```

## Usage

### Coordinates

A hex is two integers, `q` and `r`. The third cube coordinate, `s`, is derived, so that
`q + r + s` is always zero and no value can hold an inconsistent triple.

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/coordinates-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/coordinates-light.svg">
  <img alt="A hexagon of hexes labelled with their axial coordinates, their three cube coordinates, and their distance from the center" src="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/coordinates-light.svg">
</picture>
</p>

```go
a := hex.Pt(2, -1) // axial (q, r)

a.S()      // -1, the derived cube coordinate
a.QRS()    // 2, -1, -1
a.Length() // 2, steps from the origin

a.String()             // "(2,-1)"
hex.ParseHex("(2,-1)") // the string back to the hex
```

### Arithmetic

Hexes add, subtract and scale like vectors, each coordinate on its own.

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/arithmetic-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/arithmetic-light.svg">
  <img alt="Two hexes a and b with their sum, their difference, and a multiplied by two and by minus one" src="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/arithmetic-light.svg">
</picture>
</p>

```go
a, b := hex.Pt(2, -1), hex.Pt(-1, 2)

a.Add(b)        // Hex{1, 1}
a.Subtract(b)   // Hex{3, -3}
a.Multiply(2)   // Hex{4, -2}
a.DistanceTo(b) // 3, the length of the difference
a.Lerp(b, 0.5)  // Hex{1, 0}, the hex halfway along
```

### Directions

Six directions, named for the cube axis they move along and ordered by increasing angle:
clockwise as drawn on a screen with Y pointing down. Each has a compass alias in both layouts.

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/directions-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/directions-light.svg">
  <img alt="The six directions around a hex in the pointy-top and the flat-top layout, with their compass aliases" src="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/directions-light.svg">
</picture>
</p>

```go
dir := hex.QPlus
dir.Opposite() // QMinus
dir.Turn(1)    // SMinus, one step of increasing angle
dir.Hex()      // Hex{1, -1}, the step as a hex; Offset gives it as a geom.Vector[int]
dir.Angle()    // 5π/3

hex.Directions()                // all six, by increasing angle, as a fresh array
hex.DirectionOffsets()          // the Offset of each, in the same order
hex.DirectionFromAngle(math.Pi) // SPlus, the nearest of the six
hex.ParseDirection("QPlus")     // the name back to the constant, "None" to DirectionNone
```

The angles in both panels are the ones `Angle` returns, and they follow the pointy-top layout.
In a flat-top layout every direction lies 30° further along, so take π/6 off a flat-top angle
before `DirectionFromAngle`.

### Neighbors

The neighbors come back as arrays indexed by direction, so they allocate nothing.
`DirectionTo` splits the plane into the six wedges the directions point into.

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/neighbors-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/neighbors-light.svg">
  <img alt="The six neighbors and six diagonal neighbors of a hex, numbered by direction, and every hex around it coloured by the direction toward it" src="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/neighbors-light.svg">
</picture>
</p>

```go
h := hex.Pt(0, 0)

h.Neighbors()                // the six adjacent hexes, indexed by Direction
h.Neighbor(hex.QPlus)        // Hex{1, -1}
h.DiagonalNeighbors()        // the six hexes two steps away, each between two neighbors
h.DirectionTo(hex.Pt(4, -1)) // SMinus, the direction nearest to another hex
```

A hex exactly between two directions, as a diagonal neighbor is, takes the one the other is a sixth-turn past.

### Range, ring, spiral and line

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/traversal-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/traversal-light.svg">
  <img alt="A filled range of radius 2, the ring of radius 2 numbered in order, the spiral of radius 2 numbered from the center outward, and a line between two hexes" src="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/traversal-light.svg">
</picture>
</p>

```go
h.Range(2)  // all 19 hexes within radius 2
h.Ring(2)   // the 12 hexes at exactly radius 2, in order around
h.Spiral(2) // the same hexes as Range, ordered center-outward

a, b := hex.Pt(-3, 1), hex.Pt(3, -2)
a.Line(b) // the 7 hexes from a to b in a straight line
```

`Ring(1)` is exactly `Neighbors()`. A negative radius gives `nil`, a zero radius the center alone.

Each of them has a `Seq` form that walks the hexes without a buffer:

```go
for near := range h.SpiralSeq(5) { // nearest first, stops at the break
	if occupied[near] {
		break
	}
}
```

To reuse a buffer, collect the `Seq` form with the standard library's `slices.AppendSeq`: it
works with every one of them and allocates nothing once the buffer has room.

```go
buffer := make([]hex.Hex, 0, hex.RangeLen(2)) // 19 hexes within radius 2; RingLen sizes a ring
buffer = slices.AppendSeq(buffer[:0], h.RangeSeq(2))
```

### Line of sight

Sight runs along the exact segment between two hex centers, in integers. A blocker the segment
cuts through blocks it, however small the corner; one it only touches at a corner does not; and
where the segment runs along the edge between two hexes, it takes both of them to block.

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/sight-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/sight-light.svg">
  <img alt="Three cases of a line of sight: through a blocking hex is hidden, along an edge beside one blocking hex is seen, along an edge between two blocking hexes is hidden" src="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/sight-light.svg">
</picture>
</p>

```go
a, b := hex.Pt(0, 0), hex.Pt(2, 2)

a.HasLineOfSight(b, []hex.Hex{hex.Pt(1, 1)})               // false, through it
a.HasLineOfSight(b, []hex.Hex{hex.Pt(1, 0)})               // true, one beside the edge
a.HasLineOfSight(b, []hex.Hex{hex.Pt(1, 0), hex.Pt(0, 1)}) // false, both beside the edge

a.HasLineOfSightFunc(b, blocked) // the same rule, asking blocked(hex) along the segment
```

The rule reads the same from either end, so two hexes always see each other or neither does,
and neither end ever blocks. It is stricter than `Line`, which picks one hex per step.

### Field of view

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/field-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/field-light.svg">
  <img alt="The field of view of a hex among walls: lines to seven candidates, seen or hidden, and every hex within radius 5, seen or hidden" src="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/field-light.svg">
</picture>
</p>

```go
a.FieldOfView(candidates, blocking)    // the candidates a has line of sight to; also FieldOfViewSeq
a.FieldOfViewFunc(candidates, blocked) // also FieldOfViewFuncSeq
a.FieldOfViewRange(5, blocked)         // every visible hex within radius 5; also FieldOfViewRangeSeq
```

The slice forms cost one pass over the blockers per target, whatever the distance, and are the
faster ones on a small board or with few blockers. The `Func` forms take a
`blocked func(hex.Hex) bool`, a lookup in a map or a grid, and ask it only about the hexes the
segment meets, so they cost the distance however many hexes block: the forms for a large or
dense map. Both give the same answer.

`FieldOfViewRange` is the field over a whole range: what a unit sees, or a blast reaches. It
gives the hexes of `Spiral(radius)` that `FieldOfViewFunc` keeps, but sweeps the rings outward
once and asks `blocked` about each hex at most once, where a line to every hex walks the near
ones again and again. At a radius of 40 it is about fifteen times faster; use `FieldOfViewFunc`
for a handful of candidates.

### Rotation and reflection

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/transform-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/transform-light.svg">
  <img alt="A hex turned around the origin in six steps, the same hex turned around another hex, and the hex reflected across the q-, r- and s-axis" src="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/transform-light.svg">
</picture>
</p>

```go
a, b := hex.Pt(3, -1), hex.Pt(1, 0)

a.Turn(1)          // Hex{1, 2}, 60° around the origin
a.TurnAround(b, 1) // Hex{2, 1}, 60° around b; negative steps turn the other way
a.ReflectQ()       // Hex{3, -2}, q kept, r and s swapped
a.ReflectR()       // Hex{-2, -1}, r kept, q and s swapped
a.ReflectS()       // Hex{-1, 3}, s kept, q and r swapped
```

### Coordinate systems

Axial coordinates are the ones to compute with; the offset and doubled systems are the ones a
rectangular map is stored in. Every system converts both ways without loss.

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/systems-pointy-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/systems-pointy-light.svg">
  <img alt="The same patch of pointy-top hexes indexed in the Axial, OffsetOddR, OffsetEvenR and DoubleWidth coordinate systems" src="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/systems-pointy-light.svg">
</picture>
</p>

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/systems-flat-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/systems-flat-light.svg">
  <img alt="The same patch of flat-top hexes indexed in the Axial, OffsetOddQ, OffsetEvenQ and DoubleHeight coordinate systems" src="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/systems-flat-light.svg">
</picture>
</p>

```go
a := hex.Pt(2, -1)
a.To(hex.OffsetOddR)                // geom.Point[int]{1, -1}; or hex.OffsetOddR.To(a)
hex.OffsetOddR.From(geom.Pt(1, -1)) // back to Hex{2, -1}

hex.CoordinateSystems()            // all seven, as a fresh array
hex.ParseCoordinateSystem("Axial") // the name back to the constant
```

In an offset system the step to a neighbor depends on the parity of the row or column. Ask the
system for it:

```go
hex.OffsetOddR.Offsets(index)                     // all 6, indexed by Direction
hex.OffsetOddR.Offset(index, hex.PointyTopEast)   // just one
hex.OffsetOddR.Neighbor(index, hex.PointyTopEast) // the index that offset steps to
hex.OffsetOddR.Neighbors(index)                   // all 6 indices beside it, indexed by Direction
```

`Direction.Offset()` is the axial case of the same thing, equal to `hex.Axial.Offset(index, direction)`.

### Fractional hexes

`FractionalHex` is the float companion, for interpolation and for landing a position from pixel
space on the grid. `Round` gives the hex a point lies in, keeping `q + r + s = 0`:

<p align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/rounding-dark.svg">
  <source media="(prefers-color-scheme: light)" srcset="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/rounding-light.svg">
  <img alt="Points across a patch of hexes, each coloured by the hex it rounds to, and one fractional hex with the center of the hex it rounds to" src="https://raw.githubusercontent.com/gravitton/hexagon/refs/heads/main/docs/images/rounding-light.svg">
</picture>
</p>

```go
f := hex.FracPt(1.4, -1.8)
f.Round()                                     // Hex{1, -2}
f.Length()                                    // 1.8, without rounding to a hex
hex.FracPt(0, 0).Lerp(hex.FracPt(3, -1), 0.5) // FractionalHex{1.5, -0.5}

hex.Pt(2, -1).Float()                  // the other way, exactly
hex.ParseFractionalHex("(1.40,-1.80)") // the form String prints, two decimals

f.Turn(1)                          // a sub-hex position turns and reflects like a Hex
f.TurnAround(hex.FracPt(1, 0), -1) // also ReflectQ, ReflectR and ReflectS
```

It also has `Add`, `Subtract`, `Multiply`, `DistanceTo`, `Equal` and `IsZero`, with the tolerance
of `geom.Epsilon` on the float comparisons, and `IsFinite`, the test for a position `Round` can
take: it panics on a NaN or an infinity.

### Interop

Pixel layout lives in [`gravitton/grid`](https://github.com/gravitton/grid), the shapes and
vectors this package builds on in [`gravitton/geometry`](https://github.com/gravitton/geometry).

```go
a.Point()         // geom.Point[int]{2, -1}
a.Vector()        // geom.Vector[int]{2, -1}, a hex as a step
a.Float().Point() // geom.Point[float64]{2, -1}; Vector gives it as a geom.Vector[float64]

json.Marshal(hex.Pt(2, -1))  // {"q":2,"r":-1}
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
everything it leaves out: pixel layout, storage of a value per hex, pathfinding, the hexagon as a
drawn shape, and sets of cells with union, intersection and difference, as `grid.Region`, over
the indices square and hexagonal grids share.

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

[link-release]:             https://github.com/gravitton/hexagon/releases
[link-contributors]:        https://github.com/gravitton/hexagon/contributors
[link-licence]:             ./LICENSE.md
[link-workflow]:            https://github.com/gravitton/hexagon/actions
[link-go-dev-reference]:    https://pkg.go.dev/github.com/gravitton/hexagon
[link-coverage]:            https://coveralls.io/github/gravitton/hexagon
