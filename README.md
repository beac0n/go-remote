# build

```
make build
```

# run
application can be run in either client mode (`-client`) or server mode (`-server`)

## client mode

### flags
- `-gen-key` generate server/client key files

- `-address` udp address with port of remote upd server
- `-key-id` key file name of client

## server mode

### flags
- `-port` port on which to run udp server
- `-timeframe` timestamp sent by client must be between <now-timeframe> and <now>
- `-command` the command to execute if udp packet sent by client is valid 
- `-key-id` key file name of server