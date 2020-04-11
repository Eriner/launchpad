package launchpad

type LightEffect int64
type LightColor int64

const (
	// Effects are lighting effects that can be applied to the pads
	// using SysEx
	EffectOff    LightEffect = 0x80
	EffectStatic LightEffect = 0x90
	EffectFlash  LightEffect = 0x91
	EffectPulse  LightEffect = 0x92
)

// Light represents the state of a pad light
type Light struct {
	Effect LightEffect
	Color  LightColor
	Coord  Coordinate
	R      int8
	G      int8
	B      int8

	DisplayLocked bool
}

// ToggleDisplayLock is used by HitFuncs to mark a light
// to not be updated by the Grid during redraw cycles.
// This is useful if you want a layer of lights to persist over
// another layer.
func (l *Light) ToggleDisplayLock() {
	l.DisplayLocked = !l.DisplayLocked
}

func (l *Light) RGB(r, g, b int8) {
	if r < 0 {
		r = -r
	}
	if g < 0 {
		g = -g
	}
	if b < 0 {
		b = -b
	}
	l.R = r
	l.G = g
	l.B = b
}
