// Package hex provides immutable hexagonal grid math for games and UI.
//
// No method changes its receiver: each returns a new value, or the receiver when nothing changes.
//
// # Types
//
//   - [Hex] is a hex of the grid in axial coordinates, and [FractionalHex] a position between hex
//     centers, for interpolation and for rounding a position from pixel space onto the grid.
//   - Enums: [Direction], the six steps to a neighbor, and [CoordinateSystem], the seven ways to
//     address a grid by two integers.
//
// Hex covers neighbors, distance, the traversals [Hex.Range], [Hex.Ring], [Hex.Spiral] and
// [Hex.Line], line of sight and field of view, and turning and reflection in cube space. Package
// hextest holds test assertions. Pixel layout lives in github.com/gravitton/grid, and the shapes
// and vectors this package builds on in github.com/gravitton/geometry.
//
// # Coordinates
//
// A hex is stored as the axial pair (q, r); s is -q-r, derived on demand, so no value can hold an
// inconsistent cube triple. Directions are ordered by increasing angle from [SMinus] (flat-top
// southeast, pointy-top east): counterclockwise in math coordinates, where Y grows upward, which
// appears clockwise as drawn on a screen with Y pointing down. Every Turn and the order of
// [Directions], [Hex.Neighbors] and Ring follow it, and negative steps go the other way.
//
// # Enums
//
// Every Direction wraps into [SMinus, QPlus], negatives included. Only [DirectionNone] stands
// outside the six: it steps nowhere, turns to itself and has no angle. [CoordinateSystemNone] is
// its counterpart for the seven systems, and every value outside them counts as none. Both enums
// have IsNone, String, a Parse function and text marshalling by name.
//
// # Equality
//
// Axial coordinates are integers, so [Hex.Equal] is exact and a Hex works as a map key with ==.
// [Hex.Compare] orders by q and then by r: a total order, but an arbitrary one in space, following
// neither distance nor angle. [FractionalHex.Equal] and IsZero compare within [geom.Epsilon] for
// float64.
//
// # Allocations
//
// Only a method whose result is a slice allocates, once: Range, Ring, Spiral and Line with the
// exact capacity of the result, [Hex.FieldOfView] and [Hex.FieldOfViewFunc] with room for every
// candidate. [Hex.FieldOfViewRange] is the one exception: it allocates room for the range and, in
// every form, one buffer more, proportional to the radius, unless blocked is nil or the radius is
// under two. Each has a Seq form, an iterator that walks the same hexes in the same order without a
// buffer for them, and slices.AppendSeq collects it into a buffer the caller reuses, so a loop
// allocates nothing once the buffer has room; [RangeLen] and [RingLen] size one. Every array a
// method returns, from Directions and Neighbors to [CoordinateSystem.Offsets], is a copy the caller
// owns, and [Hex.HasLineOfSight] builds nothing.
//
// # Visibility
//
// HasLineOfSight asks whether the straight segment between two hex centers is clear, exactly, in
// integers. A blocking hex the segment runs through blocks it, however small the corner it cuts; a
// hex the segment only touches at a corner does not; and where the segment runs along the edge
// between two hexes, both must block. Neither end counts. The rule is the same from either end, so
// sight is symmetric, and FieldOfView is the candidates it holds for: the field symmetric shadow
// casting computes. Line is the traversal, one hex per step, so a blocker off the Line can still
// clip the sight line.
//
// [Hex.HasLineOfSightFunc] and FieldOfViewFunc apply the same rule to a function reporting whether
// a hex blocks, and ask it only about the hexes that can block the segment, each at most once. The
// slice forms cost a pass over the blockers for each target, and one more for each blocker the
// segment touches along an edge, and the Func forms the distance to it, so the slice forms suit a
// small board or few blockers and the Func forms a large or dense map.
//
// FieldOfViewRange is the field of every hex within a radius, the hexes of Spiral that
// FieldOfViewFunc keeps. It asks about a hex at most once, so it costs the size of the range where
// a line to every hex costs that times the radius.
//
// All of them are exact while the square of the distance between the hexes fits an int a few
// times over, far beyond any map on a 64-bit platform.
//
// # Panics
//
// Degenerate inputs answer rather than panic: a negative radius gives nil from Range, Ring, Spiral
// and FieldOfViewRange and nothing from their Seq forms, a zero radius gives the center alone, a
// line from a hex to itself is that one hex, no candidates give nil from FieldOfView and
// FieldOfViewFunc, and DirectionNone gives the hex itself from [Hex.Neighbor] and
// [Hex.DiagonalNeighbor], and is the answer of [Hex.DirectionTo] for the hex itself.
//
// The only panics are for arguments with no meaning:
//
//   - a system outside the seven, CoordinateSystemNone included, to CoordinateSystem.Offsets,
//     [CoordinateSystem.Offset], [CoordinateSystem.Neighbor], [CoordinateSystem.Neighbors],
//     [CoordinateSystem.To] and [CoordinateSystem.From], and to [Hex.To] through them: the absence
//     of a layout has nothing to convert, and the panic value is an error that errors.Is matches
//     to [ErrUnknownCoordinateSystem];
//   - a NaN or infinite coordinate to [FractionalHex.Round], and to [Hex.Lerp] through it a t that
//     is NaN or infinite or carries the position past the range of float64: a non-finite position
//     lies on no hex, and [FractionalHex.IsFinite] tests for one.
//
// A finite position beyond the range of int rounds to a platform-dependent hex, and a radius or a
// distance too large to allocate the slice for is outside the contract.
//
// # Reproducibility
//
// Every product in FractionalHex arithmetic is rounded before it is added, so no fused
// multiply-add changes the result and the same inputs give the same bits, and the same rounded
// hex, on every architecture.
package hex
