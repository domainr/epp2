package config

import (
	"context"
	"crypto/tls"
	"net"
	"time"

	"github.com/domainr/epp2/internal"
	"github.com/domainr/epp2/schema"
)

type Options interface {
	EPPOptions(internal.Internal)
}

// Config is an optimized form of EPP options,
// suitable for passing via the call stack.
type Config struct {
	Context   context.Context
	KeepAlive time.Duration
	Timeout   time.Duration
	Dialer    ContextDialer
	TLSConfig *tls.Config
	Pipeline  int
	Schemas   schema.Schemas
}

func (*Config) EPPOptions(internal.Internal) {}

func (cfg *Config) Join(opts ...Options) {
	for _, src := range opts {
		switch src := src.(type) {
		case nil:
			continue
		case Context:
			cfg.Context = src.Context
		case KeepAlive:
			cfg.KeepAlive = time.Duration(src)
		case Timeout:
			cfg.Timeout = time.Duration(src)
		case Dialer:
			cfg.Dialer = src.ContextDialer
		case *TLSConfig:
			cfg.TLSConfig = (*tls.Config)(src)
		case Pipeline:
			cfg.Pipeline = int(src)
		case Schemas:
			cfg.Schemas = append(cfg.Schemas, schema.Schemas(src)...)
		}
	}
}

type Setter1[T any] func(T) Options
type SetterN[T any] func(...T) Options
type Setter[T any] interface {
	Setter1[T] | SetterN[T]
}

func GetOption[T any, F Setter[T]](opts Options, setter F) (T, bool) {
	// TODO
	var zero T
	return zero, false
}

type (
	Context   struct{ context.Context } // epp.WithContext
	KeepAlive time.Duration             // epp.WithKeepAlive
	Timeout   time.Duration             // epp.WithTimeout
	Dialer    struct{ ContextDialer }   // epp.WithDialer
	TLSConfig tls.Config                // epp.WithTLS
	Pipeline  int                       // epp.WithPipeline
	Schemas   schema.Schemas            // epp.WithSchema
)

func (Context) EPPOptions(internal.Internal)    {}
func (KeepAlive) EPPOptions(internal.Internal)  {}
func (Timeout) EPPOptions(internal.Internal)    {}
func (Dialer) EPPOptions(internal.Internal)     {}
func (*TLSConfig) EPPOptions(internal.Internal) {}
func (Pipeline) EPPOptions(internal.Internal)   {}
func (Schemas) EPPOptions(internal.Internal)    {}

// ContextDialer is any type with a DialContext method that returns ([net.Conn], [error]).
type ContextDialer interface {
	DialContext(ctx context.Context, network, addr string) (net.Conn, error)
}
