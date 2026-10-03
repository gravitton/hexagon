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
// # Coordinates
//
// A hex is stored as the axial pair (q, r); s is -q-r and is derived on demand, so no value can
// hold an inconsistent cube triple. Directions are ordered by increasing angle from SMinus
// (flat-top SE, pointy-top E): counterclockwise in the standard math convention where Y grows
// upward, which appears clockwise as drawn on a screen with Y pointing down. Direction.Turn,
// Hex.Turn, FractionalHex.Turn and Ring all turn that way; negative steps go back.
//
// # Enums
//
// Every out-of-range Direction wraps into [SMinus, QPlus], negatives included. Only
// DirectionNone stands outside the six: it steps nowhere, turns to itself, and has no angle.
// CoordinateSystemNone is its counterpart for the seven systems, and every value outside them
// counts as none. Both enums have IsNone, String, a Parse function and text marshalling by name.
//
// # Equality
//
// Axial coordinates are integers, so Hex.Equal is exact and a Hex works as a map key with ==.
// Hex.Compare orders by q and then by r in the cmp.Compare convention: a total order, but an
// arbitrary one in space, following neither distance nor angle. FractionalHex.Equal compares
// within geom.Epsilon for float64.
//
// # Allocations
//
// Only a method whose result is a slice allocates, once, with the exact capacity of the result:
// Range, Ring, Spiral, Line and FieldOfView. Each has an Append form that appends to a buffer the
// caller reuses, as strconv.AppendInt does, so a loop allocates nothing once the buffer has room.
// Neighbors and DiagonalNeighbors return an array, and HasLineOfSight walks the line without
// building it.
//
// # Panics
//
// Degenerate inputs answer rather than panic: a negative radius gives nil from Range, Ring and
// Spiral, a zero radius gives the center alone, a line from a hex to itself is that one hex, and
// DirectionNone gives the hex itself from Neighbor and DiagonalNeighbor, and is the answer of
// DirectionTo for the hex itself.
//
// CoordinateSystem.Offsets, Offset, Neighbor, To and From panic for a system outside the seven,
// CoordinateSystemNone included, as geom panics for an unknown Orientation: the absence of a
// layout has nothing to convert. FractionalHex.Round panics for a NaN or infinite coordinate
// through geom.Cast, and Hex.Lerp with it for a NaN or infinite t: a non-finite position lies
// on no hex. These are the only panics.
//
// # Reproducibility
//
// Every product in FractionalHex arithmetic is rounded before it is added, so the compiler
// cannot fuse it into a multiply-add on arm64 or amd64 v3, and the same inputs give the same
// bits, and the same rounded hex, on every architecture.
package hex
