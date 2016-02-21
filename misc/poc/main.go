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

	// connect to circuit network
	c := client.Dial(os.Args[1], nil)

	// check for prereqs
	prereqs := []string{"route", "brctl", "qemu-system-x86_64", "dnsmasq", "tunctl", "ip"}
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
	debug("prereq check complete")
	// set up networks
	count := 0
	for _, r := range c.View() {
		count += 1
		err := setup_network(r, count)
		if err {
			abort("[%s] network setup error: %v", r.ServerID(), err)
		}
	}
	debug("network setup complete")
	//
	//
}

func setup_network(host client.Anchor, id int) bool {

	//sid := host.ServerID()
	// ids
	ip_base := "10.0.0." + strconv.Itoa(id)
	//ip_new := "10.0.1." + strconv.Itoa(id)
	eth0 := "eth0"
	eth00 := "eth0.1"
	br0 := "br_poc"
	tap0 := "tap_poc"
	orquit := " || exit 1"
	// commands
	commands := []string{}
	commands = append(commands, "ip link add link "+eth0+" "+eth00+" type vlan id 1")
	commands = append(commands, "ip link set "+eth00+" up")
	commands = append(commands, "brctl addbr "+br0)
	commands = append(commands, "ip link set "+br0+" up")
	commands = append(commands, "ip addr add "+ip_base+"/24 dev "+br0)
	commands = append(commands, "ip tuntap add dev "+tap0+" mode tap multi_queue")
	//commands = append(commands, "ip addr add "+ip_new+"/24 dev "+tap0)
	commands = append(commands, "ip link set "+tap0+" up")
	commands = append(commands, "brctl addif "+br0+" "+tap0+" "+eth00)

	for _, v := range commands {
		cmd := v + orquit
		out, err := runShell(host, cmd)
		debug("[%s] "+cmd, host.ServerID())
		if err != nil {
			debug("[%s] %s ", host.ServerID(), err)
		} else {
			debug("[%s] "+out, host.ServerID())
		}
	}
	return false
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
