package dataunit

import (
	"bytes"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	clientConn, serverConn := Pipe()

	c := NewClient(clientConn)

	s := NewServer(serverConn)
	go echoServer(t, s)

	sem := make(chan struct{}, 100)

	for i := 0; i < 1000; i++ {
		i := i
		sem <- struct{}{}
		go func() {
			time.Sleep(randDuration(10 * time.Millisecond))
			req := []byte(strconv.FormatInt(int64(i), 10))
			res, err := c.SendDataUnit(req)
			if err != nil {
				t.Errorf("SendDataUnit(): err == %v", err)
			}
			if !bytes.Equal(req, res) {
				t.Errorf("SendDataUnit(): got %s, expected %s", string(req), string(res))
			}
			<-sem
		}()
	}
}

// echoServer implements a rudimentary EPP data unit server that echoes
// back each received request.
func echoServer(t *testing.T, s Server) {
	for {
		if t.Failed() {
			return
		}
		req, w, err := s.ReceiveDataUnit()
		if err != nil {
			t.Errorf("echoServer: Next(): err == %v", err)
			return
		}
		time.Sleep(randDuration(10 * time.Millisecond))
		err = w.WriteDataUnit(req)
		if err != nil {
			t.Errorf("echoServer: WriteDataUnit(): err == %v", err)
			return
		}
	}
}

func randDuration(max time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(max)))
}
