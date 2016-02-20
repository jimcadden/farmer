// Copyright 2016
//
// Authors:
//   2016 Jim Cadden <jmcadden@bu.edu>

/*

farmer: Proof of Concept

Implements the basic network isolation and VM allocation mechanism on a remote
machine.

Stages:
- check remote prerequisites
- set up networks
- boot VM read stdout/stderr

*/
package main

import (
	"os"

	"github.com/jmcadden/circuit/client"
)

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

	for _, r := range c.View() {
		check_prereqs(c, r)
	}
}

func check_prereqs(c *client.Client, a client.Anchor) (ret bool) {

	ret = true

	println("checking %s:", a.Addr())

	// QEMU
	qemu_check := client.Cmd{
		Path:  "/usr/bin/which",
		Args:  []string{"qemu-system-x86_64"},
		Scrub: true,
	}
	qemu_loc, err := a.Walk([]string{"loc", "qemu"}).MakeProc(qemu_check)
	if err != nil {
		println("qmeu_check start error", err.Error())
		os.Exit(1)
	}
	//stdout := qemu_loc.Stdout().Read()
	//	println("location of qemu:%s", string(stdout))

	if err := qemu_loc.Peek().Exit; err != nil {
		println("qmeu_check runtime error", err.Error())
		os.Exit(1)
	}
	qemu_loc.Stdin().Close()
	return
}
