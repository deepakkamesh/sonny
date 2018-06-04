package navigator

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/deepakkamesh/go-roomba/constants"
	"github.com/deepakkamesh/sonny/devices"
	"github.com/golang/glog"
)

const roombaRadius float64 = 117.5 // Radius between roomba center and wheel in mm.

type AutoDrive struct {
	*Ogrid
	sonny *devices.Sonny
}

func NewAutoDrive(s *devices.Sonny) *AutoDrive {
	return &AutoDrive{
		NewOgrid(),
		s,
	}
}

func (s *AutoDrive) Start() error {
	s.StartGrid()
	return nil
}

func (s *AutoDrive) TestDrive() error {
	//	s.Ogrid.ResetMap()

	for _, c := range []int{0, 200} {

		d, err := s.MoveForward(c)
		if err != nil {
			return err
		}
		effDist := (c + d) / 10
		posture, err := s.sonny.Heading()
		if err != nil {
			return err
		}
		effAngle := 360 - posture
		glog.Infof("EffAngle %v EffDist %v", effAngle, effDist)

		x := (math.Cos(effAngle*math.Pi/180) * float64(effDist))
		y := (math.Sin(effAngle*math.Pi/180) * float64(effDist))

		s.curr_x += int(x)
		s.curr_y += int(y)

		glog.Infof("X:%v,Y:%v", s.curr_x, s.curr_y)

		// Take readings at different angles.
		for _, i := range []int{0, -60, -50, -30} {
			if _, err := s.Turn(float64(i)); err != nil {
				return err
			}
			if err := s.ForwardSweep(); err != nil {
				return err
			}
		}

		// turn back to original position.
		if _, err := s.Turn(140); err != nil {
			return err
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

// ForwardSweep does a lidar forward sweep and updates the map.
func (s *AutoDrive) ForwardSweep() error {
	minAngle := 20
	shiftAngle := 5 // reduce.TODO

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
		time.Sleep(100 * time.Millisecond)
	}
	if failure {
		return fmt.Errorf("Failed to update map: %v", err)
	}

	posture, err := s.sonny.Heading()
	if err != nil {
		return err
	}
	//glog.Infof("Posture: %v Reading:%v", posture, rangeReading)
	prettyPrint(s.curr_x, s.curr_y, posture, rangeReading)
	return s.Ogrid.UpdateMap(rangeReading, minAngle, shiftAngle, 360-posture, color.RGBA{255, 0, 0, 20})
}

// getEncoderReading returns the current encoder reading.
func (s *AutoDrive) getEncoderReading(enc byte) (int16, error) {
	p, err := s.sonny.Sensors(enc)
	if err != nil {
		return 0, err
	}
	return int16(p[0])<<8 | int16(p[1]), nil
}

// Turn rotates the bot by angle in degrees and returns the delta in degrees.
// positve angle rotates clockwise.
func (s *AutoDrive) Turn(angle float64) (int, error) {
	if angle == 0 {
		return 0, nil
	}
	vel := 50 //speed in mm/s.

	// Get starting encoder reading.
	encStart, err := s.getEncoderReading(constants.SENSOR_LEFT_ENCODER)
	if err != nil {
		return 0, err
	}

	// Calculate cirumference to drive. θ radian = circumfrence/radius.
	c := roombaRadius * math.Abs(angle) * math.Pi / 180
	driveTime := float32(c) / float32(vel)

	glog.V(3).Infof("driveTime(s)%v Circumfrence(mm):%v Angle:%v", driveTime, c, angle)

	rvel := -1 * vel
	lvel := vel
	if angle < 0 {
		rvel = vel
		lvel = -1 * vel
	}
	// Turn the bot.
	if err := s.sonny.DirectDrive(int16(rvel), int16(lvel)); err != nil {
		return 0, err
	}
	time.Sleep(time.Duration(driveTime*1000) * time.Millisecond)
	if err := s.sonny.DirectDrive(0, 0); err != nil {
		return 0, err
	}

	// Calculate if we overshot or undershot landing and return delta.
	encEnd, err := s.getEncoderReading(constants.SENSOR_LEFT_ENCODER)
	if err != nil {
		return 0, err
	}

	// TODO: need to return the right delta from compass.
	_ = encStart
	_ = encEnd
	return 0, nil
}

// MoveForward moves the bot forward in mm  . It returns the delta
// of movement in mm. If positive it overshot the location and if negative it undershot.
func (s *AutoDrive) MoveForward(desiredDist int) (delta int, err error) {
	if desiredDist == 0 {
		return 0, nil
	}
	vel := 300 // Speed in mm/s.

	// Get starting readings.
	// TODO: Ignoring RIGHT ENcoder. Need to refactor to include both.
	encStart, err := s.getEncoderReading(constants.SENSOR_LEFT_ENCODER)
	if err != nil {
		return 0, err
	}

	driveTime := float32(desiredDist) / float32(vel) // Assume fixed speed mm/s

	// Move bot (time*speed).
	if err := s.sonny.DirectDrive(int16(vel), int16(vel)); err != nil {
		return 0, err
	}
	time.Sleep(time.Duration(driveTime*1000) * time.Millisecond)
	if err := s.sonny.DirectDrive(0, 0); err != nil {
		return 0, err
	}

	// Calculate if we overshot or undershot landing and return delta.
	encEnd, err := s.getEncoderReading(constants.SENSOR_LEFT_ENCODER)
	if err != nil {
		return 0, err
	}

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
