package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"net/http"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/multiformats/go-multiaddr"
	"golang.org/x/exp/slog"

	"github.com/gin-gonic/gin"
)

type Daemon struct {
	host       host.Host
	rendezvous string
	listenHost string // mdns listening address
	listenPort int    // mdns listening port

	webListenPort int // web listening port

	mu    sync.Mutex
	peers map[peer.ID]*Peer
}

type option func(*Daemon)

const (
	kDefaultRendezvous    = "watchmen.daemon"
	kDefaultListenHost    = "0.0.0.0"
	kDefaultListenPort    = 13578
	kDefaultWebListenPort = 13579

	kVersion = "daemon-0.0.1"
)

func (d *Daemon) Run() {
	d.initOptionWithDefaultWhenNeeded()
	d.initHost()
	d.initMdns()
	d.serve()
}

// HandlePeerFound - 实现 mdns.Notifee 接口, 处理' 发现mdns 服务'
func (d *Daemon) HandlePeerFound(pi peer.AddrInfo) {
	slog.Info("found peer", "peer", pi.ID, "addr", pi.Addrs)
	_, ok := d.peers[pi.ID]
	if !ok {
		// not found, create
		d.peers[pi.ID] = &Peer{
			p:         pi,
			FoundAt:   time.Now(),
			UpdatedAt: time.Now(),
		}
		return
	}
	// found, update
	d.peers[pi.ID].UpdatedAt = time.Now()
}

func (d *Daemon) Init(opts ...option) {
	for _, opt := range opts {
		opt(d)
	}
}

func WithHost(host string) option {
	return func(d *Daemon) {
		d.listenHost = host
	}
}
func WithPort(port int) option {
	return func(d *Daemon) {
		d.listenPort = port
	}
}
func WithWebPort(port int) option {
	return func(d *Daemon) {
		d.webListenPort = port
	}
}

func WithRendezvous(rendezvous string) option {
	return func(d *Daemon) {
		d.rendezvous = rendezvous
	}
}

func (d *Daemon) initOptionWithDefaultWhenNeeded() {
	if d.rendezvous == "" {
		d.rendezvous = kDefaultRendezvous
	}
	if d.listenHost == "" {
		d.listenHost = kDefaultListenHost
	}
	if d.listenPort < 1 {
		d.listenPort = kDefaultListenPort
	}
	if d.webListenPort < 1 {
		d.webListenPort = kDefaultWebListenPort
	}
}

func (d *Daemon) initHost() {
	r := rand.Reader
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		panic(err)
	}
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(
		fmt.Sprintf("/ip4/%s/tcp/%d", d.listenHost, d.listenPort),
	)

	d.host, err = libp2p.New(
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		panic(err)
	}
	slog.Debug("mdns is listening", "addr", sourceMultiAddr)

}

// initMdns - 初始化 mdns 发现
func (d *Daemon) initMdns() {
	ser := mdns.NewMdnsService(d.host, d.rendezvous, d)
	if err := ser.Start(); err != nil {
		panic(ser)
	}

}

func (d *Daemon) serve() {
	router := gin.New()
	router.GET("/version", d.handleVersion)
	router.GET("/component/all", d.componentAll)
	addr := d.webListenAddr()

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT,
	)
	defer stop()
	go func() {
		slog.Debug("web server is listening", "addr", addr)
		err := srv.ListenAndServe()
		if err != nil {
			slog.Error("failed to serve.", err, "addr", addr)
			stop()
		}
	}()
	<-ctx.Done()

}

func (d *Daemon) webListenAddr() string {
	return fmt.Sprintf("127.0.0.1:%d", d.webListenPort)
}

func (d *Daemon) handleVersion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": kVersion,
	})
}
func (d *Daemon) componentAll(c *gin.Context) {

}
