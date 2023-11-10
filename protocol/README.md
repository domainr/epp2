# Overview

Package protocol implements low-level [EPP](https://datatracker.ietf.org/doc/rfc5730/) client and server connections.

## Usage

Open an EPP client connection:

```go
tlsConfig := &tls.Config{ServerName: "epp.example.com"}
tlsConn, err := tls.Dial("tcp", "epp.example.com:700", tlsConfig)
if err != nil {
	// handle error
}
client := protocol.NewClient(&protocol.NetConn{tlsConn})
// ...
```

Wait for the initial `<greeting>` from the EPP server:

```go
greeting, err := client.Greeting(context.Background())
if err != nil {
	// handle error
}
// ...
```

Send an EPP `<hello>` and wait for the new `<greeting>`:

```go
greeting, err := client.Hello(context.Background())
if err != nil {
	// handle error
}
// ...
```

Process an EPP `<command>` and wait for a `<response>`:

```go
greeting, err := client.Command(context.Background(), &epp.Command{...})
if err != nil {
	// handle error
}
// ...
```
