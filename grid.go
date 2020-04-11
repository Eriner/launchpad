package launchpad

import "time"

// NewGrid provides a grid state machine for an opened Launchpad device
func NewGrid(lp Launchpad) (*Grid, error) {
	g := &Grid{
		Pads:        make(map[Coordinate]*Pad),
		renderDelay: defaultRenderDelay,
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

func (c *Coordinate) XY() (x, y int) {
	x = int(*c) % 10
	y = int(*c) / 10
	return
}

func Coord(x, y int) Coordinate {
	return Coordinate((y * 10) + x)
}

// Grid is a state-machine representation of the desired Pad grid state
type Grid struct {
	Pads map[Coordinate]*Pad
	// renderDelay is how long the state machine pauses after
	// each full-grid redraw
	renderDelay time.Duration
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
