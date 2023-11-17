package wire

type Writer interface {
	WriteDataUnit([]byte) error
}

type Queue interface {
	Next() ([]byte, Writer, error)
	Close() error
}

type queue struct {
	conn Conn
}

func NewQueue(conn Conn) Queue {
	return &queue{
		conn: conn,
	}
}

func (q *queue) Close() error {
	// TODO: clean up pending transactions
	return q.conn.Close()
}

func (q *queue) Next() ([]byte, Writer, error) {
	return nil, nil, nil
}
