# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).


## [Unreleased](https://github.com/gravitton/hexagon/compare/v1.3.0...master)


## [v1.3.0 (2026-08-26)](https://github.com/gravitton/hexagon/compare/v1.2.0...v1.3.0)
### Changed
- Require Go 1.27
- `Direction` constants are renumbered to run by increasing angle — `SMinus`, `RPlus`, `QMinus`, `SPlus`, `RMinus`, `QPlus` — matching `geom.Direction`, so a positive `Rotate` step is counterclockwise in math coordinates and clockwise as drawn on a screen with Y pointing down. `Directions` and every `DirectionsOffset*` table follow the new order (**breaking**)
- `Hex.Rotate` and `Hex.RotateAround` now turn in the same sense as `Direction.Rotate`, as they always did before the reorder — no change for callers, but the two are now guaranteed consistent and tested as such
- `Hex.Ring` is ordered by increasing angle and starts from the `SMinus` corner (angle 0) rather than `QMinus`, making `Ring(1)` exactly `Neighbors()` (**breaking**)
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
- Package documentation no longer describes `Layout`, removed in v1.2.0


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
