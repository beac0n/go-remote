# go-remote
This project aims to provide a replacement/addition for tools like knock.
It does so, by allowing the client to send a single (specially prepared) UDP packet to the server.
If the UDP packet contains the correct data, the server will execute the configured commands.

## build
```
make build
```

## run
application can be run in either client mode or server mode (`-server`)

### client mode

#### flags
- `-gen-key` generate key file
- `-address` udp address with port of remote upd server (e.g. `127.0.0.1:8080`)
- `-key` path to key file or base64 encoded key (e.g. `./1599914182660964612.key`)

### server mode

#### flags
- `-server` needed to indicate server mode
- `-port` port on which to run udp server, default: `8080`
- `-timeframe` (unit: seconds) timestamp sent by the client must be between `now - timeframe` and `now`, default: `5`
- `-key` path to key file or base64 encoded key (e.g. `./1599914182660964612.key`)
- `-tmpfs` path to tmpfs directory, read by `go-remote-command-executor`
