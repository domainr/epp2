package dataunit

import (
	"bytes"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	clientConn, serverConn := Pipe()

	c := NewClient(clientConn)

	s := NewServer(serverConn)
	go echoServer(t, s)

	sem := make(chan struct{}, 2)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		if t.Failed() {
			break
		}
		req := []byte(strconv.FormatInt(int64(i), 10))
		sem <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(randDuration(10 * time.Millisecond))
			res, err := c.SendDataUnit(req)
			if err != nil {
				t.Errorf("SendDataUnit(): err == %v", err)
				t.Fail()
			}
			if !bytes.Equal(req, res) {
				t.Errorf("SendDataUnit(): got %s, expected %s", string(res), string(req))
				t.Fail()
			}
			<-sem
		}()
	}
	wg.Wait()
}

// echoServer implements a rudimentary EPP data unit server that echoes
// back each received request.
func echoServer(t *testing.T, s Server) {
	sem := make(chan struct{}, 10)
	for {
		if t.Failed() {
			break
		}
		sem <- struct{}{}
		req, w, err := s.ReceiveDataUnit()
		if err != nil {
			t.Errorf("echoServer: ReceiveDataUnit(): err == %v", err)
			t.Fail()
		}
		go func() {
			defer func() { <-sem }()
			time.Sleep(randDuration(10 * time.Millisecond))
			err = w.WriteDataUnit(req)
			if err != nil {
				t.Errorf("echoServer: WriteDataUnit(): err == %v", err)
				t.Fail()
			}
		}()
	}
}

func randDuration(max time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(max)))
}
