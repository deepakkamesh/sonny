package navigator

import (
	"fmt"
	"math"
	"time"

	"github.com/deepakkamesh/go-roomba/constants"
	"github.com/deepakkamesh/sonny/devices"
	"github.com/golang/glog"
)

const roombaRadius float64 = 117.5 // Radius between roomba center and wheel in mm.

type AutoDrive struct {
	posture int
	*Ogrid
	sonny *devices.Sonny
}

func NewAutoDrive(s *devices.Sonny) *AutoDrive {
	return &AutoDrive{
		0,
		NewOgrid(),
		s,
	}
}

func (s *AutoDrive) TestDrive() error {
	s.Ogrid.ResetMap()
	for i := 0; i < 360; i += 90 {
		s.posture = i
		if i > 0 {
			if _, err := s.Turn(float64(90)); err != nil {
				return err
			}
		}
		if err := s.ForwardSweep(); err != nil {
			return err
		}
	}
	return nil
}

// ForwardSweep does a lidar forward sweep and updates the map.
func (s *AutoDrive) ForwardSweep() error {
	minAngle := 20
	shiftAngle := 10 // reduce.TODO

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

	return s.Ogrid.UpdateMap(rangeReading, minAngle, shiftAngle, s.posture)
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

// MoveForward moves the bot forward by number of cells. It returns the delta
// of movement in mm. If positive it overshot the location and if negative it undershot.
func (s *AutoDrive) MoveForward(cells int) (delta int, err error) {
	vel := 300 // Speed in mm/s.

	// Get starting readings.
	// TODO: Ignoring RIGHT ENcoder. Need to refactor to include both.
	encStart, err := s.getEncoderReading(constants.SENSOR_LEFT_ENCODER)
	if err != nil {
		return 0, err
	}

	desiredDist := cells * cellSz * 10               // Desired distance to move in mm.
	driveTime := float32(desiredDist) / float32(vel) // Assume fixed speed of 500 mm/s

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
