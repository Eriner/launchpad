// middleware provides some example HitHandler middlewares
package middleware

import (
	"time"

	"github.com/eriner/launchpad"
)

// SimulatedFeedback lights the buttons with RGB values for a duration
func SimulatedFeedback(next launchpad.HitHandler, r, g, b int8, t time.Duration) launchpad.HitHandler {
	return launchpad.HitFunc(func(p *launchpad.Pad) error {
		// save current state
		cr := p.Light.R
		cg := p.Light.G
		cb := p.Light.B
		p.Light.RGB(r, g, b)
		time.Sleep(t)
		p.Light.RGB(cr, cg, cb)
		next.Apply(p)
		return nil
	})
}

// SimulatedFeedbackInverted inverts the colors of the pressed button for a duration
func SimulatedFeedbackInverted(next launchpad.HitHandler, t time.Duration) launchpad.HitHandler {
	return launchpad.HitFunc(func(p *launchpad.Pad) error {
		// save current state
		cr := p.Light.R
		cg := p.Light.G
		cb := p.Light.B
		p.Light.RGB(127-cr, 127-cg, 127-cb)
		time.Sleep(t)
		p.Light.RGB(cr, cg, cb)
		next.Apply(p)
		return nil
	})
}

// SimulatedFeedbackPulseToggle causes the lights to pulse when pressed
func SimulatedFeedbackPulseToggle(next launchpad.HitHandler) launchpad.HitHandler {
	return launchpad.HitFunc(func(p *launchpad.Pad) error {
		switch p.Light.Effect {
		case launchpad.EffectStatic:
			p.Light.Effect = launchpad.EffectPulse
		case launchpad.EffectPulse:
			p.Light.Effect = launchpad.EffectStatic
		}
		next.Apply(p)
		return nil
	})
}
