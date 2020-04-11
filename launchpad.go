// launchpad provides interfaces for Launchpad devices, which currently include:
//   - Launchpad (original)
//   - Launchpad mk2
//   - Launchpad X
//
package launchpad

import (
	"sync"
	"time"
)

type Launchpad interface {
	Close() error
	Clear() error
	Listen() <-chan Coordinate
	Light(Light) error
	LightSysEx([]Light) error
}

var (
	defaultRenderDelay = 50 * time.Millisecond
)

// This is designed similarly to http.HandlerFunc
type HitHandler interface {
	Apply(*Pad) error
}

// HitFunc is an adapter to use arbitrary Go functions
// to interact with Pads when they are hit.
type HitFunc func(*Pad) error

// Apply returns f(p)
func (f HitFunc) Apply(p *Pad) error {
	return f(p)
}

func NewPad() *Pad {
	return &Pad{
		SingleTapHandler: HitFunc(func(p *Pad) error {
			return nil
		}),
		DoubleTapHandler: HitFunc(func(p *Pad) error {
			return nil
		}),
		hitFuncMu: &sync.Mutex{},
	}
}

// Pads are the buttons on the Launchpad device
type Pad struct {
	Light
	// HitFuncs are triggered when the pad has been presesed
	SingleTapHandler HitHandler
	DoubleTapHandler HitHandler
	// Only one HitFunc should ever be launched at a time.
	hitFuncMu *sync.Mutex
	// rollingHitsRecord contains a rolling record of hit events for a given
	// time window
	rollingHitsRecord []time.Time
}
