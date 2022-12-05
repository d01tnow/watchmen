/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"
)

func main() {
	fl, err := NewFlock("/tmp/watchmen-daemon.pid")
	if err != nil {
		panic(err)
	}
	if err := fl.TryLock(); err != nil {
		fmt.Println(err)
		os.Exit(kExitCodeAlreadyRunning)
	}

	defer fl.Close()

	c := parseFlag()

	var d Daemon
	d.Init(WithRendezvous(c.Rendezvous),
		WithPort(c.MdnsListenPort),
		WithWebPort(c.WebListenPort),
	)

	d.Run()
}
