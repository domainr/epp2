package dataunit

import (
	"bytes"
	"context"
	"errors"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestEchoClientAndServer(t *testing.T) {
	ctx := &testContext{Context: context.Background()}
	defer ctx.TestDone()

	clientConn, serverConn := Pipe()

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

// echoServer implements a rudimentary EPP data unit server that echoes
// back each received request.
func echoServer(t *testing.T, ctx context.Context, s *Server) {
	sem := make(chan struct{}, 10)
	for {
		if t.Failed() {
			break
		}
		err := ctx.Err()
		if err != nil {
			if err != errTestDone {
				t.Error(err)
			}
			break
		}
		sem <- struct{}{}
		go func() {
			defer func() { <-sem }()

			reqCtx, cancel := context.WithCancel(ctx)
			defer cancel()

			req, w, err := s.ServeDataUnit(reqCtx)
			if err != nil {
				t.Errorf("echoServer: ServeDataUnit(): err == %v", err)
				t.Fail()
			}
			time.Sleep(randDuration(10 * time.Millisecond))
			err = w.RespondDataUnit(reqCtx, req)
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

type testContext struct {
	context.Context
	err atomic.Value
}

func (ctx *testContext) TestDone() {
	ctx.err.Store(errTestDone)
}

func (ctx *testContext) Err() error {
	err := ctx.err.Load()
	if err != nil {
		return err.(error)
	}
	return ctx.Context.Err()
}

var errTestDone = errors.New("test done")
