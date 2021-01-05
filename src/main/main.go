package main

import (
	"flag"
	"go-remote/src/client"
	"go-remote/src/server"
	"log"
)

func main() {
	keyUsage := "path to key file or base64 encoded key"
	keyArgument := flag.String("key", "", keyUsage)

	// client flags
	doGenKey := flag.Bool("gen-key", false, "generate key pair")
	addressArgument := "udp address"
	address := flag.String("address", "", addressArgument)

	// isServer flags
	isServer := flag.Bool("server", false, "run in isServer mode")
	port := flag.String("port", "8080", "udp port")
	timeFrame := flag.Int64("timeframe", int64(5), "timestamp in request must not be older than this timeframe (in seconds)")

	tmpfsDirUsage := "path to tmpfs directory, where 'start' file is expected"
	tmpfsDir := flag.String("tmpfs", "", tmpfsDirUsage)

	flag.Parse()

	pleaseProvideA := "Please provide a "
	if *keyArgument == "" {
		log.Fatal(pleaseProvideA + keyUsage)
	}

	if *isServer {
		if *tmpfsDir == "" {
			log.Fatal(pleaseProvideA + tmpfsDirUsage)
		}

		server.Run(*port, *keyArgument, *timeFrame, *tmpfsDir, nil)
	} else {
		if *address == "" {
			log.Fatal(pleaseProvideA + addressArgument)
		}

		client.Run(*doGenKey, *keyArgument, *address)
	}
}
