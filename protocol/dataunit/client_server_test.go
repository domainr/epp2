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

func TestEchoClientAndServer(t *testing.T) {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(errTestDone)

	clientConn, serverConn := net.Pipe()
	c := &Client{Conn: clientConn}
	s := &Server{Conn: serverConn}

	var mu sync.Mutex
	go echoServer(t, ctx, s, &mu)

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
			t.Logf("data received from server connection: %s\n", res)
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
	go echoServer(t, ctx, s, &mu)

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
	go echoServer(t, ctx, s, &mu)

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
	err := echoServer(t, serverCtx, s, &mu)
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
func echoServer(t *testing.T, ctx context.Context, s *Server, mu *sync.Mutex) error {
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
			t.Logf("data received from client connection: %s\n", data)
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
