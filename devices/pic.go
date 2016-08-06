package devices

import (
	"errors"
	"fmt"
	"log"
	"time"

	p "github.com/deepakkamesh/sonny/protocol"
	"github.com/tarm/serial"
)

const TIMEOUT = 500 // Controller Response Timeout.

type result struct {
	pkt []byte
	err error
}

type request struct {
	pkt []byte
	ret chan result
}

type buffer struct {
	tmstmp time.Time
	ret    chan result
	pkt    []byte
}
type Controller struct {
	port   *serial.Port
	in     chan request    // Channel to recieve command.
	out    chan []byte     // Channel to recieve response from controller.
	cmdBuf map[byte]buffer // Command buffer to maintain state of currently executing commands.
}

func NewController(tty string, baud int) *Controller {

	c := &serial.Config{Name: tty, Baud: baud}
	port, err := serial.OpenPort(c)
	if err != nil {
		log.Fatalf("Error opening tty port %v", err)
	}

	return &Controller{
		port:   port,
		in:     make(chan request),
		out:    make(chan []byte),
		cmdBuf: make(map[byte]buffer),
	}
}

func (m *Controller) Start() {
	go m.run()
	go m.read()
}

func (m *Controller) read() {

	for {
		buf := make([]byte, 16)
		_, err := m.port.Read(buf)
		if err != nil {
			log.Printf("Error reading from tty %s", err)
			continue
		}
		if !p.VerifyChecksum(buf) {
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
		// Process commands.
		case c := <-m.in:
			fmt.Printf("Executing Command: %d on device %d\n", c.pkt, p.DeviceID(c.pkt[0]))
			d := p.DeviceID(c.pkt[0])
			m.cmdBuf[d] = buffer{
				ret:    c.ret,
				tmstmp: time.Now(),
				pkt:    c.pkt,
			}

			// Send the command to the controller.
			h := p.Header(c.pkt)
			m.port.Write([]byte{h})
			m.port.Write(c.pkt)

		// Check for timeouts.
		case <-tick.C:
			for _, v := range m.cmdBuf {
				if time.Since(v.tmstmp) > TIMEOUT*time.Millisecond {
					v.ret <- result{
						pkt: nil,
						err: errors.New("Timeout on device"),
					}
					delete(m.cmdBuf, p.DeviceID(v.pkt[0]))
				}
			}

		// Process return data from controller.
		case data := <-m.out:
			deviceID := p.DeviceID(data[0])
			v, ok := m.cmdBuf[deviceID]
			if !ok {
				// TODO: add default handler here.
				log.Printf("Failed to find a handler for device: %d packet: %v", p.DeviceID(data[0]), data)
				continue
			}

			switch p.StatusCode(data[0]) {
			case p.ACK:
				continue
			case p.ACK_DONE, p.DONE:
				v.ret <- result{
					pkt: data,
					err: nil,
				}
				close(v.ret)
				delete(m.cmdBuf, deviceID)
			case p.ERR:
				v.ret <- result{
					pkt: nil,
					err: p.Error(data[1]),
				}
				close(v.ret)
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
