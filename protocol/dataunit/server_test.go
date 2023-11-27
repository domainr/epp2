package dataunit

import (
	"bytes"
	"math/rand"
	"testing"
)

func TestServer(t *testing.T) {
	clientConn, serverConn := Pipe()

	s := NewServer(serverConn)
	go echoServer(t, s)

	const str = "nomagicnumbersupmysleeverightnow"
	for i := 0; i < 1000; i++ {
		a := rand.Intn(len(str) - 1)
		b := rand.Intn(len(str))
		req := []byte(str[min(a, b):max(a, b)])
		testRequest(t, clientConn, req, req)
	}
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
