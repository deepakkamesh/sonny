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

	"github.com/golang/glog"
)

// Occupancy Grid references:
// https://github.com/udacity/RoboND-OccupancyGridMappingAlgorithm/
// https://github.com/markcsie/OccupancyGridMapping/

// size of the cell in centimeters. Cell is mapScale * CellSz square. The smaller the cell
// the more accurate the space calculation, but takes more memory.
// Note: Size of house is 700 cm x 1800 cm. MaxX, MaxY based on sz/scale + some buffer.
const (
	mapScale  = 5    // size of each cell is 5x5 cm.
	maxRealX  = 700  // Size of grid in cm.
	maxRealY  = 1800 // Siz of Grid in cm.
	maxX      = maxRealX / mapScale
	maxY      = maxRealY / mapScale
	DEG2RAD   = math.Pi / 180
	MAX_RANGE = 150 // max dist for a lidar beam in cm.

	Zmax  = 4000          // Max range of lidar sensor in cm.
	Zmin  = 0             // Min Range of lidar sensor.
	alpha = 20            // Width if obstacle in cm. ~2 ft.
	beta  = 0.5 * DEG2RAD // width of beam. Lidar ~0.5 degree.
	l0    = 0             // log odd unknown.
	locc  = 2             // log odd occupied.
	lfree = -2            // log off free.
)

// struct cell represents a single cell in the occupancy grid.
type cell struct {
	logOdd  float64   // Log Odds of the cell being occupied.
	lastUpd time.Time // Last updated timestamp for cell.
}

// ogrid represents the occupancy grid.
type Ogrid struct {
	cells [maxX][maxY]cell
	x     int // current x location of bot.
	y     int // current y location of bot.
	maxX  int
	maxY  int
}

// NewOGrid returns a initialized Ogrid structure.
func NewOgrid() *Ogrid {
	return &Ogrid{
		cells: [maxX][maxY]cell{},
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
	s.x = maxX - 210/mapScale // Roughly positioned in the end room. 200 cm from each axis.
	s.y = maxY - 210/mapScale // Roughly positioned in the end room. 200 cm from each axis.

	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			s.cells[x][y].logOdd = 0.0
		}
	}
}

// Placeholder for any goroutines to be started.
func (s *Ogrid) StartGrid() {
	s.ResetMap()
}

// pose is the deviation from x Axis.
// Older and basic without probabilities.
func (s *Ogrid) UpdateMap(rangeReading []int32, startAngle int, deltaAngle int, pose float64) error {

	// From the returned beam update the cell location.
	for i := 0; i < len(rangeReading); i++ {
		servoAngle := startAngle + i*deltaAngle // relative to the robot body.
		beamAngle := pose + float64(servoAngle)
		d := float64(rangeReading[i])

		// If lidar returns value > MAX_RANGE ignore it because we want to be able to
		// detect objects of 5cm width assuming a angle of 3-5 degrees between beams.
		// if trunc then dont set locc on the last cell.
		trunc := false
		if d > MAX_RANGE {
			d = MAX_RANGE
			trunc = true
		}

		// X,Y coord of obstacle in the global frame of reference (applying scale).
		// the sign of X,Y rotates the grid.
		X := math.Cos(beamAngle*DEG2RAD) * d
		Y := math.Sin(beamAngle*DEG2RAD) * d
		Xocc := int(math.Ceil(X/mapScale)) + s.x
		Yocc := int(math.Ceil(Y/mapScale)) + s.y

		// Set the free cells along the beam.
		freePoints := bresenham(s.x, s.y, Xocc, Yocc)
		for i, pt := range freePoints {
			if pt.X >= maxX || pt.Y >= maxY {
				glog.Warningf("X or Y exceeded grid max limits: x%v y%v got (x,y) = (%v,%v)",
					maxX, maxY, pt.X, pt.Y)
				continue
			}
			if i == len(freePoints)-1 && !trunc {
				s.cells[pt.X][pt.Y].logOdd += locc
				continue
			}
			s.cells[pt.X][pt.Y].logOdd += lfree
		}
	}
	return nil
}

// PrintMap() returns a png map of the environment.
// Occupany Grid approach.
func (s *Ogrid) PrintMap() (*bytes.Buffer, error) {

	m := 1 // Number of pixels per cell. m x m.
	img := image.NewRGBA(image.Rect(0, 0, maxX*m, maxY*m))

	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {

			p := 1 - 1/(1+math.Exp2(s.cells[x][y].logOdd))
			if p == 0.5 {
				// Unknown space.
				fillCell(img, x, y, m, color.RGBA{100, 100, 100, 90})
				continue
			}
			// cells are filled based on probabilty of being occupied.
			fillCell(img, x, y, m, color.RGBA{200, 10, 10, uint8(255 * p)})
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

/******* Utility Functions ***********/

// fillCell renders the cell at x, y with a size of scale x scale.
func fillCell(img *image.RGBA, x int, y int, scale int, c color.RGBA) {
	for i := 0; i < scale; i++ {
		for j := 0; j < scale; j++ {
			img.Set(x*scale+i, y*scale+j, c)
		}
	}
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

/* TODO: Delete below once the new mapping func is confirmed good.
/*
// UpdateOgridMap updates the occupancy grid map.
// robotTheta - Angle in degrees from robot and X axis
// sensorData - LIDAR range readings.
// startAngle - Start of LIDAR angle (deg).
// deltaAngle - Angle between each reading (deg).
func (s *Ogrid) UpdateOgridMap(sensorData []int32, startAngle int, deltaAngle int, robotTheta float64) error {

	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {

			// xi, yi center of mass of robot.
			//		xi := x*mapScale + mapScale/2 - robotOffSetX
			//		yi := -(y*mapScale + mapScale/2) + robotOffSetY
			//		r := math.Sqrt(math.Pow(float64(xi-s.x*mapScale), 2) + math.Pow(float64(yi-s.y*mapScale), 2))

			r := math.Sqrt(math.Pow(float64(x-s.x), 2) + math.Pow(float64(y-s.y), 2))

			if r <= Zmax {
				lo := s.invSensorModel(r, x, y, robotTheta*DEG2RAD, sensorData, startAngle, deltaAngle)
				fmt.Println(lo) // TODO: remove
				s.cells[x][y].logOdd += lo - l0
			}
		}
	}
	return nil
}

// inverseSensorModel for Lidar sensor.
// theta in radians.
func (s *Ogrid) invSensorModel(r float64, xi, yi int, theta float64, sensorData []int32, startAngle, deltaAngle int) float64 {

	phi := math.Atan2(float64(s.y-yi), float64(s.x-xi)) - theta

	var (
		sensorTheta float64
		minDelta    float64 = -1
		Zk          float64
		thetaK      float64
	)

	for i := 0; i < len(sensorData); i++ {
		// Convert scale from -90 to +90.
		sensorTheta = float64(startAngle+(i*deltaAngle)-90) * DEG2RAD
		if math.Abs(phi-sensorTheta) < minDelta || minDelta == -1 {
			Zk = float64(sensorData[i]) / mapScale
			thetaK = sensorTheta
			minDelta = math.Abs(phi - sensorTheta)
		}
	}

	switch {
	case r > math.Min(Zmax, Zk+alpha/2) || math.Abs(phi-thetaK) > beta/2 || Zk > Zmax || Zk < Zmin:
		return l0
	case Zk < Zmax && math.Abs(r-Zk) < alpha/2:
		return locc
	case r <= Zk:
		return lfree
	}

	// Should never come here.
	return 666.66

}

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
