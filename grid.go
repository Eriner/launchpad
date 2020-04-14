package launchpad

import (
	"time"
)

// NewGrid provides a grid state machine for an opened Launchpad device.
// This grid synchronizes its state to the launchpad after UseGrid has been
// called.
func NewGrid(lp Launchpad) (*Grid, error) {
	g := &Grid{
		Pads:        make(map[Coordinate]*Pad),
		renderDelay: defaultRenderDelay,
		taps:        make(chan Tap, 1024),
		tapChs:      make([]chan Tap, 0),
		tapCount:    make(map[Coordinate]int),
		lastTap:     make(map[Coordinate]time.Time),
		isDepressed: make(map[Coordinate]bool),
	}
	for x := 1; x < 10; x++ {
		for y := 1; y < 10; y++ {
			// This is our default Pad initializer
			coord := Coord(x, y)
			light := &Light{Coord: coord,
				Effect: EffectStatic,
			}
			pad := NewPad()
			pad.Light = *light
			g.Pads[coord] = pad
		}
	}
	return g, nil
}

// Coordinates are a position on the pad
type Coordinate int64

// XY returns the X and Y
func (c *Coordinate) XY() (x, y int) {
	x = int(*c) % 10
	y = int(*c) / 10
	return
}

// Coord converts X and Y coordinates into type Coordinate
func Coord(x, y int) Coordinate {
	return Coordinate((y * 10) + x)
}

// Grid is a state-machine made of Pads that represents  of the desired
// Pad grid state.
type Grid struct {
	Pads map[Coordinate]*Pad
	// renderDelay is how long the state machine pauses after
	// each full-grid redraw.
	renderDelay time.Duration
	// taps is fanned out by Taps() and records grid tap events.
	taps chan Tap
	// tapChs stores channels that we fan out in Taps()
	tapChs []chan Tap
	// tapCount maintains a record of tap times for coordinates in the
	// last 200ms
	tapCount map[Coordinate]int
	// lastTap records the last time a coordinate was tapped so we only
	// ever process the latest tap event.
	lastTap map[Coordinate]time.Time
	// isDepressed is true if a button is pressed down, false when button
	// is lifted.
	isDepressed map[Coordinate]bool
}

// Pad returns a pad for a given set of X and Y coordinates
func (g *Grid) Pad(x, y int) *Pad {
	return g.Pads[Coord(x, y)]
}

// Clear resets all elements in the state machine to their default state
func (g *Grid) Clear() {
	//TODO

}

func (g *Grid) Close() error {
	//TODO: when a grid is closed, we should Clear()
	return nil
}

// Taps returns a channel of tap events associated with a grid.
func (g *Grid) Taps() chan Tap {
	tapCh := make(chan Tap)
	if len(g.tapChs) != 0 {
		g.tapChs = append(g.tapChs, tapCh)
		return tapCh
	}
	// The code below here is the setup for the Taps() function
	// and will be executed on first invocation of Taps()
	g.tapChs = append(g.tapChs, tapCh)
	go func(gg *Grid) {
		for {
			tap := <-gg.taps
			go func(grid *Grid, t Tap) {
				// 200ms is the time we wait to determine if this was a
				// single or double tap
				time.Sleep(200 * time.Millisecond)
				//log.Println(t.HoldDuration)
				t.DecisionTime = time.Now()
				switch tc := g.tapCount[t.Coordinate]; tc {
				case 0: //NOTE: this often occurs after double taps
					return
				case 1:
					t.Type = SingleTap
					g.tapCount[t.Coordinate] = 0
				case 2:
					t.Type = DoubleTap
					g.tapCount[t.Coordinate] = 0

				}
				for _, tapCh := range grid.tapChs {
					tapCh <- t
				}

			}(gg, tap)
		}

	}(g)
	return tapCh
}
