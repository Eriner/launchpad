package lpx

import (
	"time"

	"github.com/eriner/launchpad"
	"github.com/pkg/errors"
	"github.com/rakyll/portmidi"
)

var (
	// msgDeviceInquiry doesn't follow the normal message pattern
	msgDeviceInquiry = []byte{0xf0, 0x7e, 0x7f, 0x06, 0x01, 0xf7}

	sysExPrefix = []byte{0xf0, 0x00, 0x20, 0x29, 0x02, 0x0c}
)

var (
	ErrWrongMode = errors.New("launchpad: Launchpad X is not in the correct mode")
)

type Function byte
type Layout byte
type ProgramMode byte
type DeviceMode byte
type Aftertouch byte
type AftertouchType byte
type AftertouchThreshold byte

const (
	// FunctionMode changes the launchpad device mode to Standalone or DAW
	//
	// Mode is the Launchpad X's current mode.
	//
	// The Launchpad X offers two modes, Standalone and DAW.
	// These modes are mutually exclusive. Only one mode can
	// be active at a time.
	FunctionMode   Function   = 0x10
	ModeStandalone DeviceMode = 0x00
	ModeDAW        DeviceMode = 0x01

	// FunctionLayout changes the launchpad layout based on the mode
	FunctionLayout   Function = 0x00
	LayoutSession    Layout   = 0x00 // only selectable in DAW mode
	LayoutNote       Layout   = 0x01
	LayoutCustom1    Layout   = 0x04 // drum rack by factory default
	LayoutCustom2    Layout   = 0x05 // keys by factory default
	LayoutCustom3    Layout   = 0x06 // lighting mode in drum rack layout by factory default
	LayoutCutsom4    Layout   = 0x07 // lighting mode in session layout by factory default
	LayoutDAWFaders  Layout   = 0x0d // only selectable in DAW mode
	LayoutProgrammer Layout   = 0x7f

	// FunctionProgrammer changes the launchpad device
	//
	// Ref: When selecting Live mode with this message, Launchpad X switches
	// to Session layout, or Note mode when not in DAW mode.
	// When selecting Programmer mode using this SysEx message, the Setup entry
	// (holding down Session for half a second) is disabled. To return the Launchpad
	// X to normal operation, use this SysEx message to switch back to Live mode.
	FunctionProgramMode   Function    = 0x0e
	ProgramModeLive       ProgramMode = 0x00
	ProgramModeProgrammer ProgramMode = 0x01

	Black launchpad.LightColor = 0x00
	White launchpad.LightColor = 0x03
	Red   launchpad.LightColor = 0x05

	FunctionRGB         Function = 0x03
	FunctionLEDFeedback Function = 0x0a

	FunctionAftertouch        Function            = 0x0b
	AftertouchTypePolymorphic AftertouchType      = 0x00 // - 0: Polyphonic Aftertouch (Key Pressure events, A0h – AFh).
	AftertouchTypeChannel     AftertouchType      = 0x01 // - 1: Channel Aftertouch (Channel Pressure events, D0h – DFh)
	AftertouchTypeOff         AftertouchType      = 0x02
	AftertouchThresholdLow    AftertouchThreshold = 0x00
	AftertouchThresholdMed    AftertouchThreshold = 0x01
	AftertouchThresholdHigh   AftertouchThreshold = 0x02

	sysExSuffix byte = 0xf7
)

// Launchpad represents a device with input and output MIDI and DAW streams.
type Launchpad struct {
	// MIDI contains the device MIDI controller IO streams
	MIDI
	// DAW contains the device DAW controller IO streams
	DAW
	// mode is the device's current mode, either Standalone or DAW.
	// only one mode can be active at a time..
	mode DeviceMode

	// the input and output streams will be changed from MIDI to DAW
	// when ModeDAW() is called.
	inputStream  *portmidi.Stream
	outputStream *portmidi.Stream

	AppVersion  []byte
	BootVersion []byte
}

// Hit represents physical touches to Launchpad buttons.
type Hit struct {
	X int
	Y int
}

// Open opens a connection Launchpad and initializes an input and output
// stream to the currently connected device. If there are no
// devices are connected, it returns an error.
func Open() (*Launchpad, error) {
	midi := &MIDI{}
	if err := midi.Open(); err != nil {
		return nil, err
	}
	daw := &DAW{}
	if err := daw.Open(); err != nil {
		return nil, err
	}
	lp := &Launchpad{MIDI: *midi,
		DAW: *daw,
	}
	// by default, we use Standalone mode and provide
	// MIDI input and outputstreams
	if err := lp.Mode(ModeStandalone); err != nil {
		return nil, err
	}
	/*
		// we also get the app and boot versions
		if err := lp.DAW.outputStream.WriteSysExBytes(portmidi.Time(), msgDeviceInquiry); err != nil {
			lp.Close()
			return nil, err
		}
		time.Sleep(10 * time.Millisecond)
		resp, err := lp.inputStream.ReadSysExBytes(5)
		if err != nil {
			panic(err)
		}
		//TODO: response does not correspond with "device inquiry message" in manual
		// needs to be investigated and fixed
		//spew.Dump(resp)
		_ = resp
	*/
	return lp, nil
}

func (l *Launchpad) Close() error {
	// reset everything
	l.Mode(ModeStandalone)
	l.ProgramMode(ProgramModeLive)
	l.inputStream = nil
	l.outputStream = nil
	var retErr error
	if err := l.MIDI.Close(); err != nil {
		retErr = err
	}
	if err := l.DAW.Close(); err != nil {
		retErr = err
	}
	return retErr
}

// Mode switches the Launchpad X into between Standalone mode and DAW mode
func (l *Launchpad) Mode(m DeviceMode) error {
	l.inputStream = l.MIDI.inputStream
	l.outputStream = l.MIDI.outputStream
	if err := l.msg(FunctionMode, []byte{byte(m)}); err != nil {
		// if we're ever unable to switch modes, the device is broken. and we need to abort
		if cErr := l.Close(); err != nil {
			return errors.Wrap(err, cErr.Error())
		}
		return err
	}
	return nil
}

func (l *Launchpad) Layout(lay Layout) error {
	return l.msg(FunctionLayout, []byte{byte(lay)})
}

func (l *Launchpad) ProgramMode(pm ProgramMode) error {
	return l.msg(FunctionProgramMode, []byte{byte(pm)})
}

func (l *Launchpad) Test() {
	if err := l.MIDI.outputStream.WriteShort(0x80, 0x51, 0x00); err != nil {
		panic(err)
	}
}

func (l *Launchpad) Light(light launchpad.Light) error {
	err := l.MIDI.outputStream.WriteShort(int64(light.Effect), int64(light.Coord), int64(light.Color))
	time.Sleep(5 * time.Millisecond)
	return err
}

func (l *Launchpad) LightSysEx(lights []launchpad.Light) error {
	var colorspec []byte
	for _, light := range lights {
		colorspec = append(colorspec, LightRGBSysEx(&light)...)
	}
	err := l.msg(FunctionRGB, colorspec)
	return err
}

func LightRGBSysEx(light *launchpad.Light) []byte {
	/*
		Host => Launchpad X:
		Hex: F0h 00h 20h 29h 02h 0Ch 03h <colourspec> [<colourspec> [...]] F7h Dec: 240 0 32 41 2 12 3 <colourspec> [<colourspec> [...]] 247
		The <colourspec> is structured as follows:
			- Lighting type (1 byte)
			- LED index (1 byte)
			- Lighting data (1 – 3 bytes)
		Lighting types:
			- 0: Static colour from palette, Lighting data is 1 byte specifying palette entry.
			- 1: Flashing colour, Lighting data is 2 bytes specifying Colour B and Colour A.
			- 2: Pulsing colour, Lighting data is 1 byte specifying palette entry.
			- 3: RGB colour, Lighting data is 3 bytes for Red, Green and Blue (127: Max, 0: Min).
	*/
	var out []byte
	switch light.Effect {
	case launchpad.EffectStatic:
		out = append(out, 0x03) //
		out = append(out, Colorspec(light.Coord, light.R, light.G, light.B)...)
	case launchpad.EffectPulse:
		out = append(out, 0x02)
		out = append(out, byte(light.Coord))
		out = append(out, approximatePalatte(light.R, light.G, light.B))

	}
	return out

}

//TODO: implement this function. This is just a placeholder
func approximatePalatte(r, g, b int8) byte {
	return byte(White)
}

// Clear sends the clear command to the device. Note that calling Clear on the device
// will be overwritten by the state of any launchpad.Grid elements
func (l *Launchpad) Clear() error {
	//TODO
	return nil
}

func (l *Launchpad) Aftertouch(attype AftertouchType, atthresh AftertouchThreshold) error {
	var args []byte
	args = append(args, byte(attype), byte(atthresh))
	return l.msg(FunctionAftertouch, args)
}

// LEDFeedback configures the device's LEDFeedback setting.
// BUG: I haven't been able to get this to work for some reason.
func (l *Launchpad) LEDFeedback(internal, external bool) error {
	var i, e byte
	if internal {
		i = 0x01
	}
	if external {
		e = 0x01
	}
	var args []byte
	args = append(args, i, e)
	return l.msg(FunctionLEDFeedback, args)
}

// Listen returns launchpad button presses
func (l *Launchpad) Listen() <-chan launchpad.Tap {
	ch := make(chan launchpad.Tap)
	go func(pad *Launchpad, ch chan launchpad.Tap) {
		for {
			time.Sleep(5 * time.Millisecond)
			hits, err := pad.Read()
			if err != nil {
				continue
			}
			for i := range hits {
				ch <- hits[i]
			}
		}
	}(l, ch)
	return ch
}

// Read returns events from the MIDI stream. This includes button presses
func (l *Launchpad) Read() (taps []launchpad.Tap, err error) {
	var evts []portmidi.Event
	if evts, err = l.MIDI.inputStream.Read(64); err != nil {
		return
	}
	for _, evt := range evts {
		i := int(evt.Data1)
		x := i % 10
		y := i / 10
		tap := launchpad.Tap{
			Time:       time.Now(),
			X:          x,
			Y:          y,
			Coordinate: launchpad.Coord(x, y),
		}
		taps = append(taps, tap)
	}
	return
}

// msg sends messages to the launchpad over the DAW interface, leaving MIDI open for use
func (l *Launchpad) msg(function Function, args []byte) error {
	err := l.DAW.outputStream.WriteSysExBytes(portmidi.Time(), msg(function, args))
	time.Sleep(5 * time.Millisecond)
	return err
}

// Colorspec creates a single light command for execution with LightSysEx
func Colorspec(c launchpad.Coordinate, r, g, b int8) []byte {
	var colorspec []byte
	colorspec = append(colorspec, byte(c))
	if r < 0 {
		r = -r
	}
	if g < 0 {
		g = -g
	}
	if b < 0 {
		b = -b
	}
	rgb := []byte{byte(r), byte(g), byte(b)}
	colorspec = append(colorspec, rgb...)
	return colorspec
}

// msg builds SysEx messages into the appropriate format
func msg(function Function, args []byte) []byte {
	msg := sysExPrefix
	msg = append(msg, byte(function))
	msg = append(msg, args...)
	msg = append(msg, sysExSuffix)
	return msg
}
