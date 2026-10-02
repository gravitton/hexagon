// Package hex provides hexagonal grid math for games and UI.
//
// A [Hex] is an axial coordinate pair (a 2-D projection of cube coordinates), with the third
// cube coordinate derived rather than stored, and every method returns a new value.
// [FractionalHex] is its floating-point companion, for interpolation and for rounding a
// position from pixel space onto the grid.
//
// The package covers neighbor lookup and [Direction], distance, area and perimeter traversal
// (range, ring, spiral), line drawing, line of sight and field of view, rotation and
// reflection in cube space, and conversion between seven [CoordinateSystem] layouts.
//
// Pixel layout lives in github.com/gravitton/grid; the shapes and vectors this package builds
// on live in github.com/gravitton/geometry.
//
// # Allocations
//
// Only a method whose result is a slice allocates, once, with the exact capacity of the result:
// Range, Ring, Spiral, Line and FieldOfView. Each has an Append form that appends to a buffer the
// caller reuses, as strconv.AppendInt does, so a loop allocates nothing once the buffer has room.
// Neighbors returns an array, and HasLineOfSight walks the line without building it.
//
// # Panics
//
// CoordinateSystem.Offsets, Offset, To and From panic for a system outside the seven,
// CoordinateSystemNone included, as geom panics for an unknown Orientation: the absence of a
// layout has nothing to convert. These are the only panics: every other guard returns a value
// the type can express, such as DirectionNone, nil from Range for a negative radius, or the hex
// itself from Neighbor for DirectionNone.
//
// # Reproducibility
//
// Every product in FractionalHex arithmetic is rounded before it is added, so the compiler
// cannot fuse it into a multiply-add on arm64 or amd64 v3, and the same inputs give the same
// bits, and the same rounded hex, on every architecture.
package hex
