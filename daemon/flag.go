package main

import "flag"

type config struct {
	Rendezvous     string
	MdnsListenHost string
	MdnsListenPort int
	WebListenPort  int
}

func parseFlag() config {
	var c config
	flag.StringVar(&c.Rendezvous, "rendezvous", kDefaultRendezvous, "rendezvous string")
	flag.StringVar(&c.MdnsListenHost, "h", kDefaultListenHost, "mdns listening address")
	flag.IntVar(&c.MdnsListenPort, "p", kDefaultListenPort, "mdns listening port")
	flag.IntVar(&c.WebListenPort, "P", kDefaultWebListenPort, "web listening port")
	flag.Parse()
	return c
}
