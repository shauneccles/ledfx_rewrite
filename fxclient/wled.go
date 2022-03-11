package fxclient

import (
	"net"
	"sync"
)

// WLED TODO populate with proper fields
type WLED struct {
	mu sync.RWMutex

	Virtuals map[string]*Virtual

	// conn TODO populate with a udpconn
	conn *net.UDPConn
}

func NewWLED(ip string) (w *WLED, err error) {
	// TODO Initialize UDP session with WLED
	return &WLED{
		mu:       sync.RWMutex{},
		Virtuals: make(map[string]*Virtual),
	}, nil
}
