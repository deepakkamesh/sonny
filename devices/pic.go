/* Package devices provides the device driver layer for Sonny.
*  Currently the package supports the controller and Pi.
 */
package devices

import (
	"errors"
	"fmt"
	"time"

	p "github.com/deepakkamesh/sonny/protocol"
	"github.com/golang/glog"
	"github.com/tarm/serial"
)

const (
	Vdd float32 = 3.2
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

const TIMEOUT = 800 * 1000 * 1000 // Controller response timeout in nanoseconds.

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

// CMDBuffer
type CMDBuffer struct {
	req    request
	tmstmp int64
}

type Controller struct {
	port  *serial.Port
	in    chan request // Channel to recieve command.
	got   chan []byte  // Channel to recieve packet from tty.
	done  chan byte    // Channel to signal the goroutine is finished.
	quit  chan struct{}
	quitR chan struct{}
}

// NewController returns a new initialized controller.
func NewController(tty string, baud int) (*Controller, error) {

	c := &serial.Config{Name: tty, Baud: baud, ReadTimeout: 1000 * time.Millisecond}
	port, err := serialOpen(c)
	if err != nil {
		return nil, err
	}
	return &Controller{
		port:  port,
		in:    make(chan request),
		got:   make(chan []byte),
		done:  make(chan byte),
		quit:  make(chan struct{}),
		quitR: make(chan struct{}),
	}, nil
}

// Start runs the controller.
func (m *Controller) Start() {
	go m.newRun()
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
				glog.Warningf("Failed to read header from tty: %v", err)
				continue
			}
			sz := p.PacketSz(header)
			if sz == 0 {
				glog.Warning("Got a zero sized packet from controller. Discarding..")
				continue
			}
			pkt := make([]byte, sz)
			n, err = serialRead(m.port, pkt)
			if err != nil {
				glog.Warningf("Failed to read data from tty: %v", err)
				continue
			}
			if n != int(sz) {
				glog.Warningf("Expected to recieve %d bytes from tty, got %d %v", sz, n, pkt)
				continue
			}
			c := p.Checksum(header)
			if !p.VerifyChecksum(pkt, c) {
				glog.Warningf("Checksum mismatch, discarding packet: %v", p.PktPrint(pkt))
				continue
			}
			// Packet is valid. Send to the main goroutine for dispatch.
			m.got <- pkt
		}
	}
}

func (m *Controller) newRun() {
	cmdBuf := make(map[byte]*CMDBuffer)
	t := time.NewTicker(TIMEOUT * time.Nanosecond)

	for {
		glog.Flush()
		select {
		case <-m.quit:
			glog.Info("Shutting down pic processor")

		// Process command requests.
		case c := <-m.in:
			h := p.Header(c.pkt)
			dev := p.DeviceID(c.pkt[0])
			if _, ok := cmdBuf[dev]; ok {
				glog.Warningf("Device %v busy. Dropping packet.", dev)
				c.ret <- result{
					pkt: nil,
					err: errors.New(fmt.Sprintf("Error: device %v busy", dev)),
				}
				continue
			}
			cmdBuf[dev] = &CMDBuffer{
				req:    c,
				tmstmp: time.Now().UnixNano(),
			}
			time.Sleep(50 * time.Millisecond)
			glog.V(2).Infof("Sending command %v to device %v on tty", p.PktPrint(c.pkt), dev)
			if _, err := serialWrite(m.port, []byte{h}); err != nil {
				glog.Warningf("Failed to write header %v to serial", h)
			}
			if _, err := serialWrite(m.port, c.pkt); err != nil {
				glog.Warningf("Failed to write packet %v to serial", p.PktPrint(c.pkt))
			}

		// Timeout handler.
		case <-t.C:
			now := time.Now().UnixNano()
			for _, b := range cmdBuf {
				if now-b.tmstmp > TIMEOUT {
					req := b.req
					dev := p.DeviceID(req.pkt[0])
					req.ret <- result{
						pkt: nil,
						err: errors.New("timeout waiting response from controller"),
					}
					delete(cmdBuf, dev)
					glog.Warningf("Timeout controller device %v packet %v", dev, p.PktPrint(req.pkt))
				}
			}

		// TTY data handler.
		case pkt := <-m.got:
			dev := p.DeviceID(pkt[0])
			buf, ok := cmdBuf[dev]
			if !ok {
				glog.Warningf("No registered channel found for device %b, packet:%v", dev, p.PktPrint(pkt))
				continue
			}
			c := buf.req
			switch p.StatusCode(pkt[0]) {
			case p.ACK:
				glog.V(2).Infof("Recieved ACK from %v", p.PktPrint(pkt))
				continue
			case p.ACK_DONE:
				glog.V(2).Infof("Recieved ACK DONE from %v", p.PktPrint(pkt))
				c.ret <- result{
					pkt: pkt,
					err: nil,
				}
			case p.ERR:
				glog.V(2).Infof("Recieved ERR from %v", p.PktPrint(pkt))
				c.ret <- result{
					pkt: nil,
					err: p.Error(pkt[1]),
				}
			case p.DONE:
				glog.V(2).Infof("Recieved DONE from %v", p.PktPrint(pkt))
				c.ret <- result{
					pkt: pkt,
					err: nil,
				}
			}
			glog.V(2).Infof("Request fulfilled for dev %v", dev)
			delete(cmdBuf, dev)

		}
	}
}

/****** Available Functions on Controller ******/

// LEDBlink blinks the LED for duration (in ms) and for the number of times.
func (m *Controller) LEDBlink(duration uint16, times byte) error {
	ret := make(chan result)
	m.in <- request{
		pkt: []byte{p.CMD_BLINK<<4 | p.DEV_LED,
			byte(duration >> 8),
			byte(duration & 0xF),
			times,
		},

		ret: ret,
	}
	return (<-ret).err
}

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

// Motor turns the motor by turns forward if fwd is true or back if false. if
// turns = 0 stop movement, turns < 0 moves indefinitely .
func (m *Controller) Motor(turns int16, fwd bool) (error, uint16) {
	dir := p.CMD_FWD
	if !fwd {
		dir = p.CMD_BWD
	}

	ret := make(chan result)
	m.in <- request{
		pkt: []byte{dir<<4 | p.DEV_MOTOR,
			byte(turns >> 8),
			byte(turns & 0xFF),
		},
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
	ret := make(chan result)
	m.in <- request{
		pkt: []byte{cmd<<4 | p.DEV_LED},
		ret: ret,
	}
	return (<-ret).err
}

// Ping returns nil if the controller is available.
func (m *Controller) Ping() error {
	ret := make(chan result)
	m.in <- request{
		pkt: []byte{p.CMD_PING<<4 | p.DEV_ADMIN},
		ret: ret,
	}
	return (<-ret).err
}

// RotateServo rotates servo by angle.
func (m *Controller) ServoRotate(servo byte, angle int) error {
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
	ret := make(chan result)

	m.in <- request{
		pkt: []byte{
			p.CMD_ROTATE<<4 | p.DEV_SERVO,
			byte(duty >> 8),
			byte(duty & 0xFF),
			byte(period >> 8),
			byte(period & 0xFF),
			servo,
		},
		ret: ret,
	}
	return (<-ret).err // Wait for response on ack.
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
}
