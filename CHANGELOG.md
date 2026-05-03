# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).


## [Unreleased](https://github.com/gravitton/hexagon/compare/v1.2.0...master)


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
