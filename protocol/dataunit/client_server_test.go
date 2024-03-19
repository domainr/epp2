package dataunit

import (
	"bytes"
	"context"
	"errors"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"testing"
	"time"
)

// TestEchoClientAndServer validates the EPP server is able to read the data
// provided in the client request and to respond back correctly.
//
// It does this by passing an instance of a `Server` to `echoServer()`, which
// runs in a goroutine. The echo server runs forever or until the passed in
// context is cancelled. Inside the echo server it checks for
// `testing.T.Failed()` and stops if there is a test failure. Otherwise, it
// handles 10 requests at a time, spinning up each request in a new goroutine.
// Each goroutine handles a response to the client request. Each request is
// sending a unique 'data unit' to be processed and the response is sent back to
// the `serverConn` which .
//
// We use a `net.Pipe()` to simulate the client/server request/response flows.
// i.e. a write to `clientConn` is readable via `serverConn` (and vice-versa).
//
// So c.ExchangeDataUnit() writes data to the `clientConn` which causes a copy
// of the data to be written into `serverConn` for it to read. Then writing to
// `serverConn` causes a copy of that data to be written into `clientConn` for
// it to read. Simulating a bidirectional network connection.
func TestEchoClientAndServer(t *testing.T) {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(errTestDone)

	clientConn, serverConn := net.Pipe()
	c := &Client{Conn: clientConn}
	s := &Server{Conn: serverConn}

	var mu sync.Mutex
	go echoServer(ctx, t, s, &mu)

	sem := make(chan struct{}, 2)
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		mu.Lock()
		failed := t.Failed()
		mu.Unlock()

		if failed {
			break
		}
		data := []byte(strconv.FormatInt(int64(i), 10))
		sem <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			time.Sleep(randDuration(10 * time.Millisecond))
			res, err := c.ExchangeDataUnit(ctx, data)
			if err != nil {
				mu.Lock()
				t.Errorf("ExchangeDataUnit(): err == %v", err)
				mu.Unlock()
			}
			if !bytes.Equal(data, res) {
				mu.Lock()
				t.Errorf("ExchangeDataUnit(): got %s, expected %s", string(res), string(data))
				mu.Unlock()
			}
			// t.Logf("data received from server connection: %s\n", res)
			<-sem
		}()
	}
	wg.Wait()
}

func TestClientContextDeadline(t *testing.T) {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(errTestDone)

	clientConn, serverConn := net.Pipe()
	c := &Client{Conn: clientConn}
	s := &Server{Conn: serverConn}
	var mu sync.Mutex
	go echoServer(ctx, t, s, &mu)

	wantErr := errors.New("test deadline exceeded")
	ctx, cancel2 := context.WithDeadlineCause(ctx, time.Now(), wantErr)
	defer cancel2()

	_, err := c.ExchangeDataUnit(ctx, []byte("hello"))
	if err != wantErr {
		mu.Lock()
		t.Errorf("ExchangeDataUnit(): err == %v, expected %v", err, wantErr)
		mu.Unlock()
	}
}

func TestClientContextCancelled(t *testing.T) {
	wantErr := errors.New("client context canceled")
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(wantErr)

	clientConn, serverConn := net.Pipe()
	c := &Client{Conn: clientConn}
	s := &Server{Conn: serverConn}
	var mu sync.Mutex
	go echoServer(ctx, t, s, &mu)

	_, err := c.ExchangeDataUnit(ctx, []byte("hello"))
	if err != wantErr {
		mu.Lock()
		t.Errorf("ExchangeDataUnit(): err == %v, expected %v", err, wantErr)
		mu.Unlock()
	}
}

func TestServerContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(errTestDone)

	clientConn, serverConn := net.Pipe()
	c := &Client{Conn: clientConn}
	s := &Server{Conn: serverConn}

	go c.ExchangeDataUnit(ctx, []byte("hello"))

	wantErr := errors.New("server context canceled")
	serverCtx, cancel := context.WithCancelCause(context.Background())
	cancel(wantErr)

	var mu sync.Mutex
	err := echoServer(serverCtx, t, s, &mu)
	if err != wantErr {
		mu.Lock()
		t.Errorf("echoServer(): err == %v, expected %v", err, wantErr)
		mu.Unlock()
	}
}

func TestMultipleResponseError(t *testing.T) {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(errTestDone)

	clientConn, serverConn := net.Pipe()
	c := &Client{Conn: clientConn}
	s := &Server{Conn: serverConn}

	go c.ExchangeDataUnit(ctx, []byte("hello"))

	req, r, err := s.ServeDataUnit(ctx)
	if err != nil {
		t.Errorf("ExchangeDataUnit(): err == %v", err)
	}
	err = r.RespondDataUnit(ctx, req)
	if err != nil {
		t.Errorf("RespondDataUnit(): err == %v", err)
	}
	err = r.RespondDataUnit(ctx, req)
	wantErr := MultipleResponseError{Index: 0, Count: 2}
	if err != wantErr {
		t.Errorf("RespondDataUnit(): err == %v, expected %v", err, wantErr)
	}
}

// echoServer implements a rudimentary EPP data unit server that echoes
// back each received request.
func echoServer(ctx context.Context, t *testing.T, s *Server, mu *sync.Mutex) error {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(errTestDone)

	sem := make(chan struct{}, 10)
	for {
		mu.Lock()
		failed := t.Failed()
		mu.Unlock()

		if failed {
			return errTestFailed
		}
		err := context.Cause(ctx)
		if err != nil {
			return err
		}
		sem <- struct{}{}
		go func() {
			defer func() { <-sem }()

			reqCtx, cancel := context.WithCancel(ctx)
			defer cancel()

			data, r, err := s.ServeDataUnit(reqCtx)
			if err != nil {
				if err == errTestDone {
					return
				}
				mu.Lock()
				t.Errorf("echoServer: ServeDataUnit(): err == %v", err)
				mu.Unlock()
			}
			// t.Logf("data received from client connection: %s\n", data)
			time.Sleep(randDuration(10 * time.Millisecond))
			err = r.RespondDataUnit(reqCtx, data)
			if err != nil {
				if err == errTestDone {
					return
				}
				mu.Lock()
				t.Errorf("echoServer: RespondDataUnit(): err == %v", err)
				mu.Unlock()
			}
		}()
	}
}

func randDuration(max time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(max)))
}

var (
	errTestDone   = errors.New("test done")
	errTestFailed = errors.New("test failed")
)
