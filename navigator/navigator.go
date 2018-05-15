/* Package navigator provides self driving capability for Sonny based on a simple occupancy
grid algorith */
package navigator

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"time"

	"github.com/deepakkamesh/go-roomba/constants"
	"github.com/deepakkamesh/sonny/devices"
	"github.com/golang/glog"
)

const cellSz = 35 // size of the cell in centimeters. Cell is cellSz x CellSz square.
const maxX = 100  // number of grid X coordinates
const maxY = 100  // number of grid Y.

// struct cell represents a single cell in the occupancy grid.
type cell struct {
	occupied int       // Is the cell occupied 1 = occupied 0 = unoccupied, -1 = unknown
	lastUpd  time.Time // Last updated timestamp for cell.
	obs      uint      // number of observations.
	posObs   float64   // Number of (positive) observations the cell is occupied.
}

// ogrid represents the occupancy grid.
type Ogrid struct {
	sonny  *devices.Sonny
	cells  [maxX][maxY]cell
	curr_x int     // current x,y location of bot.
	curr_y int     // current x,y location of bot.
	orient float64 // Current orientation of bot.
}

// NewOGrid returns a initialized Ogrid structure.
func NewOgrid(s *devices.Sonny) *Ogrid {
	return &Ogrid{
		sonny:  s,
		cells:  [maxX][maxY]cell{},
		curr_x: 50,
		curr_y: 50,
	}
}

// calcCell finds takes in the beam and angle and calculates the cell
// which contains the obstacle relative to the current cell.
// left bottom is considered 0,0.
func calcCell(line, angle uint) (x, y int) {

	switch {
	case angle == 90:
		return 0, int(line/cellSz) + 1

	// Relative cell coordinate; anything to the left of the rover is negative realtively.
	case angle < 90:
		x := math.Cos(float64(angle)*math.Pi/180) * float64(line)
		y := math.Sin(float64(angle)*math.Pi/180) * float64(line)
		return -1 * (int(x/cellSz) + 1), int(y/cellSz) + 1

	// Relative cell coordinate; anything to the right of the rover is  positive.
	case angle > 90:
		angle := 180 - angle
		x := math.Cos(float64(angle)*math.Pi/180) * float64(line)
		y := math.Sin(float64(angle)*math.Pi/180) * float64(line)
		return 1 * (int(x/cellSz) + 1), int(y/cellSz) + 1
	}
	return
}

//ResetMap resets the grid
func (s *Ogrid) ResetMap() {

	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			s.cells[x][y].occupied = -1
			s.cells[x][y].obs = 0
			s.cells[x][y].posObs = 0
		}
	}
}

// MoveForward moves the bot forward by number of cells. It returns the delta
// of movement in mm. If positive it overshot the location and if negative it undershot.
func (s *Ogrid) MoveForward(cells int) (delta int, err error) {
	vel := 300 // Speed in mm/s.

	// Get starting readings.
	// TODO: Ignoring RIGHT ENcoder. Need to refactor to include both.
	p, err := s.sonny.Sensors(constants.SENSOR_LEFT_ENCODER)
	if err != nil {
		return 0, err
	}
	encStart := int16(p[0])<<8 | int16(p[1])

	desiredDist := cells * cellSz * 10               // Desired distance to move in mm.
	driveTime := float32(desiredDist) / float32(vel) // Assume fixed speed of 500 mm/s

	glog.Infof("%v %v", desiredDist, driveTime)
	// Move bot (time*speed).
	if err := s.sonny.DirectDrive(int16(vel), int16(vel)); err != nil {
		return 0, err
	}
	time.Sleep(time.Duration(driveTime*1000) * time.Millisecond)
	if err := s.sonny.DirectDrive(0, 0); err != nil {
		return 0, err
	}

	// Calculate if we overshot or undershot landing and return delta.
	p, err = s.sonny.Sensors(constants.SENSOR_LEFT_ENCODER)
	if err != nil {
		return 0, err
	}
	encEnd := int16(p[0])<<8 | int16(p[1])
	switch {
	case encEnd > encStart:
		distTravelled := float32(encEnd-encStart) * math.Pi * 72.0 / 508.8
		return (desiredDist - int(distTravelled)), nil

	case encEnd < encStart:
		distTravelled := float32(32767-encStart+encEnd) * math.Pi * 72.0 / 508.8
		return (desiredDist - int(distTravelled)), nil
	}

	return
}

// UpdateMap updates the occupany grid map based on lidar readings
func (s *Ogrid) UpdateMap() error {

	minAngle := 20
	shiftAngle := 10 // reduce.TODO

	s.ResetMap() //TODO: Remove after testing.

	// Get forward sweep radar readings with retry in case of I2C failures.
	var (
		err          error
		rangeReading []int32
		failure      bool = true
	)
	for i := 0; i < 3; i++ {
		rangeReading, err = s.sonny.ForwardSweep(shiftAngle)
		if err == nil {
			failure = false
			break
		}
		glog.Warningf("Forward sweep failed retry#%v: %v", i, err)
		time.Sleep(1 * time.Second)
	}
	if failure {
		return fmt.Errorf("Failed to update map: %v", err)
	}

	// From the returned beam update the cell location.
	// TODO: Update all the cells till the occupied cells as non occupied.
	for i := 0; i < len(rangeReading); i++ {
		x, y := calcCell(uint(rangeReading[i]), uint(minAngle+i*shiftAngle))
		// TODO: This should take into account orientation of bot to the grid.
		xAbs := s.curr_x + x
		yAbs := s.curr_y + y

		s.cells[xAbs][yAbs].occupied = 1
		s.cells[xAbs][yAbs].obs += 1
		s.cells[xAbs][yAbs].posObs += 1
		if x == 0 {
			continue
		}

		// Calculate free cells till occupied tell.
		m := float32(y) / float32(x)
		b := float32(yAbs) - m*float32(xAbs)
		//	glog.Infof("XX %v= %v %v %v %v\n", rangeReading[i], xAbs, yAbs, m, b)
		for j := s.curr_y; j < yAbs; j++ {
			xF := (float32(j) - float32(b)) / m
			//glog.Infof("xy %v %v %v", j, yF, m)
			s.cells[int(xF)][j].occupied = 0
		}

	}

	return nil
}

// GenerateMap returns a png buffer with the current map of the grid.
// TODO: This is a CPU intensive operation. Need an optimization.
func (s *Ogrid) GenerateMap() (*bytes.Buffer, error) {
	m := 2 // Number of pixels per cell. m x m.
	img := image.NewRGBA(image.Rect(0, 0, maxX*m, maxY*m))

	// TODO: color the cell based on the probability.
	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			switch s.cells[x][y].occupied {
			case 1:
				fillCell(img, x, y, m, color.RGBA{255, 0, 0, 255})
			case 0:
				fillCell(img, x, y, m, color.RGBA{0, 255, 0, 255})
			case -1:
				fillCell(img, x, y, m, color.RGBA{20, 20, 20, 255})
			}
		}
	}

	// Set rover location on map.
	fillCell(img, s.curr_x, s.curr_y, m, color.RGBA{100, 50, 0, 255})
	buff := new(bytes.Buffer)
	if err := png.Encode(buff, img); err != nil {
		return nil, err
	}

	return buff, nil
}

// fillCell renders the cell at x, y with a size of scale x scale.
func fillCell(img *image.RGBA, x int, y int, scale int, c color.RGBA) {
	for i := 0; i < scale; i++ {
		for j := 0; j < scale; j++ {
			img.Set(x*scale+i, y*scale+j, c)
		}
	}

}
