package launchpad

import (
	"log"
	"time"
)

// UseGrid launches the a grid's state machine on a given Launchpad
func UseGrid(lp Launchpad, g *Grid) {
	// start a listener to trigger HitFunc when pad is pressed
	go func(p Launchpad, g *Grid) {
		c := p.Listen()
		for {
			hit := <-c
			hitTime := time.Now()
			pad := g.Pads[hit]
			pad.rollingHitsRecord = append(pad.rollingHitsRecord, hitTime)
		}
	}(lp, g)
	// build and apply desired grid state
	go func(p Launchpad, g *Grid) {
		for {
			var lights []Light
			for _, pad := range g.Pads {
				/*
					if pad.X == 9 || pad.Y == 9 {
						continue
					}
				*/
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
	// cleanup rollingHitsRecord for pads
	go func(g *Grid) {
		for {
			for _, pad := range g.Pads {
				pad.hitFuncMu.Lock()
				for _, hit := range pad.rollingHitsRecord {
					// decision time! single, tripple, or double tap?
					if time.Since(hit) > 200*time.Millisecond {
						hitEventCount := len(pad.rollingHitsRecord) / 2
						switch hitEventCount {
						case 0:
							continue
						case 1:
							go pad.SingleTapHandler.Apply(pad)
						case 2:
							go pad.DoubleTapHandler.Apply(pad)
						}
						// clear the hits
						pad.rollingHitsRecord = make([]time.Time, 0)
					}
				}
				// we wait an extra while to prevent spillover
				// events from a hold from progressing
				pad.hitFuncMu.Unlock()
			}
			time.Sleep(50 * time.Millisecond)
		}
	}(g)
	return
}
