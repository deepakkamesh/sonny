package devices

import (
	"errors"
	"fmt"

	roomba "github.com/deepakkamesh/go-roomba"
	"github.com/deepakkamesh/go-roomba/constants"
)

// GetRoombaTelemetry returns the current value of the roomba sensors.
func GetRoombaTelemetry(r *roomba.Roomba) (data map[byte]int16, err error) {

	if r == nil {
		return nil, errors.New("roomba not initialized")
	}

	data = make(map[byte]int16)

	d, e := r.QueryList(constants.PACKET_GROUP_100)
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
func SetRoombaMode(r *roomba.Roomba, mode byte) error {
	switch mode {
	case constants.OI_MODE_OFF:
		if err := r.Power(); err != nil {
			return err
		}
	case constants.OI_MODE_PASSIVE:
		if err := r.Passive(); err != nil {
			return err
		}
	case constants.OI_MODE_SAFE:
		if err := r.Safe(); err != nil {
			return err
		}
	case constants.OI_MODE_FULL:
		if err := r.Full(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown mode %v requested", mode)
	}

	return nil
}
