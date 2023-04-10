package util

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func IsLinux() bool {
	return runtime.GOOS == "linux"
}

func IsPortInUdpSourcePorts(sourcePort int) bool {
	for _, currentPort := range getUdpSourcePorts() {
		if currentPort == int64(sourcePort) {
			return true
		}
	}
	return false
}

func getUdpSourcePorts() []int64 {
	var sourcePorts []int64

	path := "/proc/net/udp"
	procFile, err := os.Open(path)
	if err != nil {
		log.Fatal("error while reading proc file:", err)
	}
	defer procFile.Close()

	procFileReader := bufio.NewReader(procFile)
	for {
		procFileLine, err := procFileReader.ReadBytes('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal("error while reading proc file line:", err)
		}
		sourcePorts = append(sourcePorts, getSourcePortFromLine(string(bytes.Trim(procFileLine, "\t\n "))))
	}
	if len(sourcePorts) == 0 {
		log.Fatal("can't read proc file: /proc/net/udp has no content")
	}
	// Remove header line
	return sourcePorts[1:]
}

func getSourcePortFromLine(line string) int64 {
	parts := strings.Split(line, " ")
	filtered := parts[:0]
	for _, part := range parts {
		if part != "" {
			filtered = append(filtered, part)
		}
	}

	if len(filtered) <= 2 {
		log.Fatal("expected proc file line to have more entries", filtered)
	}

	localIpAndPort := filtered[1]

	if strings.Contains(localIpAndPort, ":") {
		hex := strings.Split(localIpAndPort, ":")[1]
		value, err := strconv.ParseInt(hex, 16, 64)
		if err != nil {
			log.Fatal("could not parse hex string to int64:", err)
		}

		return value
	} else {
		return 0
	}
}
