package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var br string

func main() {

	if len(os.Args) != 2 {
		os.Stderr.WriteString(usage())
		os.Exit(1)
	}

	depth, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		os.Stderr.WriteString("error: <depth> must be an integer")
		os.Exit(1)
	}

	if depth <= 0 {
		os.Exit(0)
	}

	pid := os.Getpid()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	go func() {
		for {
			fmt.Printf("[%d]: info: received signal: %v%s", pid, <-sigs, br)
		}
	}()

	fmt.Printf("[%d]: info: planted%s", pid, br)

	if depth > 1 {
		proc := exec.Command(os.Args[0], fmt.Sprintf("%d", depth-1))
		proc.Stderr = os.Stderr
		proc.Stdout = os.Stdout
		if err := proc.Start(); err != nil {
			fmt.Printf("[%d]: error: failed to start child process: %s%s", pid, err, br)
		}
	}

	for {
		<-time.After(500 * time.Millisecond)
	}
}

func init() {
	br = map[bool]string{
		true:  "\r\n",
		false: "\n",
	}[runtime.GOOS == "windows"]
}

func usage() string {
	msg := []string{
		"pstalk - Create an arbitrarily deep process tree",
		"",
		"usage: pstalk <depth>",
		"",
		"args",
		"  depth:    the depth of the process tree to create",
		"",
	}
	return strings.Join(msg, br)
}
