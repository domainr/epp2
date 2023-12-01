package epp

// Error is the interface implemented by all errors in this package.
type Error interface {
	eppError()
}

type stringError string

func (stringError) eppError() {}

func (s stringError) Error() string {
	return "epp: " + string(s)
}

// ErrClosedConnection indicates a read or write operation on a closed connection.
const ErrClosedConnection stringError = "operation on closed connection"

// ErrServerClosed indicates a [Server] has shut down or closed.
const ErrServerClosed stringError = "server closed"

// TransactionIDError indicates an invalid transaction ID.
type TransactionIDError struct {
	TransactionID string
}

func (TransactionIDError) eppError() {}

// Error implements the error interface.
func (err TransactionIDError) Error() string {
	return "epp: invalid transaction transaction ID: " + err.TransactionID
}

// DuplicateTransactionIDError indicates a duplicate transaction ID.
type DuplicateTransactionIDError TransactionIDError

func (DuplicateTransactionIDError) eppError() {}

// Error implements the error interface.
func (err DuplicateTransactionIDError) Error() string {
	return "epp: duplicate transaction ID: " + err.TransactionID
}
