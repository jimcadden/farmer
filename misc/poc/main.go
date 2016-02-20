// Copyright 2016
//
// Authors:
//   2016 Jim Cadden <jmcadden@bu.edu>
/*

farmer: Proof of Concept

Implements basic network isolation and VM allocation on a remote machine.

Steps:
- check remote prerequisites
- set up networks
- boot VM read stdout/stderr

*/
package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	//	"strings"

	"github.com/jmcadden/circuit/client"
)

func abort(format string, arg ...interface{}) {
	println("err:", fmt.Sprintf(format, arg...))
	os.Exit(1)
}

func printf(format string, arg ...interface{}) {
	fmt.Printf(format, arg...)
}

func debug(format string, arg ...interface{}) {
	println("dbg:", fmt.Sprintf(format, arg...))
}

//	usage: poc DIALIN_CIRCUIT_URL
//
func main() {
	switch len(os.Args) {
	case 2:
	//	arg := os.Args[2]
	default:
		println("usage: poc circuit://...")
		os.Exit(1)
	}
	println("circuit: dialing into", os.Args[1])
	c := client.Dial(os.Args[1], nil)

	prereqs := []string{"cat", "qemu-system-x86_64", "echo", "brctl"}

	servers := []string{}
	for _, r := range c.View() {
		err := check_prereqs(r, prereqs)
		if err {
			servers = append(servers, r.ServerID())
		}
	}
	if len(servers) > 0 {
		abort("following servers missing prereqs: %v", servers)
	}
}

// runShell executes the shell command on the given host,
// waits until the command completes and returns its output
// as a string. The error value is non-nil if the process exited in error.
func runShell(host client.Anchor, cmd string) (string, error) {
	return runShellStdin(host, cmd, "")
}

func runShellStdin(host client.Anchor, cmd, stdin string) (string, error) {
	defer func() {
		if recover() != nil {
			abort("connection to host lost")
		}
	}()
	job := host.Walk([]string{"shelljob", strconv.Itoa(rand.Int())})
	proc, _ := job.MakeProc(client.Cmd{
		Path:  "/bin/sh",
		Dir:   "/tmp",
		Args:  []string{"-c", cmd},
		Scrub: true,
	})
	go func() {
		io.Copy(proc.Stdin(), bytes.NewBufferString(stdin))
		proc.Stdin().Close() // Must close the standard input of the shell process.
	}()
	proc.Stderr().Close() // Close to indicate discarding standard error
	var buf bytes.Buffer
	io.Copy(&buf, proc.Stdout())
	stat, _ := proc.Wait()
	return buf.String(), stat.Exit
}

func check_prereqs(host client.Anchor, ps []string) (ret bool) {

	ret = false
	for _, v := range ps {
		out, err := runShell(host, "which "+v)
		if err != nil || len(out) == 0 {
			debug("[%s] missing prereq %s ", host.ServerID(), v)
			ret = true
		}
		//else {
		//	debug("[%s] %s found at %s", host.ServerID(), v, strings.TrimSpace(out))
		//}
	}
	return
}
