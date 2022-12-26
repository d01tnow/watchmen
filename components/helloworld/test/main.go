package main

import (
	"context"

	"fmswift.com.cn/watchmen/core/component"
	"github.com/smallnest/rpcx/client"
	"golang.org/x/exp/slog"
)

func main() {
	d, err := client.NewPeer2PeerDiscovery("tcp@localhost:4096", "")
	if err != nil {
		panic(err)
	}
	c := client.NewXClient("Component", client.Failtry, client.RandomSelect, d, client.DefaultOption)
	defer c.Close()

	desc := component.Description{}

	err = c.Call(context.Background(), "Description", 0, &desc)
	if err != nil {
		slog.Error("failed to call Description", err)
		return
	}
	slog.Info("component description", "name", desc.Name, "version", desc.Version)
	var stopReply component.StopReply
	err = c.Call(context.Background(), "Shutdown", component.ShutdownReasonRemaining, &stopReply)
	if err != nil {
		slog.Error("failed to call Shutdown", err)
		return
	}
	slog.Info("component shutdown", "code", stopReply.Code, "remark", stopReply.Remark)
}
