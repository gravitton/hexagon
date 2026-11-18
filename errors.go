package hex

import (
	"errors"
)

var (
	// ErrUnknownCoordinateSystem is wrapped by the panic value of a coordinate system outside the
	// seven, CoordinateSystemNone included.
	ErrUnknownCoordinateSystem = errors.New("hex: unknown coordinate system")
	// ErrNegativeChunkRadius is wrapped by the panic value of a negative radius to the chunk
	// methods.
	ErrNegativeChunkRadius = errors.New("hex: negative chunk radius")
)
