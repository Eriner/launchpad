package launchpad

import (
	"log"
	"time"
)

// UseGrid launches the a grid's state machine on a given Launchpad
func UseGrid(lp Launchpad, g *Grid) {
	// start a listener for taps, recording the tap time.
	go func(p Launchpad, g *Grid) {
		c := p.Listen()
		for {
			tap := <-c
			tap.Time = time.Now()
			g.lastTap[tap.Coordinate] = tap.Time
			g.tapCount[tap.Coordinate]++
			g.taps <- tap
		}
	}(lp, g)
	// build and apply desired grid state
	go func(p Launchpad, g *Grid) {
		for {
			var lights []Light
			for _, pad := range g.Pads {
				if pad.Light.DisplayLocked {
					continue
				}
				lights = append(lights, pad.Light)
			}
			if err := p.LightSysEx(lights); err != nil {
				//TODO: gather MIDI error count and expose in
				// grid so that we can increase the renderDelay
				log.Println(err)
			}
			time.Sleep(g.renderDelay)
		}
	}(lp, g)
	go func(g *Grid) {
		tapsCh := g.Taps()
		for {
			tap := <-tapsCh
			pad := g.Pads[tap.Coordinate]
			go func(p *Pad, t Tap) {
				//TODO: handle errors below
				switch t.Type {
				case SingleTap:
					p.SingleTapHandler.Apply(p)
				case DoubleTap:
					p.DoubleTapHandler.Apply(p)
				}
			}(pad, tap)
		}
	}(g)
	return
}
