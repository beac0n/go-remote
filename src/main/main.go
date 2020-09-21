package main

import (
	"flag"
	"go-remote/src/client"
	"go-remote/src/server"
)

func main() {
	keyFilePath := flag.String("key", "", "path to key file")

	// client flags
	doGenKey := flag.Bool("gen-key", false, "generate key pair")
	address := flag.String("address", "", "udp address")

	// isServer flags
	isServer := flag.Bool("server", false, "run in isServer mode")
	port := flag.String("port", "8080", "udp port")
	timeFrame := flag.Int64("timeframe", int64(5), "timestamp in request must not be older than this timeframe (in seconds)")
	commandStart := flag.String("command-start", "echo 'start!'", "the command to execute before the command timeout")
	commandTimeout := flag.Int64("command-timeout", int64(60), "how long to wait before executing the end command")
	commandEnd := flag.String("command-end", "echo 'end!'", "the command to execute after the command timeout")

	flag.Parse()

	if *isServer {
		server.Run(*port, *keyFilePath, *timeFrame, *commandStart, *commandTimeout, *commandEnd, make(chan bool))
	} else {
		client.Run(*doGenKey, *keyFilePath, *address)
	}
}
