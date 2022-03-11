package fxclient

import (
	"fmt"
	"sync"
)

type Handler struct {
	mu    sync.RWMutex
	WLEDs map[string]*WLED
}

func NewHandler() (h *Handler) {
	return &Handler{
		mu:    sync.RWMutex{},
		WLEDs: make(map[string]*WLED),
	}
}

func (h *Handler) AddWLED(ip string) error {
	wled, err := NewWLED(ip)
	if err != nil {
		return fmt.Errorf("error creating new WLED: %w", err)
	}
	h.WLEDs[ip] = wled
	return nil
}
