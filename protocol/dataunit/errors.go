package dataunit

import "errors"

// ErrClosedConnection indicates a read or write operation on a closed connection.
var ErrClosedConnection = errors.New("epp: operation on closed connection")
