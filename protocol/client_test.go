package protocol_test

import (
	"io"
	"testing"

	"github.com/domainr/epp2/protocol"
	"github.com/domainr/epp2/protocol/dataunit"
)

func TestClientConnectEOF(t *testing.T) {
	clientConn, serverConn := dataunit.Pipe()
	serverConn.Close()
	_, _, err := protocol.Connect(clientConn)
	if err != io.EOF {
		t.Errorf("Connect: expected io.EOF, got %v", err)
		return
	}
}
