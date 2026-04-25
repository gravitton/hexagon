// Package hex provides utilities for working with hexagonal grids.
//
// It implements axial coordinates (a 2-D projection of cube/3-D coordinates) and several common
// integer offset systems (odd-r, even-r, odd-q, even-q) along with double-
// width and double-height variants. The package also contains helpers for
// neighbor lookup, distances, ranges, and converting between axial and
// offset coordinate systems.
//
// In addition, Layout converts between hex coordinates and pixel space
// for both flat-top and pointy-top hex orientations, and can generate
// basic shapes (like a regular hexagon) for rendering.
//
// See README and unit tests for more examples.
package hex
