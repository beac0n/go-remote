package main

import (
	"flag"
	"go-remote/src/client"
	"go-remote/src/server"
	"log"
)

func main() {

	clientMode := flag.Bool("client", false, "run in client mode")
	serverMode := flag.Bool("server", false, "run in server mode")
	keyId := flag.String("key-id", "", "key id")

	// client flags
	doGenKey := flag.Bool("gen-key", false, "generate key pair")
	address := flag.String("address", "", "udp address")

	// server flags
	port := flag.String("port", "8080", "udp port")
	timeFrame := flag.Int64("timeframe", int64(5), "timestamp in request must not be older than this timeframe (in seconds)")
	commandStart := flag.String("command-start", "echo 'start!'", "the command to execute before the command timeout")
	commandTimeout := flag.Int64("command-timeout", int64(60), "how long to wait before executing the end command")
	commandEnd := flag.String("command-end", "echo 'end!'", "the command to execute after the command timeout")

	flag.Parse()

	if *clientMode == *serverMode {
		log.Fatal("either run in client mode (-client) or server mode (-server)")
	}

	if *clientMode {
		client.Run(doGenKey, keyId, address)
	} else if *serverMode {
		server.Run(port, keyId, timeFrame, commandStart, commandTimeout, commandEnd)
	}

}
