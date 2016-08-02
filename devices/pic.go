package devices

import (
	"fmt"
	"log"

	p "github.com/deepakkamesh/sonny/protocol"
	"github.com/tarm/serial"
)

type retVal struct {
	data []byte
	err  error
}
type command struct {
	data []byte
	ret  chan retVal
}

type controller struct {
	port   *serial.Port
	in     chan command         // Channel to recieve command.
	out    chan []byte          // Channel to recieve response from controller.
	cmdBuf map[byte]chan retVal // Command buffer to maintain state of currently executing commands.
}

func NewController(tty string, baud int) *controller {

	c := &serial.Config{Name: tty, Baud: baud}
	port, err := serial.OpenPort(c)
	if err != nil {
		log.Fatalf("Error opening tty port %v", err)
	}

	return &controller{
		port:   port,
		in:     make(chan command),
		out:    make(chan []byte),
		cmdBuf: make(map[byte]chan retVal),
	}
}

func (m *controller) Start() {
	go m.run()
	go m.read()
}

func (m *controller) read() {

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

func (m *controller) run() {
	for {
		select {
		// Process commands.
		case c := <-m.in:
			fmt.Printf("Executing Command: %d on device %d\n", c.data, p.DeviceID(c.data[0]))
			m.cmdBuf[p.DeviceID(c.data[0])] = c.ret
			h := p.Header(c.data)
			m.port.Write([]byte{h})
			m.port.Write(c.data)

		// Process return data from controller.
		// TODO: Add timeout for command buffer.
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
				v <- retVal{
					data: data,
					err:  nil,
				}
				close(v)
				delete(m.cmdBuf, deviceID)
			case p.ERR:
				v <- retVal{
					data: nil,
					err:  p.Error(data[1]),
				}
				close(v)
				delete(m.cmdBuf, deviceID)
			}
		}
	}
}
func (m *controller) LedOn(on bool) error {

	log.Println("LED")
	cmd := p.CMD_ON
	if !on {
		cmd = p.CMD_OFF
	}
	data := []byte{cmd<<4 | p.DEV_LED}
	ret := make(chan retVal)
	m.in <- command{
		data: data,
		ret:  ret,
	}
	return (<-ret).err
}

func (m *controller) Ping() error {
	log.Println("Pinging controller")

	cmd := []byte{p.CMD_PING<<4 | p.DEV_ADMIN}
	ret := make(chan retVal)
	m.in <- command{
		data: cmd,
		ret:  ret,
	}
	return (<-ret).err
}

func (m *controller) RotateServo(angle int) error {
	fmt.Printf("Rotate angle %d\n", angle)

	// Assemble command data.
	cmd := []byte{p.CMD_ROTATE<<4 | p.DEV_SERVO, 10}
	// Channel for return value.
	ret := make(chan retVal)

	m.in <- command{
		data: cmd,
		ret:  ret,
	}

	return (<-ret).err // Wait for response on ack.
}
