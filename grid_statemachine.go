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
			g.isDepressed[tap.Coordinate] = !g.isDepressed[tap.Coordinate]
			// when button has lifted after a press
			if !g.isDepressed[tap.Coordinate] {
				//BUGFIX: sometimes isDepressed becomes inverted from the actual pad state,
				// meaning the HoldDuration becomes time.Now().Sub(g.lastTap[tap.Coordinate])
				// I'm not quite sure why this happens, or if this is just a normal desync
				// of the Launchpad.
				// To handle this, we invert it here if we detect this.
				// button presses run on a 200ms clock.
				holdDuration := tap.Time.Sub(g.lastTap[tap.Coordinate])
				if holdDuration > 200*time.Millisecond {
					//NOTE: "desyncs" here are at least sometimes just overruns of our detection
					// window, meaning this conditional may useful if we ever create a HoldTap
					// tap type.
					//log.Println("fixed desync")
					g.isDepressed[tap.Coordinate] = !g.isDepressed[tap.Coordinate]
				}
				tap.HoldDuration = holdDuration
				g.tapCount[tap.Coordinate]++
				g.taps <- tap
			}
			g.lastTap[tap.Coordinate] = tap.Time
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
