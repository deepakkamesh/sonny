package httphandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

const (
	TEMP byte = iota
	HUMIDITY
	LDR
	PIR
	MAG
	BATT
)

const (
	I2CBUS byte = iota
	AUXPOWER
)

// sensor data struct.
type sensorData struct {
	Err        string
	Roomba     map[byte]int16
	Controller map[byte]float32
	Pi         map[byte]int
	Enabled    map[byte]bool
}

// dataCollector collects sensor data from the various sensors on the rover. Its runs
// as a goroutine independent of the websocket routine. This allows different intervals
// for different sensors. It also polls sensors only if there is a client connected.
func (m *Server) dataCollector() {
	t5s := time.NewTicker(5 * time.Second)
	t1s := time.NewTicker(1000 * time.Millisecond)
	t300s := time.NewTicker(300 * time.Millisecond)

	// Each sensor reading is an anonymous function for readability (can use return) and code flow.
	for {
		// Read sensors only when there is a connected websocket.
		if m.connCount == 0 {
			time.Sleep(200 * time.Millisecond)
			continue
		}

		select {
		case <-t5s.C:
			// DHT11 sensor.
			func() {
				if m.sonny.GetI2CBusState() != 1 || m.sonny.GetAuxPowerState() != 1 {
					return
				}
				t, h, err := m.sonny.DHT11()
				if err != nil {
					glog.Warningf("Failed to read DHT11: %v", err)
					m.data.Err = fmt.Sprintf("%v\n%v", m.data.Err, err)
					return
				}
				m.data.Controller[TEMP] = float32(t)*1.8 + 32
				m.data.Controller[HUMIDITY] = float32(h)
			}()

		case <-t1s.C:
			// LDR sensor.
			func() {
				if m.sonny.GetI2CBusState() != 1 || m.sonny.GetAuxPowerState() != 1 {
					return
				}
				time.Sleep(200 * time.Millisecond)
				l, err := m.sonny.LDR()
				if err != nil {
					glog.Warningf("Failed to read LDR: %v", err)
					m.data.Err = fmt.Sprintf("%v\n%v", m.data.Err, err)
					return
				}
				m.data.Controller[LDR] = float32(l)
			}()

			// Controller battery voltage.
			func() {
				if m.sonny.GetI2CBusState() != 1 || m.sonny.GetAuxPowerState() != 1 {
					return
				}
				time.Sleep(200 * time.Millisecond)
				b, err := m.sonny.BattState()
				if err != nil {
					glog.Warningf("Failed to read controller batt state: %v", err)
					m.data.Err = fmt.Sprintf("%v\n%v", m.data.Err, err)
					return
				}
				m.data.Controller[BATT] = float32(b)
			}()

			// Compass.
			func() {
				if !m.data.Enabled[MAG] {
					return
				}
				// TODO: to be implemented.
				//h, err := m.mag.Heading()
				h := 0
				var err error
				if err != nil {
					glog.Warningf("Failed to read Compass: %v", err)
					m.data.Err = fmt.Sprintf("%v\n%v", m.data.Err, err)
					return
				}
				m.data.Controller[MAG] = float32(h)
			}()

		case <-t300s.C:
			// AuxPower State.
			m.data.Pi[AUXPOWER] = m.sonny.GetAuxPowerState()

			// I2CBus State.
			m.data.Pi[I2CBUS] = m.sonny.GetI2CBusState()

			// PIR state.
			m.data.Controller[PIR] = float32(m.sonny.GetPIRState())

			// Roomba data.
			func() {
				d, err := m.sonny.GetRoombaTelemetry()
				if err != nil {
					glog.Warningf("Failed to read roomba sensors: %v", err)
					m.data.Err = fmt.Sprintf("%v\n%v", m.data.Err, err)
					return
				}
				m.data.Roomba = d
			}()

		}
	}
}

// dataStream is the websocket server that streams rover stats.
func (m *Server) dataStream(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		glog.Warningf("failed to upgrade conn:%v", err)
		return
	}

	m.connCount++

	defer func() {
		c.Close()
		m.connCount--
	}()

	for {
		jsMsg, err := json.Marshal(m.data)
		if err != nil {
			glog.Errorf("Failed to unmarshall: %v", err)
			continue
		}
		m.data.Err = ""

		err = c.WriteMessage(websocket.TextMessage, jsMsg)
		if err != nil {
			glog.Errorf("Failed to write: %v", err)
			return
		}
		time.Sleep(300 * time.Millisecond)
	}
}
