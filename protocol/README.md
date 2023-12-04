# Overview

Package protocol implements low-level [EPP](https://datatracker.ietf.org/doc/rfc5730/) client and server connections. The Client and Server types in this package provide an ordered queue of EPP commands with XML serialization and deserialization. Network-related features such as timeouts, keep-alives, or cancellation are the responsibility of the caller.

## Usage

Open an EPP client connection and wait for the initial `<greeting>`:

```go
cfg := &tls.Config{ServerName: "epp.example.com"}
conn, err := tls.Dial("tcp", "epp.example.com:700", cfg)
if err != nil {
	// handle error
}
client, greeting, err := protocol.Connect(context.Background(), conn)
// ...
```
