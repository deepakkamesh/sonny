/* Package navigator provides self driving capability for Sonny based on a simple occupancy
grid algorithm */
package navigator

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"math"
	"time"
)

// size of the cell in centimeters. Cell is cellSz * CellSz square. The smaller the cell
// the more accurate the space calculation, but takes more memory.
const (
	cellSz     = 10
	maxX       = 1000 // number of grid X coordinates.
	maxY       = 1000 // number of grid Y coordinates.
	stOccupied = 1    // occupancy state of cell.
	stFree     = 0    // occupancy state of cell.
	stUnknown  = -1   // occupancy state of cell.

)

// struct cell represents a single cell in the occupancy grid.
type cell struct {
	occupied int        // Is the cell occupied 1 = occupied 0 = unoccupied, -1 = unknown
	lastUpd  time.Time  // Last updated timestamp for cell.
	obs      uint       // number of observations.
	posObs   float64    // Number of (positive) observations the cell is occupied.
	occColor color.RGBA // Color of the point if occupied..
}

// ogrid represents the occupancy grid.
type Ogrid struct {
	cells   [maxX][maxY]cell
	curr_x  int // current x location of bot.
	curr_y  int // current y location of bot.
	imgChan chan chan *bytes.Buffer
}

// NewOGrid returns a initialized Ogrid structure.
func NewOgrid() *Ogrid {
	return &Ogrid{
		cells:   [maxX][maxY]cell{},
		curr_x:  500,
		curr_y:  500,
		imgChan: make(chan chan *bytes.Buffer),
	}
}

func (s *Ogrid) SetPos(x, y int) {
	s.curr_x = x
	s.curr_y = y
}

//ResetMap resets the grid
func (s *Ogrid) ResetMap() {

	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			s.cells[x][y].occupied = stUnknown
			s.cells[x][y].obs = 0
			s.cells[x][y].posObs = 0
		}
	}
}

// Placeholder for any goroutines to be started.
func (s *Ogrid) StartGrid() {
}

// UpdateMap updates the occupany grid map based on lidar readings.
// minAngle is the starting angle (degrees) in reference to the bot.
// shiftAngle is the angle between readings in degrees.
func (s *Ogrid) UpdateMap(rangeReading []int32, minAngle int, shiftAngle int, posture float64, c color.RGBA) error {

	// From the returned beam update the cell location.
	for i := 0; i < len(rangeReading); i++ {

		servoAngle := minAngle + i*shiftAngle
		effAngle := posture - 90 + float64(servoAngle)
		if effAngle < 0 {
			effAngle = 360 + effAngle
		}
		if effAngle > 360 {
			effAngle = effAngle - 360
		}

		// Get relative X,Y points from line and angle.
		line := rangeReading[i]
		x := (math.Cos(effAngle*math.Pi/180) * float64(line))
		y := (math.Sin(effAngle*math.Pi/180) * float64(line))
		// Calculate absolute X,Y points on map.
		xAbs := s.curr_x + int(x)
		yAbs := s.curr_y + int(y)

		s.cells[xAbs][yAbs].occupied = stOccupied
		s.cells[xAbs][yAbs].obs += 1
		s.cells[xAbs][yAbs].posObs += 1
		s.cells[xAbs][yAbs].occColor = c

		// TODO: if x is zero then just update Y as free till blocked cell.
		// TODO: Review these calculations.
		if x == 0 {
			continue
		}

		// Calculate free cells till occupied cell using line equation.
		// Y = mX + b  where m is slope and b is intercept.
		if xAbs-s.curr_x == 0 {
			continue
		}
		m := float32(yAbs-s.curr_y) / float32(xAbs-s.curr_x)
		b := float32(yAbs) - m*float32(xAbs)
		for j := yAbs + 1; j < s.curr_y; j++ {
			xF := (float32(j) - float32(b)) / m
			s.cells[int(xF)][j].occupied = stFree
			s.cells[int(xF)][j].obs += 1
		}
		for j := s.curr_y; j < yAbs-1; j++ {
			xF := (float32(j) - float32(b)) / m
			s.cells[int(xF)][j].occupied = stFree
			s.cells[int(xF)][j].obs += 1
		}
	}

	return nil
}

// GenerateMap() returns a png map of the environment. It wraps
// scaledMap or normalMap.
func (s *Ogrid) GenerateMap() (*bytes.Buffer, error) {
	return s.scaledMap()
}

// GenerateMap returns a png buffer with the current map of the grid.
// TODO: This is a CPU intensive operation. Need an optimization.
func (s *Ogrid) normalMap() (*bytes.Buffer, error) {

	m := 1 // Number of pixels per cell. m x m.
	img := image.NewRGBA(image.Rect(0, 0, maxX*m, maxY*m))

	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {

			// Image has 0,0 in the top left. Adjusting 0,0 to
			// bottom left for simplifying view.
			yAdj := maxY - y

			// Mark cells.
			switch s.cells[x][y].occupied {
			case stOccupied:
				fillCell(img, x, yAdj, m, s.cells[x][y].occColor)
			case stUnknown:
				fillCell(img, x, yAdj, m, color.RGBA{20, 20, 20, 10})
			case stFree:
				fillCell(img, x, yAdj, m, color.RGBA{0, 255, 0, 255})
			}

			// Set Grid lines.
			if x%cellSz == 0 || y%cellSz == 0 {
				fillCell(img, x, yAdj, m, color.RGBA{194, 194, 214, 255})
			}

		}
	}

	// Set rover location on map.
	fillCell(img, s.curr_x, maxY-s.curr_y, m, color.RGBA{100, 50, 0, 255})

	buff := new(bytes.Buffer)
	if err := png.Encode(buff, img); err != nil {
		return nil, err
	}
	return buff, nil
}

func (s *Ogrid) scaledMap() (*bytes.Buffer, error) {
	m := 1 // Number of pixels per cell. m x m.
	sz := maxX / cellSz
	img := image.NewRGBA(image.Rect(0, 0, sz*m, sz*m))

	for x := 0; x < sz; x++ {
		for y := 0; y < sz; y++ {
			occ, free, total := s.checkState(x, y)
			yAdj := sz - y
			pOcc := float32(occ) / float32(total)
			pFree := float32(free) / float32(total)

			occColor := uint8(255)
			switch {
			case pOcc > 0:
				fillCell(img, x, yAdj, m, color.RGBA{occColor, 0, 0, 255})
			case pFree > 0:
				fillCell(img, x, yAdj, m, color.RGBA{0, 255, 0, 255})
			default:
				fillCell(img, x, yAdj, m, color.RGBA{20, 20, 20, 20})
			}
		}
	}

	// Set rover location on map.
	fillCell(img, s.curr_x/cellSz, sz-s.curr_y/cellSz, m, color.RGBA{100, 50, 0, 255})

	buff := new(bytes.Buffer)
	if err := png.Encode(buff, img); err != nil {
		return nil, err
	}
	return buff, nil

}

// Check the occupancy level in  cells in (x1,y1) -> (x2,y2).
// returns number of occupied, free readings and total.
func (s *Ogrid) checkState(x, y int) (int, int, int) {
	occCnt := 0
	freeCnt := 0
	totalCnt := 0

	for i := x * cellSz; i < x*cellSz+cellSz; i++ {
		for j := y * cellSz; j < y*cellSz+cellSz; j++ {
			xi := i
			yi := j
			if i >= maxX {
				xi = maxX - 1
			}
			if j >= maxY {
				yi = maxY - 1
			}
			if s.cells[xi][yi].occupied == stOccupied {
				occCnt++
				totalCnt++
			}
			if s.cells[xi][yi].occupied == stFree {
				freeCnt++
				totalCnt++
			}
		}
	}

	return occCnt, freeCnt, totalCnt
}

// fillCell renders the cell at x, y with a size of scale x scale.
func fillCell(img *image.RGBA, x int, y int, scale int, c color.RGBA) {
	for i := 0; i < scale; i++ {
		for j := 0; j < scale; j++ {
			img.Set(x*scale+i, y*scale+j, c)
		}
	}
}
