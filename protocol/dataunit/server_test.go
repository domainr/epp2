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

	testRequest(t, clientConn, []byte("hello"), []byte("hello"))
}

// testRequest sends a request to an data unit server, and validates the response matches res.
func testRequest(t *testing.T, conn Conn, req []byte, res []byte) {
	err := conn.WriteDataUnit(req)
	if err != nil {
		t.Errorf("WriteDataUnit(): err == %v", err)
	}
	got, err := conn.ReadDataUnit()
	if err != nil {
		t.Errorf("ReadDataUnit(): err == %v", err)
	}
	if !bytes.Equal(got, res) {
		t.Errorf("ReadDataUnit(): got %s, expected %s", string(got), string(res))
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
