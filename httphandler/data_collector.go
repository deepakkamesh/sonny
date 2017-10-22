package httphandler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/deepakkamesh/sonny/devices"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

// dataCollector collects sensor data from the various sensors on the rover. Its runs
// as a goroutine independent of the websocket routine. This allows different intervals
// for different sensors. It also polls sensors only if there is a client connected.
func (m *Server) dataCollector() {
	t5s := time.NewTicker(5 * time.Second)
	t500ms := time.NewTicker(500 * time.Millisecond)
	t300ms := time.NewTicker(300 * time.Millisecond)
	t100ms := time.NewTicker(100 * time.Millisecond)

	// Each sensor reading is an anonymous function for readability and code flow.
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
				if !m.data.Enabled[TEMP] {
					return
				}
				if m.ctrl == nil {
					glog.V(3).Infof("Controller not initialized")
					return
				}
				t, h, err := m.ctrl.DHT11()
				if err != nil {
					glog.Warningf("Failed to read DHT11: %v", err)
					return
				}
				m.data.Controller[TEMP] = float32(t)*1.8 + 32
				m.data.Controller[HUMIDITY] = float32(h)
			}()

		case <-t500ms.C:
			// LDR sensor.
			time.Sleep(50 * time.Millisecond)
			func() {
				if !m.data.Enabled[LDR] {
					return
				}
				if m.ctrl == nil {
					return
				}
				l, err := m.ctrl.LDR()
				if err != nil {
					glog.Warningf("Failed to read LDR: %v", err)
					return
				}
				m.data.Controller[LDR] = float32(l)
			}()

			// Controller battery voltage.
			time.Sleep(50 * time.Millisecond)
			func() {
				if !m.data.Enabled[BATT] {
					return
				}
				if m.ctrl == nil {
					return
				}
				b, err := m.ctrl.BattState()
				if err != nil {
					glog.Errorf("Failed to read controller batt state: %v", err)
					return
				}
				m.data.Controller[BATT] = float32(b)
			}()

		case <-t300ms.C:
			// Compass.
			func() {
				if !m.data.Enabled[MAG] {
					return
				}
				if m.mag == nil {
					return
				}
				// TODO: to be implemented.
				//h, err := m.mag.Heading()
				h := 0
				var err error
				if err != nil {
					glog.Warningf("Failed to read Compass: %v", err)
					return
				}
				m.data.Controller[MAG] = float32(h)
			}()

			// PIR sensor.
			func() {
				if !m.data.Enabled[PIR] {
					return
				}
				if m.pir == nil {
					return
				}
				m.data.Controller[PIR] = float32(*m.pir)
			}()

		case <-t100ms.C:
			// Roomba data.
			func() {
				d, err := devices.GetRoombaTelemetry(m.roomba)
				if err != nil {
					glog.Warningf("Failed to read roomba sensors: %v", err)
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

		err = c.WriteMessage(websocket.TextMessage, jsMsg)
		if err != nil {
			glog.Errorf("Failed to write: %v", err)
			return
		}
		time.Sleep(300 * time.Millisecond)
	}
}
