package epp

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/domainr/epp2/internal/config"
	"github.com/domainr/epp2/protocol"
	"github.com/domainr/epp2/schema/epp"
)

type Client interface {
	// Login(username, password, newPassword string) error
	// Logout() error
	Close() error
}

type client struct {
	conn     net.Conn
	client   protocol.Client
	greeting epp.Body
}

func Dial(network, addr string, opts ...Options) (Client, error) {
	var cfg config.Config
	cfg.Join(opts...)

	ctx := cfg.Context
	if ctx == nil {
		ctx = context.Background()
	}

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

	return connect(conn, cfg)
}

func Connect(conn net.Conn, opts ...Options) (Client, error) {
	var cfg config.Config
	cfg.Join(opts...)
	return connect(conn, cfg)
}

func connect(conn net.Conn, cfg config.Config) (Client, error) {
	ctx := cfg.Context
	if ctx == nil {
		ctx = context.Background()
	}

	c, greeting, err := protocol.Connect(ctx, conn)
	if err != nil {
		return nil, err
	}

	return &client{
		conn:     conn,
		client:   c,
		greeting: greeting,
	}, nil
}

func (c *client) Close() error {
	// TODO: handle pending transactions
	return c.conn.Close()
}
