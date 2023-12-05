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
	go echoServer(t, ctx, s)

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
			ctx, cancel := context.WithCancel(ctx)
			defer cancel()
			time.Sleep(randDuration(10 * time.Millisecond))
			res, err := c.ExchangeDataUnit(ctx, req)
			if err != nil {
				t.Errorf("ExchangeDataUnit(): err == %v", err)
				t.Fail()
			}
			if !bytes.Equal(req, res) {
				t.Errorf("ExchangeDataUnit(): got %s, expected %s", string(res), string(req))
				t.Fail()
			}
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
	go echoServer(t, ctx, s)

	wantErr := errors.New("test deadline exceeded")
	ctx, cancel2 := context.WithDeadlineCause(ctx, time.Now(), wantErr)
	defer cancel2()

	_, err := c.ExchangeDataUnit(ctx, []byte("hello"))
	if err != wantErr {
		t.Errorf("ExchangeDataUnit(): err == %v, expected %v", err, wantErr)
		t.Fail()
	}
}

func TestClientContextCancelled(t *testing.T) {
	wantErr := errors.New("client context canceled")
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(wantErr)

	clientConn, serverConn := net.Pipe()
	c := &Client{Conn: clientConn}
	s := &Server{Conn: serverConn}
	go echoServer(t, ctx, s)

	_, err := c.ExchangeDataUnit(ctx, []byte("hello"))
	if err != wantErr {
		t.Errorf("ExchangeDataUnit(): err == %v, expected %v", err, wantErr)
		t.Fail()
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

	err := echoServer(t, serverCtx, s)
	if err != wantErr {
		t.Errorf("echoServer(): err == %v, expected %v", err, wantErr)
	}
}

// echoServer implements a rudimentary EPP data unit server that echoes
// back each received request.
func echoServer(t *testing.T, ctx context.Context, s *Server) error {
	ctx, cancel := context.WithCancelCause(ctx)
	defer cancel(errTestDone)

	sem := make(chan struct{}, 10)
	for {
		if t.Failed() {
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

			req, w, err := s.ServeDataUnit(reqCtx)
			if err != nil {
				if err == errTestDone {
					return
				}
				t.Errorf("echoServer: ServeDataUnit(): err == %v", err)
			}
			time.Sleep(randDuration(10 * time.Millisecond))
			err = w.RespondDataUnit(reqCtx, req)
			if err != nil {
				if err == errTestDone {
					return
				}
				t.Errorf("echoServer: WriteDataUnit(): err == %v", err)
			}
		}()
	}
}

func randDuration(max time.Duration) time.Duration {
	return time.Duration(rand.Int63n(int64(max)))
}

var errTestDone = errors.New("test done")
var errTestFailed = errors.New("test failed")
