package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
)

const usage = `
pstalk - Create an arbitrarily deep process tree

usage: pstalk <depth>

args
	depth:    the depth of the process tree to create
`

func main() {

	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "%s", usage)
		os.Exit(1)
	}

	depth, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: <depth> must be an integer")
		os.Exit(1)
	}

	if depth <= 0 {
		os.Exit(1)
	}

	pid := os.Getpid()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	exit := make(chan struct{}, 1)
	go func() {
		for {
			s := <-sigs
			fmt.Printf("[%d]: info: received signal: %v\n", pid, s)
			if s == syscall.SIGINT || s == syscall.SIGTERM || s == syscall.SIGKILL {
				exit <- struct{}{}
				return
			}
		}
	}()

	fmt.Printf("[%d]: info: planted\n", pid)

	if depth > 1 {
		proc := exec.Command(os.Args[0], fmt.Sprintf("%d", depth-1))
		proc.Stderr = os.Stderr
		proc.Stdout = os.Stdout
		if err := proc.Start(); err != nil {
			fmt.Printf("[%d]: error: failed to start child process: %s\n", pid, err)
		}
	}

	<-exit
}
