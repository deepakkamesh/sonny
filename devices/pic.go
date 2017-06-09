/* Package devices provides the device driver layer for Sonny.
*  Currently the package supports the controller and Pi.
 */
package devices

import (
	"errors"
	"fmt"

	p "github.com/deepakkamesh/sonny/protocol"
	"github.com/kidoman/embd"
)

const (
	Vdd float32 = 3.2
)

type Controller struct {
	bus     embd.I2CBus // Initialized I2C bus.
	address byte        // I2c Address of the pic controller.
}

// NewController returns a new initialized controller.
func NewController(bus embd.I2CBus, address byte) *Controller {
	return &Controller{
		bus:     bus,
		address: address,
	}
}

// sendCmd sends a command, parameters to deviceID.
func (m *Controller) send(deviceID byte, pkt []byte) error {
	d := []byte{p.Header(pkt)}
	d = append(d, pkt...)
	if err := m.bus.WriteToReg(m.address, deviceID, d); err != nil {
		return err
	}
	return nil
}

func (m *Controller) recv(deviceID byte) ([]byte, error) {

	header, err := m.bus.ReadByteFromReg(0x07, deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %v", err)
	}

	if (header >> 4) == 0 {
		return nil, fmt.Errorf("got a zero sized header. Dropping quietly")
	}
	pkt := make([]byte, (header >> 4))

	if err = m.bus.ReadFromReg(m.address, deviceID, pkt); err != nil {
		return nil, fmt.Errorf("failed to read packet: %v", err)
	}
	if !p.VerifyChecksum(pkt, p.Checksum(header)) {
		return nil, fmt.Errorf("checksum failed")
	}
	if p.StatusCode(pkt[0]) == p.ERR {
		return nil, p.Error(pkt[1])
	}
	return pkt, nil
}

// Ping returns nil if the controller is available.
func (m *Controller) Ping() (err error) {
	pkt := []byte{p.CMD_PING}

	if err := m.send(p.DEV_ADMIN, pkt); err != nil {
		return fmt.Errorf("unable to send command: %v", err)
	}
	_, err = m.recv(p.DEV_ADMIN)
	return
}

// LEDOn turn on/off the LED.
func (m *Controller) LEDOn(on bool) (err error) {
	cmd := p.CMD_ON
	if !on {
		cmd = p.CMD_OFF
	}
	pkt := []byte{cmd}
	if err := m.send(p.DEV_LED, pkt); err != nil {
		return fmt.Errorf("unable to send command: %v", err)
	}
	_, err = m.recv(p.DEV_LED)
	return
}

// LEDBlink blinks the LED for duration (in ms) and for the number of times.
func (m *Controller) LEDBlink(duration uint16, times byte) (err error) {
	pkt := []byte{p.CMD_BLINK,
		byte(duration >> 8),
		byte(duration & 0xF),
		times,
	}

	if err := m.send(p.DEV_LED, pkt); err != nil {
		return fmt.Errorf("unable to send command: %v", err)
	}
	_, err = m.recv(p.DEV_LED)
	return
}

// RotateServo rotates servo by angle.
func (m *Controller) ServoRotate(servo byte, angle int) (err error) {
	const (
		deg0      float32 = 0.0007    // 0.7 ms.
		deg180    float32 = 0.0024    // 2.4 ms.
		pwmPeriod float32 = 0.020     // 20ms.
		cycle     float32 = 0.0000005 // = Fosc/4 divided by PWM prescaler.
	)

	// Ensure maximums are not exceeded.
	if angle < 0 || angle > 180 {
		return errors.New("Angle needs to be between 0 to 180 degrees")
	}
	if servo != 1 && servo != 2 {
		return errors.New("Servo should be 1 or 2")
	}

	time := deg0 + ((deg180 - deg0) * float32(angle) / 180)
	duty := uint16(time / cycle)        // On time.
	period := uint16(pwmPeriod / cycle) // PWM period.

	pkt := []byte{
		p.CMD_ROTATE,
		byte(duty >> 8),
		byte(duty & 0xFF),
		byte(period >> 8),
		byte(period & 0xFF),
		servo,
	}

	if err := m.send(p.DEV_SERVO, pkt); err != nil {
		return fmt.Errorf("unable to send command: %v", err)
	}
	_, err = m.recv(p.DEV_SERVO)
	return
}

/*
// DHT11 returns the temperature in 'C and humidity %.
func (m *Controller) DHT11() (temp, humidity uint8, err error) {
	ret := make(chan result)
	m.in <- request{
		pkt: []byte{p.CMD_STATE<<4 | p.DEV_DHT11},
		ret: ret,
	}
	res := <-ret
	if err = res.err; err != nil {
		return
	}

	humidity = uint8(res.pkt[1])
	temp = uint8(res.pkt[3])
	return
}

// LDR returns the ADC light value of the LDR sensor.
func (m *Controller) LDR() (adc uint16, err error) {
	ret := make(chan result)
	m.in <- request{
		pkt: []byte{p.CMD_STATE<<4 | p.DEV_LDR},

		ret: ret,
	}
	res := <-ret
	if err = res.err; err != nil {
		return
	}

	adc = uint16(res.pkt[1])
	adc = adc<<8 | uint16(res.pkt[2])
	return
}

// Accelerometer returns the ADC values from the accelerometer.
func (m *Controller) Accelerometer() (gx, gy, gz float32, err error) {

	ret := make(chan result)
	m.in <- request{
		pkt: []byte{p.CMD_STATE<<4 | p.DEV_ACCEL},
		ret: ret,
	}

	res := <-ret
	if err = res.err; err != nil {
		return
	}

	x := uint16(res.pkt[1])
	x = x<<8 | uint16(res.pkt[2])
	y := uint16(res.pkt[3])
	y = y<<8 | uint16(res.pkt[4])
	z := uint16(res.pkt[5])
	z = z<<8 | uint16(res.pkt[6])

	gx = (float32(x)*Vdd/1023 - Vdd/2) / 0.8
	gy = (float32(y)*Vdd/1023 - Vdd/2) / 0.8
	gz = (float32(z)*Vdd/1023 - Vdd/2) / 0.8

	return
}

// BattState returns the voltage reading for the battery.
func (m *Controller) BattState() (float32, error) {
	ret := make(chan result)
	m.in <- request{
		pkt: []byte{p.CMD_STATE<<4 | p.DEV_BATT},
		ret: ret,
	}
	res := <-ret
	if res.err != nil {
		return 0, res.err
	}

	var adc uint16
	adc = uint16(res.pkt[1])
	adc = adc<<8 | uint16(res.pkt[2])
	return 2095.104 / float32(adc), nil
}

// Distance returns the distance reading from the ultrasonic sensor.
func (m *Controller) Distance() (uint16, error) {
	ret := make(chan result)
	m.in <- request{
		pkt: []byte{p.CMD_STATE<<4 | p.DEV_US020},
		ret: ret,
	}
	res := <-ret
	if res.err != nil {
		return 0, res.err
	}

	var dist uint16
	dist = uint16(res.pkt[1])
	dist = dist<<8 | uint16(res.pkt[2])
	return dist, nil
}*/
