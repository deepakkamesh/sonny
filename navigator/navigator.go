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

	"github.com/deepakkamesh/sonny/devices"
	"github.com/golang/glog"
)

const cellSz = 35 // size of the cell in centimeters.
const maxX = 100  // number of grid X coordinates
const maxY = 100  // number of grid Y.

// struct cell represents a single cell in the occupancy grid.
type cell struct {
	occupied bool      // Is the cell occupied?
	lastUpd  time.Time // Last updated timestamp for cell.
	obs      uint      // number of observations.
	posObs   float64   // Number of observations the cell is occupied.
}

// ogrid represents the occupancy grid.
type Ogrid struct {
	sonny  *devices.Sonny
	cells  [maxX][maxY]cell
	curr_x int
	curr_y int
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
		return int(x/cellSz) + 1, int(y/cellSz) + 1
	}
	return
}

//ResetMap resets the grid
func (s *Ogrid) ResetMap() {

	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			s.cells[x][y].occupied = false
			s.cells[x][y].obs = 0
			s.cells[x][y].posObs = 0
		}
	}
}

// UpdateMap updates the occupany grid map based on lidar readings
func (s *Ogrid) UpdateMap() error {

	minAngle := 20
	shiftAngle := 5

	s.ResetMap() //TODO: Remove after testing.

	// Get forward sweep radar readings with retry in case of I2C failures..
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
	}
	if failure {
		return fmt.Errorf("Failed to update map: %v", err)
	}

	// From the returned beam update the cell location.
	// TODO: Update all the cells till the occupied cells as non occupied.
	for i := 0; i < len(rangeReading); i++ {
		x, y := calcCell(uint(rangeReading[i]), uint(minAngle+i*shiftAngle))
		xAbs := s.curr_x + x
		yAbs := s.curr_y + y
		s.cells[xAbs][yAbs].occupied = true
		s.cells[xAbs][yAbs].obs += 1
		s.cells[xAbs][yAbs].posObs += 1
	}

	return nil
}

// GenerateMap returns a png buffer with the current map.
func (s *Ogrid) GenerateMap() (*bytes.Buffer, error) {
	img := image.NewRGBA(image.Rect(0, 0, maxX, maxY))

	// TODO: color the cell based on the probability.
	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			if s.cells[x][y].occupied == true {
				img.Set(x, y, color.RGBA{255, 0, 0, 255})
				continue
			}
			img.Set(x, y, color.RGBA{255, 0, 0, 10})
		}
	}

	// Set rover location on map.
	img.Set(s.curr_x, s.curr_y, color.RGBA{100, 50, 0, 255})

	buff := new(bytes.Buffer)
	if err := png.Encode(buff, img); err != nil {
		return nil, err
	}

	return buff, nil
}
