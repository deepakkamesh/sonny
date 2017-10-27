package devices

import (
	"errors"
	"fmt"
	"time"

	roomba "github.com/deepakkamesh/go-roomba"
	"github.com/deepakkamesh/go-roomba/constants"
	"github.com/golang/glog"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
)

// Sonny is the struct that represents all the devices.
type Sonny struct {
	*Controller               // PIC controller.
	*i2c.LIDARLiteDriver      // Lidar Lite.
	*i2c.HMC6352Driver        // Magnetometer HMC5663.
	*roomba.Roomba            // Roomba controller.
	*gpio.DirectPinDriver     // GPIO port control for I2C Bus.
	*gpio.PIRMotionDriver     // PIR driver.
	pirState              int // State of PIR. 1=enabled, 0=disabled.
	i2cBusState           int // State of I2CBus. 1=enabled, 0=disabled.
	auxPowerState         int // Start of AuxPower. 1=enabled, 0=disabled.
	roombaMode            int //Roomba mode: 1 = passive, 2=safe, 3=full.
}

func NewSonny(c *Controller,
	l *i2c.LIDARLiteDriver,
	m *i2c.HMC6352Driver,
	r *roomba.Roomba,
	i2cEn *gpio.DirectPinDriver,
	p *gpio.PIRMotionDriver) *Sonny {

	return &Sonny{
		c, l, m, r, i2cEn, p, 0, 0, 0, 0,
	}
}

// GetAuxPowerState returns the state of Aux Power.
func (s *Sonny) GetAuxPowerState() int {
	return s.auxPowerState
}

// AuxPower enables/disables Auxillary power from main brush motor on Roomba.
func (s *Sonny) AuxPower(enable bool) error {
	if enable {
		s.auxPowerState = 1
		return s.MainBrush(true, true)
	}

	s.auxPowerState = 0
	return s.MainBrush(false, true)
}

// PIREventLoop subscribes to events from the PIR gpio.
func (s *Sonny) PIREventLoop() {
	if s.PIRMotionDriver == nil {
		return
	}

	pirCh := s.PIRMotionDriver.Subscribe()
	go func() {
		for {
			evt := <-pirCh
			s.pirState = evt.Data.(int)
			glog.V(3).Infof("Got pir data %v %v", evt.Name, evt.Data.(int))
		}
	}()
}

// Returns PIR state.
func (s *Sonny) GetPIRState() int {
	return s.pirState
}

// GetRoombaTelemetry returns the current value of the roomba sensors.
func (s *Sonny) GetRoombaTelemetry() (data map[byte]int16, err error) {

	if s.Roomba == nil {
		return nil, fmt.Errorf("roomba not initialized")
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

	// Inspect roomba mode. If different, reset aux power.
	prevMode := s.roombaMode
	s.roombaMode = int(data[constants.SENSOR_OI_MODE])
	// Changed into passive mode.
	if s.roombaMode != prevMode && s.roombaMode == 1 {
		s.AuxPower(false)
	}
	return
}

// GetRoombaMode returns the current roomba mode from the sensor reading.
func (s *Sonny) GetRoombaMode() int {
	return s.roombaMode
}

// SetRoombaMode sets the mode for Roomba.
func (s *Sonny) SetRoombaMode(mode byte) error {
	if s.Roomba == nil {
		return fmt.Errorf("roomba not initialized")
	}
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
	if s.Controller == nil {
		return nil, errors.New("controller not initialized")
	}
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

// I2CBusEnable enables/disables the I2C buffer chip.
// Connects the rest of the I2C devices with Pi.
func (s *Sonny) I2CBusEnable(b bool) error {
	if s.DirectPinDriver == nil {
		return fmt.Errorf("gpio I2C not initialized")
	}
	if b {
		s.i2cBusState = 1
		return s.DigitalWrite(1)
	}
	s.i2cBusState = 0
	return s.DigitalWrite(0)
}

// GetI2CBusState return 1 if I2C bus is enabled otherwise it returns 0.
func (s *Sonny) GetI2CBusState() int {
	return s.i2cBusState
}
