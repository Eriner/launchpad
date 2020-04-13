package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/eriner/launchpad"
	"github.com/eriner/launchpad/pkg/lpx"
	"github.com/eriner/launchpad/pkg/middleware"
)

func main() {
	// open the launchpad device, in this case a Launchpad X
	lp, err := lpx.Open()
	if err != nil {
		die(err)
	}
	// catch interrupts to exit programmer mode when we ctrl+C
	ic := make(chan os.Signal, 1)
	signal.Notify(ic, os.Interrupt, syscall.SIGTERM)
	go func(pad *lpx.Launchpad) {
		<-ic
		pad.Close()
		os.Exit(1)
	}(lp)
	// switch to programmer mode, which gives us control over the lights
	if err := lp.ProgramMode(lpx.ProgramModeProgrammer); err != nil {
		log.Fatalf("error setting launchpad program mode: %v", err)
	}
	//
	// Grids are state machines that hold and continually apply the
	// desired state of the button grid to the launchpad.
	//
	// Grids are composed of Pads. Pads have HitFuncs which are called
	// when buttons are pressed.
	//
	// create a new grid, testGrid, which maintains a desired grid state
	testGrid, err := launchpad.NewGrid(lp)
	if err != nil {
		die(err)
	}
	// In theory, we could have mulitple grids or devices. UseGrid activates a grid on
	// a launchpad.
	launchpad.UseGrid(lp, testGrid)
	// loop over all the devices to set a HitFunc, which is the function that activates on
	// a button press. Note that button press events are limited to one every 200 milliseconds.
	// If it has been less than 200ms since the last time HitFunc was called, it will be not execute.
	for x := 1; x < 9; x++ {
		for y := 1; y < 9; y++ {
			pad := testGrid.Pad(x, y)
			// Set all of the lights to Red.
			pad.Light.RGB(127, 0, 0)

			// Demonstration of using middleware to wrap a handler for single tap events
			pad.SingleTapHandler = logTap(
				middleware.SimulatedFeedbackInverted(
					pad.SingleTapHandler, time.Second*3,
				))
			// And another for double-tap events, but without the logTap middleware func
			pad.DoubleTapHandler = middleware.SimulatedFeedbackPulseToggle(pad.DoubleTapHandler)
		}
	}
	// here we override the double-tap handler for the bottom left pad.
	pad := testGrid.Pad(1, 1)
	pad.DoubleTapHandler = launchpad.HitFunc(func(p *launchpad.Pad) error {
		log.Println("overridden double-tap: no pulsing for this corner!")
		return nil
	})

	// Now that we have assigned handlers for our button presses and
	// the state machine is running, we can just sleep forever
	select {}
}

// logTap is an example of how to create middleware for pad hit event handlers
//
// logTap will print the X and Y positions of a pad when pressed (and wrapped
// around a handler).
func logTap(next launchpad.HitHandler) launchpad.HitHandler {
	return launchpad.HitFunc(func(p *launchpad.Pad) error {
		x, y := p.Light.Coord.XY()
		log.Printf("tap! Pos: X=%d, Y=%d", x, y)
		next.Apply(p)
		return nil
	})
}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
