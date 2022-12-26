package main

import "flag"

type args struct {
	daemonPort int
	listenPort int
}

func parseFlag() *args {
	var a args
	flag.IntVar(&a.daemonPort, "daemon-port", 0, "daemon listening port")
	flag.IntVar(&a.listenPort, "p", 0, "listen port")
	flag.Parse()

	return &a
}
