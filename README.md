# go-remote

This project aims to provide a replacement/addition for tools like knock.
It does so, by allowing the client to send a single (specially prepared) UDP packet to the server.
If the UDP packet contains the correct data, the server will execute the configured commands.

## build

Make sure to use the exact same binary on client and server side.

```
make build
```

## run
application can be run in either client mode or server mode (`-server`)

### client mode

#### flags
- `-gen-key` generate server/client key files
- `-address` udp address with port of remote upd server (e.g. `127.0.0.1:8080`)
- `-key` path to client key file (e.g. `./1599914182660964612.client`)

### server mode

#### flags
- `-server` needed to indicate server mode
- `-port` port on which to run udp server, default: `8080`
- `-timeframe` (unit: seconds) timestamp sent by client must be between `now - timeframe` and `now`, default: `5`
- `-command-start` the command to execute if udp packet sent by client is valid, default: `echo "start!"`
- `-command-timeout` (unit: seconds) the timeout to wait after command-start was executed, default: `60`
- `-command-end` the command to execute after command-timeout is over, default: `echo "end!"`
- `-key` path to server key file (e.g. `./1599914182660964612.server`)
