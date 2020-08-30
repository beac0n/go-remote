# go-remote

This project aims to provide a replacement/addition for tools like knock.
It does so, by allowing the client to send a single (specially prepared) UDP packet to the server.
If the UDP packet contains the correct data, the server will execute the configured command.2

## build

```
make build
```

## run
application can be run in either client mode (`-client`) or server mode (`-server`)

### client mode

#### flags
- `-gen-key` generate server/client key files

- `-address` udp address with port of remote upd server
- `-key-id` key file name of client

### server mode

#### flags
- `-port` port on which to run udp server
- `-timeframe` timestamp sent by client must be between <now-timeframe> and <now>
- `-command` the command to execute if udp packet sent by client is valid 
- `-key-id` key file name of server
