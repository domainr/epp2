package dataunit

import (
	"sync"
	"testing"
)

func TestServer(t *testing.T) {
	clientConn, serverConn := Pipe()

	var wg sync.WaitGroup
	defer wg.Wait()
	wg.Add(1)
	go func() {
		clientConn.WriteDataUnit([]byte("hello"))
		res, err := clientConn.ReadDataUnit()
		if err != nil {
			t.Errorf("ReadDataUnit(): err == %v", err)
		}
		if string(res) != "world" {
			t.Errorf("ReadDataUnit(): got %s, expected %s", string(res), "world")
		}
		wg.Done()
	}()

	s := NewServer(serverConn, 1)

	req, w, err := s.Next()
	if err != nil {
		t.Errorf("Next(): err == %v", err)
	}
	if string(req) != "hello" {
		t.Errorf("Next(): got %s, expected %s", string(req), "hello")
	}
	err = w.WriteDataUnit([]byte("world"))
}
