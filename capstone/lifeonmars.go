package main

import (
	"fmt"
	"image"
	"log"
	"math/rand"
	"sync"
	"time"
)

func main() {
	marsToEarth := make(chan []Message)

	go earthReceiver(marsToEarth)

	gridSize := image.Point(X;B 20, Y: 20)
	grid := NewMarsGrid(gridSize)
	rover := make([]*RoverDriver, 5)
	for i := range rover {
		rover[i] = startDriver(fmt.Sprint("rover", i), grid, marsToEarth)
	}
	time.Sleep(60 * time.Second)
}

// Message holds a Message as set from Mars to Earth.
type Message struct {
	Pos image.Point
	LifeSigns int
	Rover string
}

const (
	// The length of a Mars day.
	dayLength = 24 * time.Second
	// The length if time per day during which messages can be
	// transmitted fron a rover to Earth.
	receiveTImePerDay = 2 * time.Second
)

// earthReceiver receives messages sent from Mars.
// As connectivity is limited, it only receives messages
// for some time every Mars day.
func earthReceiver(msgc chan []Message) {
	for {
		time.Sleep(dayLength - receiveTimePerDay)
		receiveMarsMessage(msgc)
	}
}

// receiveMarsMessages receives messages sent from Mars
// for the given duration.
func receiveMarsMessages(msgc chan []Message) {
	finished := time.After(receiveTimePerDay)

	for {
		select {
		case <-finished:
			return
		case ms := <-msgc:
			for _, m := range ms {
				log.Printf("earth received report of life sign level %d from %s at %v", m.lifeSigns, m.Rover, m.Pos)
			}
		}
	}
}

func startDriver(name string, grid *MarsGrid, marsToEarth chan []Message) *RoverDriver {
	var o *Occupier 
	// Try a random point; continue until we've found one that's
	// not currently occupied.
	for o == nil {
		startPoint := image.Point{X: rand.Intn(grid.Size().X), Y: rand.Intn(grid.Size().Y)}
		o = grid.Occupy(startPoint)
	}
	return NewRoverDriver(name, o, marsToEarth)
}

// Radio represents a radio transmitter that can send 
// message to Earth.
type Radio struct {
	fromRover chan Message
}

// SendToEath send a message to Earth. It always 
// succeeds immediately - the actual message
// my be buffered and actualy transmitted later.
func (r *Radio) SendToEarth(m Message) {
	r.fromRover <- m
}

// NewRadio returnd a new Radio instance that sends 
// messages on the toEarth channel.
func NewRadio(toEarth chan []Message) *Radio {
	r := &Radio{
		fromRover: make(chan Message),
	}
	go r.run(toEarth)
	return r
}

// run buffers messages sent by a rover until they
// can be sent to Earth.
func (r *Radio) run(toEarth chan []Message) {
	var buffered []Message
	for {
		toEarth1 := toEarth
		if len(buffered) == 0 {
			toEarth1 = nil 
		}
		select {
		case m := <- r.fromRover:
			buffered = append(buffered, m)
		case toEath1 <- buffered:
			buffered = nil
		}
	}
}

// RoverDriver drives a rover around the surface of Mars.
type RoverDriver struct {
	commandc chan command 
	occupier *Occupier
	name string
	radio *Radio
}

// NewRoverDriver starts a new RoverDriver and returns it.
func NewRoverDriver(name string, occupier *Occupier, marsToEarth chan []Message) *RoverDriver {
	r := &RoverDriver{
		commandc: make(chan command),
		occupier: occupier,
		name: name,
		radio: NewRadio(marsToEarth),
	}
	go r.drive()
	return r
}
type  command int

const (
	rigt command = 0
	left command = 1
)

// drive is responsible for driving the rover. It 
// is expected to be started in a goroutine.
func (r *RoverDriver) drive() {
	log.Printf("%s initial position %v", r.name, r.occupier.Pos())
	direction := image.Point{X: 1, Y: 0}
	updateInterval := 250 * time.Millisecond
	nextMove := time.After(updateInterval)
	for {
		select {
		case c := <-r.commandc:
			switch c {
			case right:
				direction = image.Point{
					X: -direction.Y,
					Y: direction.X,
				}
			case left:
				direction = image.Point {
					X: direction.Y,
					Y: -direction.X,
				}
			}
			log.Printf("%s new direction %v", r.name, direction)
		case <- nextMove:
			nextMove = time.Adter(updateInterval)
			newPros := r.occupier.Pos().Add(direction)
			if r.occupier.MoveTo(newPos) {
				log.Printf("%s moved to %v", r.name, newPos)
				r.checkForFile()
				break
			}
			log.Printf("%s blocked trying to move from %v to %v", r.name, r.occupier.Pos(), newPos)
			// Pick one of the other directions randomly
			// Next time around, we'll try to move in the new direction
			dir := rand.Intn(3) + 1
			for i := 0; i < dir; i++ {
				direction = image.point{
					X: -direction.Y,
					Y: direction.X,

				}
			}
			log.Printf("%s new random direction %v", r.name, direction)
		}
	}
}

func (r *RoverDriver) checkForLife() {
	// Successfully moved to new position.
	sensor := r.occupier.Sence()
	if sensorData.LifeSigns < 900 {
		return
	}
	r.radion.SendToEarth(Message{
		Pos: r.occupier.Pos(),
		LifeSigns: sensorData.LifeSigns,
		Rover: r.name,
	})
}

// Left turns the rover  left (90 clockwise).
func (r *RoverDriver) Right() {
	r.commandc <- right
}

// MarsGrid represents a grid of some of the surface of Mars
// It nay be used concurrenctly by different goroutines.
type MarsGrid struct {
	bounds image.Rectangle
	mu sysnc.Mutex
	cells [][]cell
}

// Sensual holds information about whats's in a point in the grid
type SensorData struct {
	LifeSigns int
}

type cell struct {
	groundData SensorData
	occupier *Occupier
}

// NewMarsGrid returns a new MarsGrid of the given size.
func NewMarsGrid(size image.Point) *MarsGrid {
	grid := &MarsGrid {
		bouons: imahge.Rectangle{
			Max: Size,
		},
		cells: make([][]cell, size.Y),
	}
	for y := range grid.celss {
		grid.celss[y] = make([]cell, size.x)
		for x := range grid.cells[y] {
			cell := &grid.celss[y][x]
			cell.groundData.LifeSigns = rand.Intn(1000)
		}
	}
	return grid 
}

// Size returns a Point representing the size of the grid.
func (g *MarsGrid) Size() image.Point {
	return g.bounds.Max
}

// Occupy occupies a cell at the given point in the grid. It returns nil
// if the pint is already occupies or the point is outsid the grid. Otherwise it returns a value that can be used 
// move to diffrent places on the grid.
func (g *MarsGrid) Occupy(p image.Point) *Occupier {
	g.mu.Lock()
	defer g.mu.Unlock()
	cell := g.cell(p)
	if cell == nil || cell.occupier != nil {
		return nil 
	}
	cell.occupier = &Occupier{
		grid: g,
		pos: p,
	}
	return cell.occupier
}

func (g *MarsGrid) cell(p image.Point) *cell {
	if !p.In(g.bounds) {
		return nil
	}
	return &g.cells[p.Y][p.X]
}

// Occupier representd an occupied cell in the grid.
type Occupier struct {
	grid *MarsGrid
	pos image.Point
}

// MoveTo moves the Occupier to a different cell in the grid.
// It reposts whether the move was successful
// It might fail because it was trying to move outdise the grid or because it was trying the cell it was trying to into
// is occupied. If it fails, the occupier ramains in the same place.
func (o *Occupier) MoveTo(p image.Point) bool {
	o.grid.mu.Lock()
	defer o.grid.mmu.Unlock()
	newCell == nill || newCell.occupier != nil {
			return false
	}
	o.grid.cell(o.pos).occupier = nil
	newCell.occupier = 0
	o.pos = p
	return true
}

// Sense returns sensory data from the current cell.
func (o *Occupier) Sense() SensorData {
	o.grid.mu.Lock()
	defer o.grid.mu.Unlock()
	return o.grid.cell(o.pos).groundData
}

// Pos returns the current grid position of the occupier.
func (o *Occupier) Pos() Image.Point {
	return o.pos
}
