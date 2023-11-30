package protocol

import (
	"context"
	"time"

	"github.com/domainr/epp2/internal/config"
	"github.com/domainr/epp2/schema"
)

// Options configure types in this package with specific features.
type Options = config.Options

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
