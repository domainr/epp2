package epp

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/domainr/epp2/protocol/dataunit"
)

type Option any

func WithKeepalive(d time.Duration) Option {
	return nil
}

func WithDialer(dialContext func(ctx context.Context, network, addr string) (net.Conn, error)) Option {
	return nil
}

func WithTLS(cfg *tls.Config) Option {
	return nil
}

type Client interface {
	Login(username, password, newPassword string) error
	Logout() error
	Close() error
}

type client struct{}

func NewClient(conn dataunit.Conn, options ...Option) Client {
	return nil
}

func Dial(addr string, options ...Option) (Client, error) {
	return nil, nil
}
