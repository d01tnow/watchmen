package main

import (
	"fmt"

	"fmswift.com.cn/watchmen/core/component"
)

type evh struct {
	d component.Discoverer
}

func (e *evh) SetDiscoverer(d component.Discoverer) {
	e.d = d
}

func (e *evh) OnStart() {
	fmt.Println("evn started")
}
func (e *evh) OnShutdown() {

	fmt.Println("evn shutdown")
}

func main() {
	arg := parseFlag()
	fmt.Println("Hello World! daemon port:", arg.daemonPort)
	h := evh{}
	component.Serve(&h,
		component.WithNameAndVersion("HelloWrold", 100),
	)
}
