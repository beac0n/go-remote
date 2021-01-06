package main

import (
	"flag"
	"go-remote/src/client"
	"go-remote/src/server"
	"os"
)

func main() {
	keyBase64 := flag.String("key", "", "base64 encoded key")

	// client flags
	doGenKey := flag.Bool("gen-key", false, "generate base64 encoded aes key")
	address := flag.String("address", "", "udp address")

	// isServer flags
	isServer := flag.Bool("server", false, "run in isServer mode")
	port := flag.String("port", "8080", "udp port")
	timeFrame := flag.Int64("timeframe", int64(5), "timestamp in request must not be older than this timeframe (in seconds)")

	tmpfsDir := flag.String("tmpfs", "", "path to tmpfs directory, containing 'start' file")

	flag.Parse()

	if *isServer {
		if *tmpfsDir == "" || *keyBase64 == "" {
			flag.Usage()
			os.Exit(1)
		}

		server.Run(*port, *keyBase64, *timeFrame, *tmpfsDir, nil)
	} else {
		client.Run(*doGenKey, *keyBase64, *address)
	}
}
