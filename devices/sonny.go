package devices

import (
	"errors"
	"fmt"
	"os"
	"time"

	roomba "github.com/deepakkamesh/go-roomba"
	"github.com/deepakkamesh/go-roomba/constants"
	"github.com/golang/glog"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/drivers/i2c"
)

// Sonny is the struct that represents all the devices.
type Sonny struct {
	*Controller                                 // PIC controller.
	*i2c.LIDARLiteDriver                        // Lidar Lite.
	*i2c.QMC5883Driver                          // Magnetometer QMC5883.
	*roomba.Roomba                              // Roomba controller.
	i2cEn                 *gpio.DirectPinDriver // GPIO port control for I2C Bus.
	*gpio.PIRMotionDriver                       // PIR driver.
	lidarEn               *gpio.DirectPinDriver // Lidar enable gpio. Pull high to disable.
	*Video
	pirState        int          // State of PIR. 1=enabled, 0=disabled.
	i2cBusState     int          // State of I2CBus. 1=enabled, 0=disabled.
	auxPowerState   int          // Start of AuxPower. 1=enabled, 0=disabled.
	roombaMode      int          // Roomba mode: 1 = passive, 2=safe, 3=full.
	auxPowerOnInit  func() error // initialization to execute after Aux Power is on.
	auxPowerOffInit func() error // initialization to execute after Aux Power is off.
}

func NewSonny(
	c *Controller,
	l *i2c.LIDARLiteDriver,
	m *i2c.QMC5883Driver,
	r *roomba.Roomba,
	i2cEn *gpio.DirectPinDriver,
	p *gpio.PIRMotionDriver,
	le *gpio.DirectPinDriver,
	v *Video,
) *Sonny {

	return &Sonny{
		c, l, m, r, i2cEn, p, le, v, 0, 0, 0, 0,
		func() error { return nil },
		func() error { return nil },
	}
}

// GetAuxPowerState returns the state of Aux Power.
func (s *Sonny) GetAuxPowerState() int {
	return s.auxPowerState
}

// SetAuxPostInit sets the initialization routines to call after aux power
// is turned on or off.
func (s *Sonny) SetAuxPostInit(fOn func() error, fOff func() error) {
	s.auxPowerOnInit = fOn
	s.auxPowerOffInit = fOff
}

// AuxPower enables/disables Auxillary power from main brush motor on Roomba.
func (s *Sonny) AuxPower(enable bool) error {
	if enable {
		if err := s.MainBrush(true, true); err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond) // Time to power up Aux.
		s.auxPowerState = 1
		if err := s.auxPowerOnInit(); err != nil {
			return err
		}
		return nil
	}

	if err := s.MainBrush(false, true); err != nil {
		return err
	}
	time.Sleep(100 * time.Millisecond) // Time to power down Aux.
	s.auxPowerState = 0
	if err := s.auxPowerOffInit(); err != nil {
		return err
	}
	return nil
}

// LidarPwrEnable enables the power to Lidar by driving the power enable pin high(on) or low(off).
func (s *Sonny) LidarPower(enable bool) error {
	if s.lidarEn == nil {
		return fmt.Errorf("lidar en pin not initialized")
	}

	if enable {
		if s.GetAuxPowerState() == 0 {
			return fmt.Errorf("Aux power not turned on; cannot power on lidar")
		}
		// Drive GPIO high to enable LIDAR.
		return s.lidarEn.DigitalWrite(1)
	}

	// Drive GPIO low to disable LIDAR.
	return s.lidarEn.DigitalWrite(0)
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

func (s *Sonny) ForwardSweep(angle, min, max int) ([]int32, error) {
	// Lock I2C bus to avoid any contention.
	s.LockI2CBus()
	defer s.UnlockI2CBus()

	if s.Controller == nil {
		return nil, errors.New("controller not initialized")
	}
	val := []int32{}

	if err := s.Controller.ServoRotate(1, 20); err != nil {
		return nil, fmt.Errorf("failed to rotate servo: %v", err)

	}
	// Sleep to allow servo to move to starting position.
	// rotation speed 100ms for 60"
	time.Sleep(400 * time.Millisecond)
	for i := min; i <= max; i += angle {
		if err := s.Controller.ServoRotate(1, i); err != nil {
			return nil, fmt.Errorf("failed to rotate servo: %v", err)
		}

		if err := s.LidarPower(true); err != nil {
			return nil, err
		}
		// Sleep to finish servo rotation prior to measuring and prevent
		// contention on I2C bus.
		time.Sleep(100 * time.Millisecond)
		// Take 3 distance measurements to eliminate any suprious readings.
		// TODO: use standard deviation to eliminate bad readings.
		dist0, err := s.Distance()
		dist1, err := s.Distance()
		dist2, err := s.Distance()
		dist := (dist0 + dist1 + dist2) / 3

		if err := s.LidarPower(false); err != nil {
			return nil, err
		}

		if err != nil {
			return nil, fmt.Errorf("failed to read lidar: %v", err)
		}
		val = append(val, int32(dist))
	}

	return val, nil
}

// I2CBusEnable enables/disables the I2C buffer chip.
// Connects the rest of the I2C devices with Pi.
func (s *Sonny) I2CBusEnable(b bool) error {
	if s.i2cEn == nil {
		return fmt.Errorf("gpio I2C not initialized")
	}
	if b {
		s.i2cBusState = 1
		return s.i2cEn.DigitalWrite(1)
	}
	s.i2cBusState = 0
	return s.i2cEn.DigitalWrite(0)
}

// GetI2CBusState return 1 if I2C bus is enabled otherwise it returns 0.
func (s *Sonny) GetI2CBusState() int {
	return s.i2cBusState
}

// CalibrateCompass runs the calibration routine on the compass and sets
// the offsets.
func (s *Sonny) CalibrateCompass() error {

	if s.QMC5883Driver == nil {
		return fmt.Errorf("Compass not enabled")
	}

	var minX, maxX, minY, maxY int16
	f, err := os.OpenFile("calibrationReading.csv", os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Reset offset.
	s.SetOffset(0, 0, 0)

	for i := 0; i < 200; i++ {
		x, y, z, err := s.RawHeading()
		if err != nil {
			return err
		}

		time.Sleep(200 * time.Millisecond)
		h := s.HeadingFromRaw(x, y, z)
		f.WriteString(fmt.Sprintf(" Reading (x,y,z,H), %v,%v,%v,%v\n", x, y, z, h))

		// Turn the bot through 360.
		if err := s.DirectDrive(50, -50); err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
		if err := s.DirectDrive(0, 0); err != nil {
			return err
		}

		if x < minX {
			minX = x
		}
		if x > maxX {
			maxX = x
		}
		if y < minY {
			minY = y
		}
		if y > maxY {
			maxY = y
		}
	}

	offX := (minX + maxX) / 2
	offY := (minY + maxY) / 2
	glog.Infof("Calibration Complete: X:(%v - %v)/2=%v, Y:(%v - %v)/2=%v", minX, maxX, offX, minY, maxY, offY)
	s.SetOffset(offX, offY, 0)
	return nil
}
