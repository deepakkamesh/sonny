/* Package devices provides the device driver layer for Sonny.
*  Currently the package supports the controller and Pi.
 */
package devices

import (
	"errors"
	"log"
	"time"

	p "github.com/deepakkamesh/sonny/protocol"
	"github.com/tarm/serial"
)

// Dependency injection for mocks.
var (
	serialOpen  = serOpen
	serialRead  = serRead
	serialWrite = serWrite
)

func serOpen(c *serial.Config) (*serial.Port, error) {
	return serial.OpenPort(c)
}

func serRead(s *serial.Port, b []byte) (int, error) {
	return s.Read(b)
}

func serWrite(s *serial.Port, b []byte) (int, error) {
	return s.Write(b)
}

const TIMEOUT = 500 // Controller response timeout in milliseconds.

// result stores the return value from the controller.
type result struct {
	pkt []byte
	err error
}

// request is the inbound command request to the controller.
type request struct {
	pkt []byte
	ret chan result
}

type Controller struct {
	port       *serial.Port
	in         chan request // Channel to recieve command.
	deviceChan map[byte]chan []byte
	quit       chan struct{}
	quitR      chan struct{}
}

// NewController returns a new initialized controller.
func NewController(tty string, baud int) (*Controller, error) {

	c := &serial.Config{Name: tty, Baud: baud, ReadTimeout: 500 * time.Millisecond}
	port, err := serialOpen(c)
	if err != nil {
		return nil, err
	}

	return &Controller{
		port:       port,
		in:         make(chan request),
		deviceChan: make(map[byte]chan []byte),
		quit:       make(chan struct{}),
		quitR:      make(chan struct{}),
	}, nil
}

// Start runs the controller.
func (m *Controller) Start() {
	go m.run()
	go m.read()
}

// Stop terminates the controller.
func (m *Controller) Stop() {
	m.quitR <- struct{}{}
	m.quit <- struct{}{}
}

// read is the loop to recieve data from the pic and sends the packet to the goroutine for processing.
func (m *Controller) read() {
	for {
		select {
		case <-m.quitR:
			return

		default:
			header := make([]byte, 1)
			n, err := serialRead(m.port, header)
			if n == 0 {
				continue
			}
			if err != nil {
				log.Printf("Failed to read header from tty: %v", err)
				continue
			}
			sz := p.PacketSz(header)
			pkt := make([]byte, sz)
			n, err = serialRead(m.port, pkt)
			if err != nil {
				log.Printf("Failed to read data from tty: %v", err)
				continue
			}
			if n != int(sz) {
				log.Printf("Expected to recieve %d bytes from tty, got %d", sz, n)
				continue
			}
			c := p.Checksum(header)
			if !p.VerifyChecksum(pkt, c) {
				log.Printf("Checksum mismatch, discarding packet: %v", pkt)
				continue
			}
			// Send the packet to goroutine handling the device.
			dev := p.DeviceID(pkt[0])
			if c, ok := m.deviceChan[dev]; ok {
				c <- pkt
				continue
			}
			log.Printf("No registered channel found for device %b, packet:%x", dev, pkt)
		}
	}
}

// run is the main loop to process input.
func (m *Controller) run() {

	for {
		select {
		case <-m.quit:
			log.Printf("Shutting down sonny")
			return

		case c := <-m.in:
			// Write command to tty.
			h := p.Header(c.pkt)
			serialWrite(m.port, []byte{h})
			serialWrite(m.port, c.pkt)
			dev := p.DeviceID(c.pkt[0])
			// Device is busy. Retry later.
			if _, ok := m.deviceChan[dev]; ok {
				go func() {
					log.Printf("Device %v busy. Retrying command", dev)
					time.Sleep(time.Millisecond * 100)
					m.in <- c
				}()
				continue
			}
			m.deviceChan[dev] = make(chan []byte)

			// Goroutine to process data from controller.
			go func() {
				t := time.NewTimer(TIMEOUT * time.Millisecond)
				for {
					select {
					// Timeout.
					case <-t.C:
						c.ret <- result{
							pkt: nil,
							err: errors.New("timeout waiting response from controller"),
						}
					// Handle data from controller.
					case d := <-m.deviceChan[dev]:
						switch p.StatusCode(d[0]) {
						case p.ACK:
							t.Reset(TIMEOUT * time.Millisecond)
							continue
						case p.ACK_DONE:
							c.ret <- result{
								pkt: d,
								err: nil,
							}
						case p.ERR:
							c.ret <- result{
								pkt: nil,
								err: p.Error(d[1]),
							}
						case p.DONE:
							c.ret <- result{
								pkt: d,
								err: nil,
							}
						}
					}
					// Done with this channel.
					close(m.deviceChan[dev])
					delete(m.deviceChan, dev)
					return
				}
			}()
		}
	}
}

// LEDBlink blinks the LED for duration (in ms) and for the number of times.
func (m *Controller) LEDBlink(duration uint16, times byte) error {
	pkt := []byte{p.CMD_BLINK<<4 | p.DEV_LED, byte(duration >> 8), byte(duration & 0xF), times}
	ret := make(chan result)
	m.in <- request{
		pkt: pkt,
		ret: ret,
	}
	return (<-ret).err
}

// LDR returns the ADC light value of the LDR sensor.
func (m *Controller) LDR() (error, uint16) {
	pkt := []byte{p.CMD_STATE<<4 | p.DEV_LDR}
	ret := make(chan result)
	m.in <- request{
		pkt: pkt,
		ret: ret,
	}
	res := <-ret
	if res.err != nil {
		return res.err, 0
	}
	return nil, uint16(res.pkt[0]<<8 | res.pkt[1])
}

// Motor turns the motor by turns forward if fwd is true or back if false.
func (m *Controller) Motor(turns uint16, fwd bool) (error, uint16) {
	dir := p.CMD_FWD
	if !fwd {
		dir = p.CMD_BWD
	}

	pkt := []byte{dir<<4 | p.DEV_MOTOR}
	ret := make(chan result)
	m.in <- request{
		pkt: pkt,
		ret: ret,
	}
	return (<-ret).err, 0
}

// LEDOn turn on/off the LED.
func (m *Controller) LEDOn(on bool) error {
	cmd := p.CMD_ON
	if !on {
		cmd = p.CMD_OFF
	}
	pkt := []byte{cmd<<4 | p.DEV_LED}
	ret := make(chan result)
	m.in <- request{
		pkt: pkt,
		ret: ret,
	}
	return (<-ret).err
}

// Ping returns nil if the controller is available.
func (m *Controller) Ping() error {
	pkt := []byte{p.CMD_PING<<4 | p.DEV_ADMIN}
	ret := make(chan result)
	m.in <- request{
		pkt: pkt,
		ret: ret,
	}
	return (<-ret).err
}

// RotateServo rotates servo by angle.
func (m *Controller) ServoRotate(servo byte, angle byte) error {
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

	// Assemble command data.
	pkt := []byte{
		p.CMD_ROTATE<<4 | p.DEV_SERVO,
		byte(duty >> 8),
		byte(duty & 0xFF),
		byte(period >> 8),
		byte(period & 0xFF),
		servo,
	}
	ret := make(chan result)

	m.in <- request{
		pkt: pkt,
		ret: ret,
	}
	return (<-ret).err // Wait for response on ack.
}
