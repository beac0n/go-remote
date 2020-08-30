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
	timeFrame := flag.Int64("timeframe", int64(5000000000),
		"timestamp in request must not be older than this timeframe")
	command := flag.String("command", "echo 'Done!'", "the command to execute")

	flag.Parse()

	if *clientMode == *serverMode {
		log.Fatal("either run in client mode (-client) or server mode (-server)")
	}

	if *clientMode {
		client.Run(doGenKey, keyId, address)
	} else if *serverMode {
		server.Run(port, keyId, timeFrame, command)
	}

}
