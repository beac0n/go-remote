# go-remote
This project aims to provide a replacement/addition for tools like knock.
It does so, by allowing the client to send a single (specially prepared) UDP packet to the server.
If the UDP packet contains the correct data, the server will execute the configured commands.

## build
```
make build
```

## run
Application can be run in either client mode or server mode (`-server`).

### client mode
#### flags
- `-gen-key` generate base64 encoded aes key
- `-address` udp address with port of remote upd server (e.g. `127.0.0.1:8080`)
- `-key` base64 encoded key (e.g. `fdUPciUFq0nTodfSzHiImOBuqBGzSsMSx411DyPMoZ4=`)

### server mode
#### flags
- `-server` needed to indicate server mode
- `-port` port on which to run udp server, default: `8080`
- `-timeframe` (unit: seconds) timestamp sent by the client must be between `now - timeframe` and `now`, default: `5`
- `-key` base64 encoded key (e.g. `fdUPciUFq0nTodfSzHiImOBuqBGzSsMSx411DyPMoZ4=`)
- `-tmpfs` path to tmpfs directory, read by `go-remote-command-executor`

### systemd integration
Install `go-remote.service` and `go-remote.socket` file (https://unix.stackexchange.com/questions/224992/where-do-i-put-my-systemd-unit-file).
Make sure to edit the `go-remote.service` and add the base64 encoded key (use `-gen-key` to generate one).
Make sure you also have `go-remote-command-executor.service` installed (https://github.com/beac0n/go-remote-command-executor).

If you want the service to run on a different port than `80`, edit `go-remote.service` and `go-remote.socket` respectively
