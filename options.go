package epp

import (
	"crypto/tls"

	"github.com/domainr/epp2/internal/config"
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
