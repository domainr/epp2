package dataunit

import (
	"bytes"
	"testing"
)

func TestServer(t *testing.T) {
	clientConn, serverConn := Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	s := NewServer(serverConn, 1)
	go echoServer(t, s)

	echoRequest(t, clientConn, []byte("hello"))
}

// echoRequest sends one request to an echo server, and validates the response is correct.
func echoRequest(t *testing.T, conn Conn, data []byte) {
	err := conn.WriteDataUnit(data)
	if err != nil {
		t.Errorf("WriteDataUnit(): err == %v", err)
	}
	res, err := conn.ReadDataUnit()
	if err != nil {
		t.Errorf("ReadDataUnit(): err == %v", err)
	}
	if !bytes.Equal(res, data) {
		t.Errorf("ReadDataUnit(): got %s, expected %s", string(res), string(data))
	}
}

// echoServer implements a rudimentary EPP data unit server that echoes
// back each received request.
func echoServer(t *testing.T, s Server) {
	for {
		if t.Failed() {
			return
		}
		req, w, err := s.Next()
		if err != nil {
			t.Errorf("echoServer: Next(): err == %v", err)
			return
		}
		err = w.WriteDataUnit(req)
		if err != nil {
			t.Errorf("echoServer: WriteDataUnit(): err == %v", err)
			return
		}
	}
}
