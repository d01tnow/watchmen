package main

import (
	"time"

	"github.com/libp2p/go-libp2p/core/peer"
)

type Peer struct {
	p         peer.AddrInfo
	FoundAt   time.Time // 最早发现时间
	UpdatedAt time.Time // 最后更新时间
}
