package navigator

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/deepakkamesh/sonny/devices"
	"github.com/golang/glog"
)

type AutoDrive struct {
	*Ogrid
	sonny devices.Platform
}

func NewAutoDrive(s devices.Platform) *AutoDrive {
	return &AutoDrive{
		NewOgrid(),
		s,
	}
}

func (s *AutoDrive) Start() error {
	s.StartGrid()
	s.ResetMap()
	return nil
}

func (s *AutoDrive) TestDrive(cmd string, val string) error {

	switch cmd {

	case "sweep":
		if err := s.ForwardSweep(3, 60, 120); err != nil {
			return err
		}

	case "reset":
		s.Ogrid.ResetMap()
		s.SetXY(500, 500)
		glog.Infof("X:%v,Y:%v", s.x, s.y)

	case "fwd":
		_, _ = strconv.Atoi(val) // dist in cm.
	case "turn":
		angle, _ := strconv.Atoi(val) // dist in cm.
		glog.Infof("turn by :%v", angle)
		if _, err := s.sonny.Turn(float64(angle)); err != nil {
			return err
		}

	case "calib":
		test := false
		if val == "test" {
			test = true
		}
		if err := s.sonny.CalibrateCompass(test); err != nil {
			glog.Errorf("err %v", err)
		}

	}
	return nil
}

func prettyPrint(x, y int, posture float64, reading []int32) {
	fmt.Printf("\ndatapt{\nposture:%v,\nx:%v,\ny:%v,\nreading:[]int32{", posture, x, y)
	for _, i := range reading {
		fmt.Printf("%v,", i)
	}
	fmt.Printf("},\n},\n")
}

// Move moves the rover forward or backward and updates the current location on the grid.
// Dist to move is provided in mm. Positive moves forward. vel is mm/s.
func (s *AutoDrive) Move(distmm, vel int) error {
	orient, err := s.sonny.TiltHeading()
	if err != nil {
		return err
	}

	if vel < 0 {
		orient += 180
	}

	dist, err := s.sonny.MoveForward(distmm, vel)
	if err != nil {
		return err
	}
	// Some sleep to let rover come to a stop.
	time.Sleep(500 * time.Millisecond)
	X := math.Cos(orient*DEG2RAD) * dist
	Y := math.Sin(orient*DEG2RAD) * dist

	x := int(math.Ceil(X)/mapScale) + s.x
	y := int(math.Ceil(Y)/mapScale) + s.y

	s.SetXY(x, y)

	glog.Infof("Posture: %0.2f X:%v,Y:%v", orient, s.x+int(x), s.y+int(y))

	return nil
}

// ForwardSweep does a lidar forward sweep and updates the map.
func (s *AutoDrive) ForwardSweep(deltaAngle, minAngle, maxAngle int) error {

	// Get forward sweep radar readings with retry in case of I2C failures.
	var (
		err          error
		rangeReading []int32
		failure      bool = true
	)
	posture, err := s.sonny.TiltHeading()
	if err != nil {
		return err
	}
	glog.Infof("Current angle: %v", posture)

	for i := 0; i < 3; i++ {
		rangeReading, err = s.sonny.ForwardSweep(deltaAngle, minAngle, maxAngle)
		if err == nil {
			failure = false
			break
		}
		glog.Warningf("Forward sweep failed retry#%v: %v", i, err)
		time.Sleep(100 * time.Millisecond)
	}
	if failure {
		return fmt.Errorf("Failed to update map: %v", err)
	}

	return s.Ogrid.UpdateMap(rangeReading, minAngle, deltaAngle, posture-90)
}
