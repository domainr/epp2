package dataunit

import (
	"testing"
)

func TestServer(t *testing.T) {
	clientConn, serverConn := Pipe()
	defer clientConn.Close()
	defer serverConn.Close()

	s := NewServer(serverConn, 1)
	go echoServer(t, s)

	err := clientConn.WriteDataUnit([]byte("hello"))
	if err != nil {
		t.Errorf("WriteDataUnit(): err == %v", err)
	}
	res, err := clientConn.ReadDataUnit()
	if err != nil {
		t.Errorf("ReadDataUnit(): err == %v", err)
	}
	if string(res) != "hello" {
		t.Errorf("ReadDataUnit(): got %s, expected %s", string(res), "world")
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
		println("replied")
		if err != nil {
			t.Errorf("echoServer: WriteDataUnit(): err == %v", err)
			return
		}
	}
}
