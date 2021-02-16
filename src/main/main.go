package main

import (
	"flag"
	"go-remote/src/client"
	"go-remote/src/server"
	"os"
)

func main() {
	keyBase64 := flag.String("key", "", "base64 encoded aes key")

	// client flags
	doGenKey := flag.Bool("gen-key", false, "generate base64 encoded aes key")
	address := flag.String("address", "", "udp address")

	// isServer flags
	isServer := flag.Bool("server", false, "run in server mode")
	port := flag.String("port", "8080", "udp port")
	timeFrame := flag.Int64("timeframe", int64(5), "timestamp in request must not be older than this timeframe (in seconds)")

	flag.Parse()

	if *isServer {
		if *keyBase64 == "" {
			flag.Usage()
			os.Exit(1)
		}

		server.Run(*port, *keyBase64, *timeFrame, nil)
	} else {
		client.Run(*doGenKey, *keyBase64, *address)
	}
}
