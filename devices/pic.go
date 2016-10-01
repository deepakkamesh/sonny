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
	done       chan byte
}

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
		done:       make(chan byte),
	}, nil
}

func (m *Controller) Start() {
	go m.run()
	go m.read()
}

func (m *Controller) Stop() {
	m.quitR <- struct{}{}
	m.quit <- struct{}{}
}

func (m *Controller) read() {
	for {
		select {
		case <-m.quitR:
			return

		case dev := <-m.done:
			close(m.deviceChan[dev])
			delete(m.deviceChan, dev)

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
			// Everything looks good. Dump the packet to goroutine handling the device.
			dev := p.DeviceID(pkt[0])
			if c, ok := m.deviceChan[dev]; ok {
				c <- pkt
				continue
			}
			log.Printf("No registered channel found for device %b", dev)

		}
	}
}

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
			deviceID := p.DeviceID(c.pkt[0])
			// Device is busy. Retry later.
			if _, ok := m.deviceChan[deviceID]; ok {
				go func() {
					log.Printf("Device %v busy. Retrying command", deviceID)
					time.Sleep(time.Millisecond * 100)
					m.in <- c
				}()
				continue
			}
			m.deviceChan[deviceID] = make(chan []byte, 2)

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
						m.done <- deviceID // Send done to close channel.
						return
					// Handle data from controller.
					case d := <-m.deviceChan[deviceID]:
						switch p.StatusCode(d[0]) {
						case p.ACK:
							t.Reset(TIMEOUT * time.Millisecond)
							continue
						case p.ACK_DONE:
							c.ret <- result{
								pkt: d,
								err: nil,
							}
							m.done <- deviceID // Send done to close channel.
							return
						case p.ERR:
							c.ret <- result{
								pkt: nil,
								err: p.Error(d[1]),
							}
							m.done <- deviceID // Send done to close channel.
							return
						case p.DONE:
							c.ret <- result{
								pkt: d,
								err: nil,
							}
							m.done <- deviceID // Send done to close channel.
							return
						}
					}
				}
			}()
		}
	}
}

func (m *Controller) LedOn(on bool) error {

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
