// launchpad provides an interface for the Launchpad X device.
//
//
package launchpad

import (
	"sync"
	"time"
)

type Launchpad interface {
	Close() error
	// Clear wipes all pads to default states and issues
	// a device clear command.
	Clear() error
	// Listen collects coordinates of pad presses
	Listen() <-chan Coordinate

	// Light applies palatte-based lights over the MIDI channel
	//NOTE: LightSysEx should be preferred over Light, as LightSysEx
	// uses the DAW interface and does not crowd the MIDI I/O.
	Light(Light) error
	// LightSysEx uses the DAW I/O interface to apply Light configurations.
	LightSysEx([]Light) error
}

var (
	// defaultRenderDelay is the amount of time the pad light
	// rendering loop waits after each cycle.
	// Reducing this value too much will cause distortion.
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

// NewPad returns an empty, default Pad
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
