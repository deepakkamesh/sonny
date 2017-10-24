package devices

import (
	"errors"
	"fmt"
	"time"

	roomba "github.com/deepakkamesh/go-roomba"
	"github.com/deepakkamesh/go-roomba/constants"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
)

// Sonny is the struct that represents all the devices.
type Sonny struct {
	*Controller           // PIC controller.
	*i2c.LIDARLiteDriver  // Lidar Lite.
	*i2c.HMC6352Driver    // Magnetometer HMC5663.
	*roomba.Roomba        // Roomba controller.
	*gpio.DirectPinDriver // GPIO port control for I2C Bus.
}

// GetRoombaTelemetry returns the current value of the roomba sensors.
func (s *Sonny) GetRoombaTelemetry() (data map[byte]int16, err error) {

	if s.Roomba == nil {
		return nil, errors.New("roomba not initialized")
	}

	data = make(map[byte]int16)
	d, e := s.Roomba.QueryList(constants.PACKET_GROUP_100)
	if e != nil {
		return nil, e
	}

	for i, p := range d {
		pktID := constants.PACKET_GROUP_100[i]
		if len(p) == 1 {
			data[pktID] = int16(p[0])
			continue
		}
		data[pktID] = int16(p[0])<<8 | int16(p[1])
	}
	return
}

// SetRoombaMode sets the mode for Roomba.
func (s *Sonny) SetRoombaMode(mode byte) error {
	switch mode {
	case constants.OI_MODE_OFF:
		if err := s.Roomba.Power(); err != nil {
			return err
		}
	case constants.OI_MODE_PASSIVE:
		if err := s.Roomba.Passive(); err != nil {
			return err
		}
	case constants.OI_MODE_SAFE:
		if err := s.Roomba.Safe(); err != nil {
			return err
		}
	case constants.OI_MODE_FULL:
		if err := s.Roomba.Full(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown mode %v requested", mode)
	}

	return nil
}

func (s *Sonny) ForwardSweep(angle int) ([]int32, error) {
	val := []int32{}

	// Sleep to allow servo to move to starting position.
	time.Sleep(320 * time.Millisecond)
	for i := 20; i <= 160; i += angle {
		if err := s.Controller.ServoRotate(1, i); err != nil {
			return nil, err
		}

		time.Sleep(40 * time.Millisecond) // Sleep to allow servo to finish turning.
		dist, err := s.Distance()
		if err == nil {
			break
		}

		val = append(val, int32(dist))
		time.Sleep(100 * time.Millisecond)
	}

	return val, nil
}
