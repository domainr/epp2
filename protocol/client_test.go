package protocol_test

import (
	"context"
	"io"
	"net"
	"testing"

	"github.com/domainr/epp2/protocol"
)

func TestClientConnectEOF(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	serverConn.Close()
	_, _, err := protocol.Connect(context.Background(), clientConn)
	if err != io.EOF {
		t.Errorf("Connect: expected io.EOF, got %v", err)
	}
}
