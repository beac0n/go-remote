package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	userPtr := flag.String("userName", "", "the name of the userName who is allowed to write to the tmpfs")
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

	tmpfsDir := "/tmp/go-remote"
	startFilePath := filepath.Join(tmpfsDir, "start")

	group, err := user.Lookup(userName)
	panicOnErr(err)
	uid, err := strconv.Atoi(group.Uid)
	panicOnErr(err)
	gid, err := strconv.Atoi(group.Gid)
	panicOnErr(err)

	panicOnErr(os.MkdirAll(tmpfsDir, os.FileMode(0200)))
	panicOnErr(syscall.Chown(tmpfsDir, uid, gid))

	_, err = os.Create(startFilePath)
	panicOnErr(err)

	panicOnErr(syscall.Chown(startFilePath, uid, gid))
	panicOnErr(syscall.Chmod(startFilePath, 0200))

	currentTime := time.Now().Local()
	panicOnErr(os.Chtimes(startFilePath, currentTime, currentTime))

	lastModTs, _ := getModTs(startFilePath)

	go func() {
		for {
			time.Sleep(10 * time.Millisecond)

			modTs, err := getModTs(startFilePath)
			if err != nil || modTs <= lastModTs {
				continue
			}

			executeCommand(startCommand, false)
			time.Sleep(time.Duration(commandTimeout) * time.Second)
			executeCommand(stopCommand, false)

			lastModTs = modTs
		}
	}()

	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, os.Interrupt)
	<-signalChannel

	_ = os.RemoveAll(tmpfsDir)
}

func panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getModTs(startFile string) (int64, error) {
	fileInfo, err := os.Stat(startFile)
	if err != nil {
		log.Println("ERROR:", err)
		return -1, err
	}

	return fileInfo.ModTime().Unix(), nil
}

func executeCommand(command string, terminateOnError bool) {
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

	if err != nil && !terminateOnError {
		log.Println("ERROR:", err)
	} else if err != nil && terminateOnError {
		log.Fatal("ERROR:", err)
	}
}
