package transport

import (
	"context"

	"github.com/domainr/epp2/schema/epp"
)

type transaction struct {
	ctx   context.Context
	reply chan reply
}

func newTransaction(ctx context.Context) (transaction, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	return transaction{
		ctx:   ctx,
		reply: make(chan reply, 1), // 1-buffered to not block
	}, cancel
}

type reply struct {
	body epp.Body
	err  error
}
