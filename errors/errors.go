package errors

// String represents a simple error string.
type String string

// Error implements the [error] interface.
func (s String) Error() string {
	return "epp: " + string(s)
}

// ClosedConnection indicates a read or write operation on a closed connection.
const ClosedConnection String = "operation on closed connection"
