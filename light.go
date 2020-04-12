package launchpad

// LightEffects are effects applied to pad lights and are one of:
// EffectOff, EffectStatic, EffectFlash, EffectPulse.
type LightEffect int64

const (
	EffectOff    LightEffect = 0x80
	EffectStatic LightEffect = 0x90
	EffectFlash  LightEffect = 0x91
	EffectPulse  LightEffect = 0x92
)

// LightColor are palette-based colors used by the device for
// the pulse and flash effects, and for writing colors over the
// MIDI interface.
// LightColors are set by the devices, as colors may differ between
// devices.
type LightColor int64

// Light represents the state of a pad light
type Light struct {
	Effect LightEffect
	Color  LightColor
	Coord  Coordinate
	R      int8
	G      int8
	B      int8

	// DisplayLocked prevents Light redraws while true
	DisplayLocked bool
}

// ToggleDisplayLock is used by HitFuncs to mark a light
// to not be updated by the Grid during redraw cycles.
// This is useful if you want a layer of lights to persist over
// another layer.
func (l *Light) ToggleDisplayLock() {
	l.DisplayLocked = !l.DisplayLocked
}

// RGB sets RGB values on a Light
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
