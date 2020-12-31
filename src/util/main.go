package util

import (
	"bytes"
	"encoding/binary"
	"log"
	"os"
	"os/exec"
	"strings"
)

func Check(err error, reason string) {
	if err == nil {
		return
	} else if reason == "" {
		os.Exit(1)
	} else {
		LogError(err, reason)
		os.Exit(1)
	}
}

func LogError(err error, reason string) {
	if err == nil {
		log.Println("ERROR: ", reason)
	} else {
		log.Println("ERROR: ", reason, ": ", err)
	}
}

func GetSourcePort(keyBytes []byte) int {
	hashedPublicKeyBytes := GetHashFromBytes(keyBytes)

	sourcePort := int(binary.LittleEndian.Uint16(hashedPublicKeyBytes))
	if sourcePort < 1024 {
		return 1024
	}

	return sourcePort
}

func ExecuteCommand(command string) {
	commandSplit := strings.Split(command, " ")
	commandSplitLen := len(commandSplit)

	var cmd *exec.Cmd
	if commandSplitLen == 0 {
		return
	} else if commandSplitLen == 1 {
		cmd = exec.Command(commandSplit[0])
	} else {
		cmd = exec.Command(commandSplit[0], commandSplit[1:]...)
	}

	var stdOutBytes bytes.Buffer
	var stdErrBytes bytes.Buffer
	cmd.Stdout = &stdOutBytes
	cmd.Stderr = &stdErrBytes

	log.Println("running command", command)
	if err := cmd.Run(); err != nil {
		log.Println("ERROR when running command:", err)
	}
	log.Println("Stdout:", stdOutBytes.String())
	log.Println("Stderr:", stdErrBytes.String())
}
