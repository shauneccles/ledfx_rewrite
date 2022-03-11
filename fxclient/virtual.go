package fxclient

import "sync"

// Virtual TODO populate with proper fields
type Virtual struct {
	mu sync.RWMutex

	Pixels []Pixel
}

func NewVirtual() (v *Virtual) {
	return &Virtual{
		mu:     sync.RWMutex{},
		Pixels: make([]Pixel, 0),
	}
}
