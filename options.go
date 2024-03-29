package epp

import (
	"context"
	"crypto/tls"
	"time"

	"github.com/domainr/epp2/internal/config"
	"github.com/domainr/epp2/schema"
)

// Options configure [TODO] with specific features.
type Options = config.Options

// ContextDialer is any type with a DialContext method that returns ([net.Conn], [error]).
type ContextDialer = config.ContextDialer

func WithDialer(d ContextDialer) Options {
	return config.Dialer{ContextDialer: d}
}

func WithTLS(cfg *tls.Config) Options {
	return (*config.TLSConfig)(cfg.Clone())
}

func WithPipeline(depth int) Options {
	return config.Pipeline(depth)
}

func WithContext(ctx context.Context) Options {
	return config.Context{Context: ctx}
}

func WithKeepAlive(d time.Duration) Options {
	return config.KeepAlive(d)
}

func WithTimeout(d time.Duration) Options {
	return config.Timeout(d)
}

func WithSchema(schemas ...schema.Schema) Options {
	return config.Schemas(schemas)
}
