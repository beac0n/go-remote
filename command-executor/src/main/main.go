package main

import (
	"bytes"
	"flag"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	userPtr := flag.String("user-name", "", "the name of the user who is allowed to write to the socket")
	startCommandPtr := flag.String("command-start", "echo \"start!\"", "the command to execute when start is triggered")
	stopCommandPtr := flag.String("command-stop", "echo \"end!\"", "the command to execute after command-timeout is over")
	commandTimeoutPtr := flag.Int64("command-timeout", int64(60), "how long to wait before executing the end command")

	flag.Parse()

	if os.Geteuid() != 0 {
		log.Fatal("ERROR: Please run this program as root")
	}

	userName := *userPtr
	startCommand := *startCommandPtr
	stopCommand := *stopCommandPtr
	commandTimeout := *commandTimeoutPtr

	if userName == "" {
		flag.Usage()
		os.Exit(1)
	}

	userCredentials, err := user.Lookup(userName)
	panicOnErr(err)
	uid, err := strconv.Atoi(userCredentials.Uid)
	panicOnErr(err)
	gid, err := strconv.Atoi(userCredentials.Gid)
	panicOnErr(err)

	socketPath := "/tmp/go-remote.sock"
	listener, err := net.Listen("unix", socketPath)
	panicOnErr(syscall.Chown(socketPath, uid, gid))
	panicOnErr(syscall.Chmod(socketPath, 0200))
	panicOnErr(err)

	go func() {
		for {
			connection, err := listener.Accept()
			if err != nil {
				continue
			}

			executeCommand(startCommand)
			time.Sleep(time.Duration(commandTimeout) * time.Second)
			executeCommand(stopCommand)

			_ = connection.Close()
		}
	}()

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel

	panicOnErr(listener.Close())
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func executeCommand(command string) {
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

	log.Println("run '" + command + "'")
	err := cmd.Run()

	if stdOutBytes.Len() > 0 {
		log.Println("Stdout:", stdOutBytes.String())
	}

	if stdErrBytes.Len() > 0 {
		log.Println("Stderr:", stdErrBytes.String())
	}

	if err != nil {
		log.Println("ERROR:", err)
	}
}
