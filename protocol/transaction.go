package protocol

import (
	"context"

	"github.com/domainr/epp2/schema/epp"
)

type transaction struct {
	ctx    context.Context
	result chan result
}

func newTransaction(ctx context.Context) (transaction, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	return transaction{
		ctx:    ctx,
		result: make(chan result, 1), // 1-buffered to not block
	}, cancel
}

type result struct {
	body epp.Body
	err  error
}
