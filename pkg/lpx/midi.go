package lpx

import (
	"errors"
	"strings"

	"github.com/rakyll/portmidi"
)

// MIDI is the MIDI interface for the Launchpad X
type MIDI struct {
	inputStream  *portmidi.Stream
	outputStream *portmidi.Stream
}

// Open the Launchpad X MIDI interface
func (m *MIDI) Open() error {
	input, output, err := discoverMIDI()
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
	m.inputStream = inStream
	m.outputStream = outStream
	return nil
}

// Close the Launchpad X MIDI interface
func (m *MIDI) Close() error {
	if err := m.inputStream.Close(); err != nil {
		return err
	}
	if err := m.outputStream.Close(); err != nil {
		return err
	}
	return nil
}

// discoverMIDI provides the Launchpad X MIDI device
func discoverMIDI() (input portmidi.DeviceID, output portmidi.DeviceID, err error) {
	in := -1
	out := -1
	for i := 0; i < portmidi.CountDevices(); i++ {
		info := portmidi.Info(portmidi.DeviceID(i))
		if strings.Contains(info.Name, "Launchpad X LPX MIDI") {
			if info.IsInputAvailable {
				in = i
			}
			if info.IsOutputAvailable {
				out = i
			}
		}
	}
	if in == -1 || out == -1 {
		err = errors.New("launchpad: no Launchpad X MIDI device is connected")
	} else {
		input = portmidi.DeviceID(in)
		output = portmidi.DeviceID(out)
	}
	return
}
