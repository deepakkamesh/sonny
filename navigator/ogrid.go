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

// size of the cell in centimeters. Cell is mapScale * CellSz square. The smaller the cell
// the more accurate the space calculation, but takes more memory.
// Note: Size of house is 640cm x 1676 cm. MaxX, MaxY based on sz/scale + some buffer.
const (
	mapScale   = 5
	maxX       = 150 // number of grid X coordinates.
	maxY       = 350 // number of grid Y coordinates.
	stOccupied = 1   // occupancy state of cell.
	stFree     = 0   // occupancy state of cell.
	stUnknown  = -1  // occupancy state of cell.
	DEG2RAD    = math.Pi / 180
)

// struct cell represents a single cell in the occupancy grid.
type cell struct {
	occupied int       // Is the cell occupied 1 = occupied 0 = unoccupied, -1 = unknown
	lastUpd  time.Time // Last updated timestamp for cell.
	obs      uint      // number of observations.
	posObs   float64   // Number of (positive) observations the cell is occupied.
}

// ogrid represents the occupancy grid.
type Ogrid struct {
	cells   [maxX][maxY]cell
	x       int // current x location of bot.
	y       int // current y location of bot.
	imgChan chan chan *bytes.Buffer
}

// NewOGrid returns a initialized Ogrid structure.
func NewOgrid() *Ogrid {
	return &Ogrid{
		cells:   [maxX][maxY]cell{},
		x:       120, // Roughly positioned in the end room.
		y:       320,
		imgChan: make(chan chan *bytes.Buffer),
	}
}

func (s *Ogrid) GetXY() (int, int) {
	return s.x, s.y
}

func (s *Ogrid) SetXY(x, y int) {
	s.x = x
	s.y = y
}

//ResetMap resets the grid
func (s *Ogrid) ResetMap() {
	s.x = 120
	s.y = 320
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
	s.ResetMap()
}

// Adapted from https://github.com/encukou/bresenham/blob/master/bresenham.py.
func bresenham(x0, y0, x1, y1 int) []image.Point {
	dx := x1 - x0
	dy := y1 - y0

	xsign := -1
	ysign := -1
	if dx > 0 {
		xsign = 1
	}
	if dy > 0 {
		ysign = 1
	}

	dx = int(math.Abs(float64(dx)))
	dy = int(math.Abs(float64(dy)))

	xx, xy, yx, yy := xsign, 0, 0, ysign
	if dx < dy {
		dx, dy = dy, dx
		xx, xy, yx, yy = 0, ysign, xsign, 0
	}

	D := 2*dy - dx
	y := 0

	pts := []image.Point{}
	for x := 0; x < (dx + 1); x++ {

		pts = append(pts, image.Point{x0 + x*xx + y*yx, y0 + x*xy + y*yy})
		if D >= 0 {
			y += 1
			D -= 2 * dx
		}
		D += 2 * dy
	}
	return pts
}

// pose is the deviation from x Axis.
func (s *Ogrid) UpdateMap(rangeReading []int32, startAngle int, deltaAngle int, pose float64) error {

	// From the returned beam update the cell location.
	for i := 0; i < len(rangeReading); i++ {

		servoAngle := startAngle + i*deltaAngle // relative to the robot body.
		beamAngle := pose + float64(servoAngle)
		d := float64(rangeReading[i])

		// X,Y coord of obstacle in the global frame of reference (applying scale).
		// the sign rotates the grid.
		X := math.Cos(beamAngle*DEG2RAD) * d
		Y := math.Sin(beamAngle*DEG2RAD) * d
		Xocc := int(math.Ceil(X)/mapScale) + s.x
		Yocc := int(math.Ceil(Y)/mapScale) + s.y

		// Set the free cells along the beam.
		freePoints := bresenham(s.x, s.y, Xocc, Yocc)
		for _, pt := range freePoints {
			s.cells[pt.X][pt.Y].occupied = stFree
			s.cells[pt.X][pt.Y].obs++
		}

		// Set the occupied cells.
		s.cells[Xocc][Yocc].occupied = stOccupied
		s.cells[Xocc][Yocc].obs += 1
		s.cells[Xocc][Yocc].posObs += 1

	}
	return nil
}

// GenerateMap() returns a png map of the environment.
func (s *Ogrid) GenerateMap() (*bytes.Buffer, error) {

	m := 1 // Number of pixels per cell. m x m.
	img := image.NewRGBA(image.Rect(0, 0, maxX*m, maxY*m))

	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {

			// Mark cells.
			switch s.cells[x][y].occupied {
			case stOccupied:
				fillCell(img, x, y, m, color.RGBA{200, 10, 10, 255})
			case stUnknown:
				fillCell(img, x, y, m, color.RGBA{0, 0, 0, 90})
			case stFree:
				fillCell(img, x, y, m, color.RGBA{0, 255, 0, 255})
			}
		}
	}

	// Set rover location on map.
	fillCell(img, s.x, s.y, m, color.RGBA{100, 50, 0, 255})

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

/* TODO: Delete below once the new mapping func is confirmed good.
// UpdateMap updates the occupany grid map based on lidar readings.
// minAngle is the starting angle (degrees) in reference to the bot.
// shiftAngle is the angle between readings in degrees.
func (s *Ogrid) UpdateMap1(rangeReading []int32, minAngle int, shiftAngle int, posture float64, c color.RGBA) error {

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
		x := math.Cos(effAngle*math.Pi/180) * float64(line)
		y := math.Sin(effAngle*math.Pi/180) * float64(line)
		// Calculate absolute X,Y points on map.
		xAbs := s.x + int(x)
		yAbs := s.y + int(y)

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
		if xAbs-s.x == 0 {
			continue
		}
		m := float32(yAbs-s.y) / float32(xAbs-s.x)
		b := float32(yAbs) - m*float32(xAbs)
		for j := yAbs + 1; j < s.y; j++ {
			xF := (float32(j) - float32(b)) / m
			s.cells[int(xF)][j].occupied = stFree
			s.cells[int(xF)][j].obs += 1
		}
		for j := s.y; j < yAbs-1; j++ {
			xF := (float32(j) - float32(b)) / m
			s.cells[int(xF)][j].occupied = stFree
			s.cells[int(xF)][j].obs += 1
		}
	}

	return nil
}

// GenerateMap() returns a png map of the environment. It wraps
// scaledMap or normalMap.
func (s *Ogrid) GenerateMap2() (*bytes.Buffer, error) {
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
				fillCell(img, x, yAdj, m, color.RGBA{200, 10, 10, 255})
			case stUnknown:
				fillCell(img, x, yAdj, m, color.RGBA{20, 20, 20, 10})
			case stFree:
				fillCell(img, x, yAdj, m, color.RGBA{0, 255, 0, 255})
			}

			// Set Grid lines.
			if x%mapScale == 0 || y%mapScale == 0 {
				fillCell(img, x, yAdj, m, color.RGBA{194, 194, 214, 255})
			}

		}
	}

	// Set rover location on map.
	fillCell(img, s.x, maxY-s.y, m, color.RGBA{100, 50, 0, 255})

	buff := new(bytes.Buffer)
	if err := png.Encode(buff, img); err != nil {
		return nil, err
	}
	return buff, nil
}

func (s *Ogrid) scaledMap() (*bytes.Buffer, error) {
	m := 1 // Number of pixels per cell. m x m.
	sz := maxX / mapScale
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
	fillCell(img, s.x/mapScale, sz-s.y/mapScale, m, color.RGBA{100, 50, 0, 255})

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

	for i := x * mapScale; i < x*mapScale+mapScale; i++ {
		for j := y * mapScale; j < y*mapScale+mapScale; j++ {
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
*/
