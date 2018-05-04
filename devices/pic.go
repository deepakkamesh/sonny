/* Package devices provides the device driver layer for Sonny.
*  Currently the package supports the controller and Pi.
 */
package devices

import (
	"container/heap"
	"errors"
	"fmt"
	"sync"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/i2c"

	p "github.com/deepakkamesh/sonny/protocol"
	"github.com/golang/glog"
)

const (
	Vdd float32 = 3.2
)

type Controller struct {
	name       string
	connector  i2c.Connector
	connection i2c.Connection
	i2c.Config
	i2cReq chan req      // channel for inbound requests for i2c.
	reqQ   PriorityQueue // Prioritized queue of requests for I2C.
	quit   chan struct{} // quit loop.
	mut    *sync.Mutex   // mutex to lock I2C bus.
}

// ret is the struct for returning response from i2c.
type ret struct {
	retData data
	err     error
}

type req struct {
	retChan  chan ret
	reqData  data
	priority int
	index    int
}

type data struct {
	deviceID byte
	pkt      []byte
}

// A PriorityQueue implements heap.Interface and holds Req structs.
type PriorityQueue []*req

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the highest, not lowest, priority so we use greater than here.
	return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*req)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// NewController returns a new initialized controller.
func NewController(a i2c.Connector, options ...func(i2c.Config)) *Controller {

	q := make(PriorityQueue, 0)
	heap.Init(&q)

	d := &Controller{
		name:      gobot.DefaultName("PIC"),
		connector: a,
		Config:    i2c.NewConfig(),
		i2cReq:    make(chan req),
		quit:      make(chan struct{}),
		reqQ:      q,
		mut:       &sync.Mutex{},
	}

	for _, option := range options {
		option(d)
	}
	return d
}

// Start initialized the PIC.
func (h *Controller) Start() (err error) {
	bus := h.GetBusOrDefault(h.connector.GetDefaultBus())
	address := h.GetAddressOrDefault(0x7)

	h.connection, err = h.connector.GetConnection(address, bus)
	if err != nil {
		return err
	}
	go h.run()
	return
}

// run starts I2C communication goroutine. The requests from channel
// are added to a priority heap and pop'ed based on decreasing priority.
// This ensures priority I2C messages (servo rotate) are send on wire firt.
func (m *Controller) run() {
	t := time.NewTicker(50 * time.Millisecond)
	for {
		select {
		case req := <-m.i2cReq:
			heap.Push(&m.reqQ, &req)

		case <-m.quit:
			return

		case <-t.C:
			if m.reqQ.Len() > 0 {
				req := heap.Pop(&m.reqQ).(*req)
				if err := m.send(req.reqData.deviceID, req.reqData.pkt); err != nil {
					req.retChan <- ret{
						err: err,
					}
					continue
				}
				pkt, err := m.recv(req.reqData.deviceID)
				if err != nil {
					req.retChan <- ret{
						err: err,
					}
					continue
				}
				req.retChan <- ret{
					retData: data{
						deviceID: req.reqData.deviceID,
						pkt:      pkt,
					},
				}
			}

		}
	}
}

func (m *Controller) get(deviceID byte, pkt []byte, priority int) ([]byte, error) {

	retChan := make(chan ret)
	m.i2cReq <- req{
		retChan: retChan,
		reqData: data{
			deviceID: deviceID,
			pkt:      pkt,
		},
		priority: priority,
	}

	retData := <-retChan
	return retData.retData.pkt, retData.err
}

// send transmits a command, parameters to deviceID.
func (m *Controller) send(deviceID byte, pkt []byte) error {
	d := []byte{p.Header(pkt)}
	d = append(d, pkt...)
	if err := m.connection.WriteBlockData(deviceID, d); err != nil {
		return err
	}
	return nil
}

func (m *Controller) recv(deviceID byte) ([]byte, error) {
	header, err := m.connection.ReadByteData(deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %v", err)
	}

	pktSz := header >> 4
	if (pktSz) == 0 {
		return nil, fmt.Errorf("got a zero sized packet. Header %b", header)
	}

	pkt := make([]byte, (pktSz))
	if _, err = m.connection.Read(pkt); err != nil {
		return nil, fmt.Errorf("failed to read packet: %v", err)
	}

	glog.V(3).Infof("Header: %v Packet size %v, content:%v", header, len(pkt), pkt)

	if !p.VerifyChecksum(pkt, p.Checksum(header)) {
		return nil, fmt.Errorf("checksum failed")
	}

	if p.StatusCode(pkt[0]) == p.ERR {
		if len(pkt) > 1 {
			return nil, p.Error(pkt[1])
		}
		return nil, fmt.Errorf("unknown error")
	}

	return pkt, nil

}

func (m *Controller) LockI2CBus()   { m.mut.Lock() }
func (m *Controller) UnlockI2CBus() { m.mut.Unlock() }

// Ping returns nil if the controller is available.
func (m *Controller) Ping() (err error) {
	if m == nil {
		return fmt.Errorf("controller not initialized")
	}

	pkt := []byte{p.CMD_PING}

	_, err = m.get(p.DEV_ADMIN, pkt, 10)
	return
}

// LEDOn turn on/off the LED.
func (m *Controller) LEDOn(on bool) (err error) {
	if m == nil {
		return fmt.Errorf("controller not initialized")
	}
	cmd := p.CMD_ON
	if !on {
		cmd = p.CMD_OFF
	}
	pkt := []byte{cmd}
	_, err = m.get(p.DEV_LED, pkt, 50)
	return
}

// LEDBlink blinks the LED for duration (in ms) and for the number of times.
func (m *Controller) LEDBlink(duration uint16, times byte) (err error) {
	if m == nil {
		return fmt.Errorf("controller not initialized")
	}
	pkt := []byte{p.CMD_BLINK,
		byte(duration >> 8),
		byte(duration & 0xF),
		times,
	}

	_, err = m.get(p.DEV_LED, pkt, 50)
	return
}

// RotateServo rotates servo by angle.
func (m *Controller) ServoRotate(servo byte, angle int) (err error) {
	if m == nil {
		return fmt.Errorf("controller not initialized")
	}
	const (
		deg0      float32 = 0.0007 // 0.7 ms.
		deg180    float32 = 0.0024 // 2.4 ms.
		pwmPeriod float32 = 0.020  // 20ms.
		// TODO: Update is pic clock is change.
		cycle float32 = 0.000001 // = Fosc/4 divided by PWM prescaler
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

	_, err = m.get(p.DEV_SERVO, pkt, 90)
	return
}

// DHT11 returns the temperature in 'C and humidity %.
func (m *Controller) DHT11() (temp, humidity uint8, err error) {
	m.mut.Lock()
	defer m.mut.Unlock()

	if m == nil {
		return 0, 0, fmt.Errorf("controller not initialized")
	}
	pkt := []byte{p.CMD_STATE}

	pkt, err = m.get(p.DEV_DHT11, pkt, 30)
	if err != nil {
		return
	}
	humidity = uint8(pkt[1])
	temp = uint8(pkt[3])
	return
}

// LDR returns the ADC light value of the LDR sensor.
func (m *Controller) LDR() (adc uint16, err error) {
	m.mut.Lock()
	defer m.mut.Unlock()

	if m == nil {
		return 0, fmt.Errorf("controller not initialized")
	}
	pkt := []byte{p.CMD_STATE}

	pkt, err = m.get(p.DEV_LDR, pkt, 10)
	if err != nil {
		return
	}
	if len(pkt) < 3 {
		err = fmt.Errorf("expected more bytes")
		return
	}
	adc = uint16(pkt[1])<<8 | uint16(pkt[2])
	return
}

// Accelerometer returns the ADC values from the accelerometer.
func (m *Controller) Accelerometer() (gx, gy, gz float32, err error) {
	m.mut.Lock()
	defer m.mut.Unlock()

	if m == nil {
		return 0, 0, 0, fmt.Errorf("controller not initialized")
	}
	pkt := []byte{p.CMD_STATE}

	pkt, err = m.get(p.DEV_ACCEL, pkt, 10)
	if err != nil {
		return
	}

	x := uint16(pkt[1])<<8 | uint16(pkt[2])
	y := uint16(pkt[3])<<8 | uint16(pkt[4])
	z := uint16(pkt[5])<<8 | uint16(pkt[6])

	gx = (float32(x)*Vdd/1023 - Vdd/2) / 0.8
	gy = (float32(y)*Vdd/1023 - Vdd/2) / 0.8
	gz = (float32(z)*Vdd/1023 - Vdd/2) / 0.8

	return
}

// BattState returns the voltage reading for the battery.
func (m *Controller) BattState() (batt float32, err error) {
	m.mut.Lock()
	defer m.mut.Unlock()

	if m == nil {
		return 0, fmt.Errorf("controller not initialized")
	}
	pkt := []byte{p.CMD_STATE}

	pkt, err = m.get(p.DEV_BATT, pkt, 10)
	if err != nil {
		return
	}

	adc := uint16(pkt[1])<<8 | uint16(pkt[2])
	batt = 2095.104 / float32(adc)
	return
}
