# BTC Handshake
This project is a simple illustration of using the btcd node to perform a handshake between an incoming and outgoing peer node.

btcd was chosen as the Bitcoin node, because of its simplicity and interoperability within the Go ecosystem.

## Usage Example

### Compile project
```
go build
```

### Start the inbound (listening) node
This will be the node that accepts requests from a connecting node and responds with an ack. 
**Note:** Go flags are used to specify the node address and port and the direction of inbound denotes that an inbound or listening node will be started.

```
# ./btc-handshake -direction=inbound -address=127.0.0.1 -port=18555
starting inbound peer on 127.0.0.1:18556!
```

### Start the outbound (connecting) node
This will be the node that connects to the listening node and waits for an ack in order to complete the handshake.
**Note:** Go flags are used to specify the node address and port and the direction of outbound denotes that an outbound or connecting node will be started.

```
# ./btc-handshake -direction=outbound -address=127.0.0.1 -port=18555
outbound version received: /btcwire:0.5.0/inbound:1.0.0/
handshake ack received!
```