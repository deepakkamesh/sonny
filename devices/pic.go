/* Package devices provides the device driver layer for Sonny.
*  Currently the package supports the controller and Pi.
 */
package devices

import (
	"errors"
	"fmt"
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

type request struct {
	pkt []byte
	ret chan result
}

type buffer struct {
	tmstmp time.Time   // Time when command was sent.
	ret    chan result // Channel to return the result.
	status byte        // Status of command.
}

type Controller struct {
	port   *serial.Port
	in     chan request    // Channel to recieve command.
	out    chan []byte     // Channel to recieve response from controller.
	cmdBuf map[byte]buffer // Command buffer to maintain state of currently executing commands.
}

func NewController(tty string, baud int) (*Controller, error) {

	c := &serial.Config{Name: tty, Baud: baud}
	port, err := serialOpen(c)

	if err != nil {
		return nil, err
	}

	return &Controller{
		port:   port,
		in:     make(chan request),
		out:    make(chan []byte),
		cmdBuf: make(map[byte]buffer),
	}, nil
}

func (m *Controller) Start() {
	go m.run()
	go m.read()
}

func (m *Controller) read() {
	header := make([]byte, 1)
	for {
		_, err := serialRead(m.port, header)
		if err != nil {
			log.Printf("Error reading header from tty: %v", err)
			continue
		}
		sz := p.PacketSz(header)
		c := p.Checksum(header)
		pkt := make([]byte, sz)
		if _, err := serialRead(m.port, pkt); err != nil {
			log.Printf("Error reading data from tty: %v", err)
			continue
		}
		if !p.VerifyChecksum(pkt, c) {
			log.Printf("Checksum mismatch, discarding packet: %v", pkt)
			continue
		}
		m.out <- pkt
		continue
	}
}

func (m *Controller) readold() {

	for {
		// TODO: This may fail if there are 2 packets within the 16 bytes.
		buf := make([]byte, 16)
		_, err := serialRead(m.port, buf)
		if err != nil {
			log.Printf("Error reading from tty %s", err)
			continue
		}
		if !p.VerifyChecksum(buf, buf[0]) {
			log.Printf("Checksum mismatch, discarding packet %v", buf)
			continue
		}
		m.out <- buf[1:] // Strip out header after checksum verification.
	}
}

func (m *Controller) run() {

	tick := time.NewTicker(500 * time.Millisecond)

	for {
		select {
		// Send commands to controller.
		case c := <-m.in:
			d := p.DeviceID(c.pkt[0])
			fmt.Printf("Executing Command: %d on device %d\n", c.pkt, d)
			m.cmdBuf[d] = buffer{
				ret:    c.ret,
				tmstmp: time.Now(),
			}

			// Send the command to the controller.
			h := p.Header(c.pkt)
			serialWrite(m.port, []byte{h})
			serialWrite(m.port, c.pkt)

		// Check for timeouts.
		case <-tick.C:
			for deviceID, buf := range m.cmdBuf {
				// timeout after ack is 20x the regular timeout.
				if (time.Since(buf.tmstmp) > TIMEOUT*time.Millisecond && buf.status != p.ACK) || (time.Since(buf.tmstmp) > 20*TIMEOUT*time.Millisecond && buf.status == p.ACK) {
					buf.ret <- result{
						pkt: nil,
						err: errors.New("Timeout on device"),
					}
					delete(m.cmdBuf, deviceID)
				}
			}

		// Process return data from controller.
		case data := <-m.out:
			deviceID := p.DeviceID(data[0])
			buf, ok := m.cmdBuf[deviceID]
			if !ok {
				// TODO: add default handler here.
				log.Printf("Failed to find a handler for device: %d packet: %v", p.DeviceID(data[0]), data)
				continue
			}

			switch p.StatusCode(data[0]) {
			case p.ACK:
				buf.status = p.ACK
			case p.ACK_DONE, p.DONE:
				buf.ret <- result{
					pkt: data,
					err: nil,
				}
				close(buf.ret)
				delete(m.cmdBuf, deviceID)
			case p.ERR:
				buf.ret <- result{
					pkt: nil,
					err: p.Error(data[1]),
				}
				close(buf.ret)
				delete(m.cmdBuf, deviceID)
			}
		}
	}
}

func (m *Controller) LedOn(on bool) error {

	log.Println("LED")
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

func (m *Controller) Ping() error {
	log.Println("Pinging controller")

	pkt := []byte{p.CMD_PING<<4 | p.DEV_ADMIN}
	ret := make(chan result)
	m.in <- request{
		pkt: pkt,
		ret: ret,
	}
	return (<-ret).err
}

func (m *Controller) RotateServo(angle int) error {
	fmt.Printf("Rotate angle %d\n", angle)

	// Assemble command data.
	pkt := []byte{p.CMD_ROTATE<<4 | p.DEV_SERVO, 10}
	// Channel for return value.
	ret := make(chan result)

	m.in <- request{
		pkt: pkt,
		ret: ret,
	}

	return (<-ret).err // Wait for response on ack.
}
