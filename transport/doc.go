/*
Package transport provides low-level EPP client and server implementations.

Opening an EPP client connection:

	tcfg := &tls.Config{ServerName: "epp.example.com"}
	tconn, err := tls.Dial("tcp", "epp.example.com:700", tcfg)
	if err != nil {
		// handle error
	}
	client := NewClient(&NetConn{Conn: tconn})
	// ...

Waiting for the initial <greeting> from the EPP server:

	greeting, err := client.Greeting(context.Background())
	if err != nil {
		// handle error
	}
	// ...

Sending a EPP <hello> and waiting for the new <greeting>:

	greeting, err := client.Hello(context.Background())
	if err != nil {
		// handle error
	}
	// ...

Processing an EPP <command> and waiting for a <response>:

	greeting, err := client.Command(&epp.Command{...})
	if err != nil {
		// handle error
	}
	// ...

*/

package transport
