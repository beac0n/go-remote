# go-remote
This project aims to provide a replacement/addition for tools like knock.
It does so, by allowing the client to send a single (specially prepared) UDP packet to the server.
If the UDP packet contains the correct data, the server will execute the configured commands.

## build
```
make build_server
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

### systemd integration
Install `go-remote.service` and `go-remote.socket` file (https://unix.stackexchange.com/questions/224992/where-do-i-put-my-systemd-unit-file).
Make sure to edit the `go-remote.service` and add the base64 encoded key (use `-gen-key` to generate one).
Make sure you also have `go-remote-command-executor.service` installed.

If you want the service to run on a different port than `80`, edit `go-remote.service` and `go-remote.socket` respectively

# go-remote-command-executor
This is the companion service for go-remote

## build
```
make build_command_executor
```

## run
### flags
- `-user` the name of the user who is allowed to write to the tmpfs
- `-command-start` the command to execute when start is triggered, default: `echo "start!"`
- `-command-timeout` (unit: seconds) the timeout to wait after command-start was executed, default: `60`
- `-command-end` the command to execute after command-timeout is over, default: `echo "end!"`

### systemd integration
Install `go-remote-command-executor.service` file (https://unix.stackexchange.com/questions/224992/where-do-i-put-my-systemd-unit-file).
Make sure to edit the `go-remote-command-executor.service` and add the start and stop command.
Make sure you also have `go-remote.service` installed.
