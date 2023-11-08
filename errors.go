package epp

import (
	"errors"
)

// ErrServerClosed indicates a [Server] has shut down or closed.
var ErrServerClosed = errors.New("epp: server closed")
