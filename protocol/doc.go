/*
Package protocol provides low-level EPP client and server implementations.

Open an EPP client connection:

	tlsConfig := &tls.Config{ServerName: "epp.example.com"}
	tlsConn, err := tls.Dial("tcp", "epp.example.com:700", tlsConfig)
	if err != nil {
		// handle error
	}
	client := protocol.NewClient(&protocol.Conn{tlsConn})
	// ...

Wait for the initial <greeting> from the EPP server:

	greeting, err := client.Greeting(context.Background())
	if err != nil {
		// handle error
	}
	// ...

Send a EPP <hello> and wait for the new <greeting>:

	greeting, err := client.Hello(context.Background())
	if err != nil {
		// handle error
	}
	// ...

Process an EPP <command> and wait for a <response>:

	greeting, err := client.Command(context.Background(), &epp.Command{...})
	if err != nil {
		// handle error
	}
	// ...
*/

package protocol
