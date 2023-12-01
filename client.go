package epp

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/domainr/epp2/internal/config"
	"github.com/domainr/epp2/protocol/dataunit"
)

type Client interface {
	// Login(username, password, newPassword string) error
	// Logout() error
	Close() error
}

type client struct {
	conn dataunit.Conn
}

func Dial(network, addr string, opts ...Options) (Client, error) {
	return DialContext(context.Background(), network, addr, opts...)
}

func DialContext(ctx context.Context, network, addr string, opts ...Options) (Client, error) {
	var cfg config.Config
	cfg.Join(opts...)
	dialer := cfg.Dialer
	if dialer == nil {
		dialer = &net.Dialer{
			KeepAlive: cfg.KeepAlive,
		}
	}
	conn, err := dialer.DialContext(ctx, network, addr)
	if err != nil {
		return nil, err
	}
	if cfg.TLSConfig != nil {
		conn = tls.Client(conn, cfg.TLSConfig)
	}
	return NewClient(&dataunit.NetConn{Conn: conn}, &cfg), nil
}

func NewClient(conn dataunit.Conn, opts ...Options) Client {
	return &client{
		conn: conn,
	}
}

func (c *client) Close() error {
	// TODO: handle pending transactions
	return c.conn.Close()
}
