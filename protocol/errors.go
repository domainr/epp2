package protocol

import (
	"errors"
)

// ErrClosedConnection indicates a read or write operation on a closed connection.
var ErrClosedConnection = errors.New("epp: operation on closed connection")

// TransactionIDError indicates an invalid transaction ID.
type TransactionIDError struct {
	TransactionID string
}

// Error implements the error interface.
func (err TransactionIDError) Error() string {
	return "epp: invalid transaction transaction ID: " + err.TransactionID

}

// DuplicateTransactionIDError indicates a duplicate transaction ID.
type DuplicateTransactionIDError TransactionIDError

// Error implements the error interface.
func (err DuplicateTransactionIDError) Error() string {
	return "epp: duplicate transaction ID: " + err.TransactionID
}
