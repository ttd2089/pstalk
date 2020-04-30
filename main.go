package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var newline string

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

	if depth == 0 {
		os.Exit(0)
	}

	fmt.Printf("planted %d%s", os.Getpid(), newline)

	if depth > 1 {
		proc := exec.Command(os.Args[0], fmt.Sprintf("%d", depth-1))
		proc.Stderr = os.Stderr
		proc.Stdout = os.Stdout
		if err := proc.Start(); err != nil {
			os.Stderr.WriteString("error: failed to start child process")
			os.Exit(1)
		}
	}

	for {
		<-time.After(500 * time.Millisecond)
	}
}

func init() {
	newline = map[bool]string{
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
	return strings.Join(msg, newline)
}
