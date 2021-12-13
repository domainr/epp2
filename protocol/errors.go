package protocol

import (
	"errors"
	"fmt"
)

// ErrClosedConnection indicates a read or write operation on a closed connection.
var ErrClosedConnection = errors.New("epp: operation on closed connection")

// TransactionIDError indicates an invalid transaction ID.
type TransactionIDError string

// Error implements the error interface.
func (err TransactionIDError) Error() string {
	return fmt.Sprintf("epp: invalid transaction ID: %q", string(err))
}
