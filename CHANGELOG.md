# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project follows [Semantic Versioning](http://semver.org/spec/v2.0.0.html) with one deliberate exception:

**a breaking change may ship in a minor release.** Breaking changes are avoided where the cost is reasonable, but a
Go major version changes the import path for every user, and for a package with a small user base that is often the
more expensive of the two. A breaking change is marked **breaking** in its entry, and every release that has one
lists them under **Breaking** at the top of its section. Renames land as a rename; no deprecated alias is kept.


## [Unreleased](https://github.com/gravitton/hexagon/compare/v1.4.0...main)

### Breaking
- **breaking** `Direction.Rotate`, `Hex.Rotate` and `Hex.RotateAround` are `Turn` and `TurnAround`, following `geom.Direction.Turn`: a count of steps is a turn, and `Rotate` is kept for an angle in radians
- **breaking** `Directions` is the function `Directions()` returning a fresh array, matching `geom.Directions()`, so a caller cannot alter the package's own list
- **breaking** `Direction(-1)` is `DirectionNone` rather than `QPlus`; every other out-of-range value still wraps into `[SMinus, QPlus]`
- **breaking** `CoordinateSystem.String` prints `"None"` for a value outside the seven, where it printed `CoordinateSystem(99)`
- **breaking** `Hex.Neighbors` returns a `[6]Hex` array rather than a slice, so it allocates nothing
- **breaking** `Hex` and `FractionalHex` have JSON tags, so a hex is stored as `{"q":1,"r":-2}` with the one-character keys `geom` uses
- **breaking** The assertions moved to the new `hextest` package as `hextest.AssertHex` and `hextest.AssertFractionalHex`, so the main package no longer imports a test library itself. They take the expected value rather than loose `q, r` coordinates, take `assert.Testing` rather than `*testing.T`, return whether the assertion held and prefix their messages per field, like every `geomtest.Assert*` helper. `AssertFracHex` is renamed `AssertFractionalHex` and compares through `geomtest.AssertNumber`, within `geom.EpsilonRelative`, so the tolerance holds at any magnitude

### Added
- `DirectionNone` and `Direction.IsNone` – the absence of a direction, matching `geom.Direction`: it steps nowhere (`Offset` and `Hex` give a zero step, `Hex.Neighbor` returns the hex itself, `CoordinateSystem.Offset` gives the zero vector), turns to itself and has no angle
- `Direction.Angle` – the angle of a direction in `[0, 2π)`, in the order of the constants, and NaN for `DirectionNone`
- `DirectionFromAngle(angle)` – the direction nearest to an angle, `DirectionNone` for NaN and ±Inf, round-tripping with `Angle` for every direction
- `ParseDirection(name)`, `Direction.MarshalText` and `UnmarshalText` – a direction is stored as `"QPlus"` in JSON and as a map key rather than as its number, like `geom.Direction`; `"None"` parses to `DirectionNone`
- `CoordinateSystemNone`, `CoordinateSystem.IsNone`, `ParseCoordinateSystem`, `MarshalText` and `UnmarshalText` – the same for the seven systems, like `geom.Orientation`
- `CoordinateSystems()` – all seven systems in order, as a fresh array
- `Hex.AppendRange`, `AppendRing`, `AppendSpiral`, `AppendLine` and `AppendFieldOfView` – the `Append` form of every slice result, appending to a buffer the caller reuses as `strconv.AppendInt` does, so a loop allocates nothing once the buffer has room; the slice methods are the `Append` form on a buffer of the exact size
- `Hex.Equal` and `FractionalHex.Equal` – exact for `Hex`, within `geom.Epsilon` for `FractionalHex`
- `Hex.Lerp(hex, t)` – the hex at the interpolated position, rounded through `FractionalHex`
- `Hex.Float` – the `FractionalHex` of a hex, the counterpart of `FractionalHex.Round`
- `FractionalHex.Length`, `DistanceTo`, `Add`, `Subtract`, `Multiply` and `IsZero` – the arithmetic `Hex` has, without rounding to a hex

### Changed
- Require `geometry` v1.15.0 and `assert` v1.6.0
- `Hex.Spiral` allocates once rather than once per ring, and `HasLineOfSight` walks the line without building it, so a line-of-sight test allocates nothing and `FieldOfView` allocates only its result
- `Hex.Neighbor` steps by `Direction.Hex`, so a direction has one lattice step and the two cannot disagree
- A `CoordinateSystem` outside the seven panics with `hex: unknown coordinate system` and the value, as `geom` panics on an unknown `Orientation`

### Fixed
- `FractionalHex.Lerp`, and through it `Hex.Lerp` and `Hex.Line`, give the same bits on every architecture: `geom.Lerp` in geometry v1.15.0 rounds its product before adding it, so arm64 and amd64 v3 no longer fuse it into a multiply-add, which could pick a different hex there

### Removed
- The indirect `github.com/gravitton/x` dependency, which geometry v1.15.0 dropped


## [v1.4.0 (2026-08-26)](https://github.com/gravitton/hexagon/compare/v1.3.0...v1.4.0)
### Added
- `Hex.Compare(hex) int` – orders hexes by Q and then by R in the `cmp.Compare` convention, the axial counterpart of `geom.Point.Compare`, for `slices.SortFunc` and `slices.BinarySearchFunc`. The order is total but arbitrary in space — it follows neither distance nor angle from the origin


## [v1.3.0 (2026-08-26)](https://github.com/gravitton/hexagon/compare/v1.2.0...v1.3.0)
### Changed
- Require Go 1.27
- `Direction` constants are renumbered to run by increasing angle — `SMinus`, `RPlus`, `QMinus`, `SPlus`, `RMinus`, `QPlus` — matching `geometry` v1.11.0, so a positive `Rotate` step is counterclockwise in math coordinates and clockwise as drawn on a screen with Y pointing down. Everything ordered by direction follows the new order: `Directions`, the offset tables, `CoordinateSystem.Offsets`, `Hex.Neighbors`, `Hex.Ring`, and `Hex.Spiral` (**breaking**)
- `Hex.Rotate` and `Hex.RotateAround` now turn in the same sense as `Direction.Rotate`, as they always did before the reorder — no change for callers, but the two are now guaranteed consistent and tested as such
- `Hex.Ring` starts from the `SMinus` corner (angle 0) instead of the `QMinus` corner, making `Ring(1)` exactly `Neighbors()`; `Hex.Spiral` follows (**breaking**)
- `Direction.NeighborOffset()` renamed to `Direction.Offset()`, matching `geom.Direction.Offset` and freeing the name from the package-level `NeighborOffset` function (**breaking**)
- Compass direction aliases spelled out in full, matching `geometry`: `FlatTopSE` is now `FlatTopSouthEast`, `PointyTopNE` is now `PointyTopNorthEast`, and so on for all twelve (**breaking**)
- `To(hex, system)` and `From(index, system)` replaced by the methods `CoordinateSystem.To(hex)` and `CoordinateSystem.From(index)`; `Hex.To(system)` is unchanged and now forwards to `system.To(hex)` (**breaking**)
- `NeighborOffsets(index, system)` and `NeighborOffset(index, system, direction)` replaced by the methods `CoordinateSystem.Offsets(index)` and `CoordinateSystem.Offset(index, direction)`, matching the method style `geometry` v1.10.0 moved to (**breaking**)
- `NeighborOffsetsAxial`, `NeighborOffsetsOffsetOddR`, `NeighborOffsetsOffsetEvenR`, `NeighborOffsetsOffsetOddQ`, `NeighborOffsetsOffsetEvenQ`, `NeighborOffsetsDoubleWidth`, and `NeighborOffsetsDoubleHeight` removed — `CoordinateSystem.Offsets` covers all seven (**breaking**)
- `DirectionsOffsetOddR`, `DirectionsOffsetEvenR`, `DirectionsOffsetOddQ`, `DirectionsOffsetEvenQ`, `DirectionsDoubleWidth`, and `DirectionsDoubleHeight` unexported, now reached through `CoordinateSystem.Offsets` (**breaking**)
- `ToAxial`, `ToOffsetOddR`, `ToOffsetEvenR`, `ToOffsetOddQ`, `ToOffsetEvenQ`, `ToDoubleWidth`, `ToDoubleHeight` and their seven `From` counterparts unexported; `CoordinateSystem.To` and `CoordinateSystem.From` are the entry points (**breaking**)
- `Directions` now lists the six `Direction` values rather than their offset vectors, matching `geom.Directions`; the vectors moved to the unexported `directionOffsets`, reachable via `Direction.Offset` (**breaking**)
- `mod6` replaced by `geom.Mod` from `geometry` v1.10.0

### Added
- `Direction.Hex() Hex` — the unit hex step in a direction, the `Hex` counterpart of `Direction.Offset`

### Fixed
- `Direction.Offset`, `Hex.Neighbor`, and `Direction.String` no longer panic or misreport for negative directions — all direction inputs now wrap into `[SMinus, QPlus]`
- `CoordinateSystem.Offset` wraps out-of-range directions instead of panicking, as the old `NeighborOffset` free function did not
- `Hex.HasLineOfSight` no longer treats the source hex as a blocker, as the documentation already promised
- `Hex.Length` computes `(|q|+|r|+|s|)/2` in integer arithmetic instead of round-tripping through `float64`
- Neighbor offsets are returned as `[6]ints.Vector` by value; the removed functions returned slices aliasing the package's own tables, so callers could mutate them


## [v1.2.0 (2026-05-03)](https://github.com/gravitton/hexagon/compare/v1.1.0...v1.2.0)
### Breaking Changes
- `H(q, r int)` constructor renamed to `Pt(q, r int)`
- `F(q, r float64)` constructor renamed to `FracPt(q, r float64)`
- `Layout` type removed from the package (moved to [`gravitton/grid`](https://github.com/gravitton/grid))

### Added
- `Hex.IsZero() bool` — reports whether the hex equals the zero value `{0, 0}`
- `CoordinateSystem.String() string` — human-readable name for each coordinate system constant
- `Direction.String() string` — human-readable name for each direction constant
- Named direction aliases for flat-top grids: `FlatTopSE`, `FlatTopNE`, `FlatTopN`, `FlatTopNW`, `FlatTopSW`, `FlatTopS`
- Named direction aliases for pointy-top grids: `PointyTopE`, `PointyTopNE`, `PointyTopNW`, `PointyTopW`, `PointyTopSW`, `PointyTopSE`


## [v1.1.0 (2026-04-25)](https://github.com/gravitton/hexagon/compare/v1.0.0...v1.1.0)
### Added
- `Hex.Ring(radius)` — hexes at exactly the given distance from a center hex, ordered counterclockwise
- `Hex.Spiral(radius)` — all hexes from center outward to radius, ring by ring (nearest-first traversal)
- `Hex.Rotate(steps)` — rotate a hex around the origin by `steps×60°` (clockwise on screen; negative = counterclockwise)
- `Hex.RotateAround(center, steps)` — same as `Rotate` but around an arbitrary center hex
- `Hex.ReflectQ()`, `Hex.ReflectR()`, `Hex.ReflectS()` — mirror a hex across the q-, r-, or s-axis in cube space
- `Direction.Opposite()` — returns the direction directly opposite (3 steps away)
- `Direction.Rotate(steps)` — advances a direction by `steps` positions, with wrap-around and negative support


## v1.0.0 (2025-10-27)
### Added
- Support for multiple hex coordinate systems
- Support direction offsets
- Support for multiple hex grid layouts that maps hexes to pixels (and back)
