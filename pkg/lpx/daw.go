package lpx

import (
	"errors"
	"strings"

	"github.com/rakyll/portmidi"
)

// DAW is the DAW interface for the Launchpad X
type DAW struct {
	inputStream  *portmidi.Stream
	outputStream *portmidi.Stream
}

// Open the Launchpad X DAW interface
func (d *DAW) Open() error {
	input, output, err := discoverDAW()
	if err != nil {
		return err
	}
	var inStream, outStream *portmidi.Stream
	if inStream, err = portmidi.NewInputStream(input, 1024); err != nil {
		return err
	}
	if outStream, err = portmidi.NewOutputStream(output, 1024, 0); err != nil {
		return err
	}
	d.inputStream = inStream
	d.outputStream = outStream
	return nil
}

// Close the Launchpad X DAW interface
func (d *DAW) Close() error {
	// When the DAW interface is closed, we must switch to session mode
	//
	// Ref: When the DAW/software exits, it should send SysEx to revert
	// the device to Standalone mode. Doing this ensures that all the state
	// is cleared, and the device remains useful as a standalone device once
	// the DAW is done using it (without power cycling to restore it).
	if err := d.outputStream.WriteSysExBytes(portmidi.Time(), msg(FunctionMode, []byte{byte(ModeStandalone)})); err != nil {
		return err
	}
	if err := d.inputStream.Close(); err != nil {
		return err
	}
	if err := d.outputStream.Close(); err != nil {
		return err
	}
	return nil
}

// discoverDAW provides the Launchpad X DAW device
func discoverDAW() (input portmidi.DeviceID, output portmidi.DeviceID, err error) {
	in := -1
	out := -1
	for i := 0; i < portmidi.CountDevices(); i++ {
		info := portmidi.Info(portmidi.DeviceID(i))
		if strings.Contains(info.Name, "Launchpad X LPX DAW") {
			if info.IsInputAvailable {
				in = i
			}
			if info.IsOutputAvailable {
				out = i
			}
		}

	}
	if in == -1 || out == -1 {
		err = errors.New("launchpad: no Launchpad X DAW device is connected")
	} else {
		input = portmidi.DeviceID(in)
		output = portmidi.DeviceID(out)
	}
	return
}
