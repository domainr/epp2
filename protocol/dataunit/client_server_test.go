package dataunit

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
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
// context is cancelled. Inside the echo server it checks for any context
// cancellations and returns early if so (it will optionally record the error if
// the error is not an expected type). Otherwise, it handles 10 requests at a
// time, spinning up each request in a new goroutine. Each goroutine handles a
// response to the client request. Each request is sending a unique 'data unit'
// to be processed and the response is sent back to the `serverConn`.
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
	n := 100
	echoServerErr := make(chan error, n)

	go echoServer(ctx, s, echoServerErr)

	sem := make(chan struct{}, 2)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		data := []byte(strconv.FormatInt(int64(i), 10))
		sem <- struct{}{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			time.Sleep(randDuration(10 * time.Millisecond))
			res, err := c.ExchangeDataUnit(ctx, data)
			// t.Logf("data received from server connection: %s\n", res)
			if err != nil {
				echoServerErr <- fmt.Errorf("ExchangeDataUnit(): err == %w", err)
			}
			if !bytes.Equal(data, res) {
				echoServerErr <- fmt.Errorf("ExchangeDataUnit(): got %s, expected %s", string(res), string(data))
			}
			<-sem
		}()
	}
	wg.Wait()
	close(echoServerErr)
	for err := range echoServerErr {
		t.Error(err)
	}
}

func TestClientContextDeadline(t *testing.T) {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(errTestDone)

	clientConn, serverConn := net.Pipe()
	c := &Client{Conn: clientConn}
	s := &Server{Conn: serverConn}
	echoServerErr := make(chan error)

	go echoServer(ctx, s, echoServerErr)

	wantErr := errors.New("test deadline exceeded")
	ctx, cancel2 := context.WithDeadlineCause(ctx, time.Now(), wantErr)
	defer cancel2()

	_, err := c.ExchangeDataUnit(ctx, []byte("hello"))
	if !errors.Is(err, wantErr) {
		echoServerErr <- fmt.Errorf("ExchangeDataUnit(): err == %w, expected %w", err, wantErr)
	}
	close(echoServerErr)
	for err := range echoServerErr {
		t.Error(err)
	}
}

func TestClientContextCancelled(t *testing.T) {
	wantErr := errCtxCancelled
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(wantErr)

	clientConn, serverConn := net.Pipe()
	c := &Client{Conn: clientConn}
	s := &Server{Conn: serverConn}
	echoServerErr := make(chan error)

	go echoServer(ctx, s, echoServerErr)

	_, err := c.ExchangeDataUnit(ctx, []byte("hello"))
	if !errors.Is(err, wantErr) {
		echoServerErr <- fmt.Errorf("ExchangeDataUnit(): err == %w, expected %w", err, wantErr)
	}
	close(echoServerErr)
	for err := range echoServerErr {
		t.Error(err)
	}
}

func TestServerContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(errTestDone)

	clientConn, serverConn := net.Pipe()
	c := &Client{Conn: clientConn}
	s := &Server{Conn: serverConn}

	go func() {
		_, _ = c.ExchangeDataUnit(ctx, []byte("hello"))
	}()

	wantErr := errors.New("server context cancelled")
	serverCtx, cancel := context.WithCancelCause(context.Background())
	cancel(wantErr)

	echoServerErr := make(chan error)

	go echoServer(serverCtx, s, echoServerErr)

	for err := range echoServerErr {
		if !errors.Is(err, wantErr) {
			t.Errorf("unexpected error: %s", err)
		}
	}
}

func TestMultipleResponseError(t *testing.T) {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(errTestDone)

	clientConn, serverConn := net.Pipe()
	c := &Client{Conn: clientConn}
	s := &Server{Conn: serverConn}

	go func() {
		_, _ = c.ExchangeDataUnit(ctx, []byte("hello"))
	}()

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
	if !errors.Is(err, wantErr) {
		t.Errorf("RespondDataUnit(): err == %v, expected %v", err, wantErr)
	}
}

// echoServer implements a rudimentary EPP data unit server that echoes
// back each received request.
func echoServer(ctx context.Context, s *Server, echoServerErr chan<- error) {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(errTestDone)

	sem := make(chan struct{}, 10)
	for {
		err := context.Cause(ctx)
		if err != nil {
			if errors.Is(err, errTestDone) || errors.Is(err, errCtxCancelled) {
				return
			}
			echoServerErr <- err
			close(echoServerErr)
			return
		}
		sem <- struct{}{}
		go func() {
			defer func() { <-sem }()

			reqCtx, cancel := context.WithCancel(ctx)
			defer cancel()

			data, r, err := s.ServeDataUnit(reqCtx)
			if err != nil {
				if errors.Is(err, errTestDone) {
					return
				}
				echoServerErr <- fmt.Errorf("echoServer: ServeDataUnit(): err == %w", err)
				return
			}
			time.Sleep(randDuration(10 * time.Millisecond))
			err = r.RespondDataUnit(reqCtx, data)
			if err != nil {
				if errors.Is(err, errTestDone) {
					return
				}
				echoServerErr <- fmt.Errorf("echoServer: RespondDataUnit(): err == %w", err)
			}
		}()
	}
}

func randDuration(max time.Duration) time.Duration {
	return time.Duration(rand.Int64N(int64(max)))
}

var (
	errCtxCancelled = errors.New("context cancelled")
	errTestDone     = errors.New("test done")
)
