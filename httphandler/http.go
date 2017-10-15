package httphandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gobot.io/x/gobot/drivers/i2c"

	roomba "github.com/deepakkamesh/go-roomba"
	"github.com/deepakkamesh/sonny/devices"
	"github.com/deepakkamesh/sonny/rpc"
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

type Server struct {
	ctrl       *devices.Controller
	lidar      *i2c.LIDARLiteDriver
	mag        *i2c.HMC6352Driver
	pir        *int
	roomba     *roomba.Roomba
	ssl        bool
	resources  string
	servoAngle map[byte]int // Map to hold state of each servo.
	servoDelta uint8
	velocity   int16
	timer      *time.Timer
	data       *sensorData
	connCount  int // Count of websockets.
}

// Struct to return JSON.
type response struct {
	Err  string
	Data interface{}
}

// sensor data struct.
type sensorData struct {
	Err        string
	Roomba     map[byte]int16
	Controller map[byte]float32
	Enabled    map[byte]bool
}

func New(d *rpc.Devices, ssl bool, resources string) *Server {
	t := time.NewTimer(500 * time.Millisecond)
	t.Stop()

	return &Server{
		ctrl:       d.Ctrl,
		lidar:      d.Lidar,
		mag:        d.Mag,
		pir:        d.Pir,
		roomba:     d.Roomba,
		ssl:        ssl,
		resources:  resources,
		servoAngle: map[byte]int{1: 90, 2: 90},
		servoDelta: 10,
		velocity:   100,
		timer:      t,
		data: &sensorData{
			Controller: make(map[byte]float32),
			Roomba:     make(map[byte]int16),
			Enabled: map[byte]bool{
				TEMP:     true,
				HUMIDITY: true,
				LDR:      true,
				PIR:      true,
				MAG:      true,
				BATT:     true,
			},
		},
	}

}

func (m *Server) Start(hostPort string) error {

	http.HandleFunc("/api/setparam/", m.SetParam)
	http.HandleFunc("/api/ping/", m.Ping)
	http.HandleFunc("/api/ledon/", m.LEDOn)
	http.HandleFunc("/api/ledblink/", m.LEDBlink)
	http.HandleFunc("/api/servorotate/", m.ServoRotate)
	http.HandleFunc("/api/move/", m.Move)
	http.HandleFunc("/api/roomba_cmd/", m.RoombaCmd)
	http.HandleFunc("/datastream", m.dataStream) // Websocket  data stream Handler.

	// Serve static content from resources dir.
	fs := http.FileServer(http.Dir(m.resources))
	http.Handle("/", fs)
	go m.dataCollector()
	return http.ListenAndServe(hostPort, nil)
}

// writeResponse writes the response json object to w. If unable to marshal
// it writes a http 500.
func writeResponse(w http.ResponseWriter, resp *response) {
	js, e := json.Marshal(resp)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	glog.V(2).Infof("Writing json response %s", js)
	w.Write(js)
}

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

			time.Sleep(50 * time.Millisecond)
			// LDR sensor.
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
			time.Sleep(50 * time.Millisecond)

			// Controller battery voltage.
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

// Ping is a http wrapper for devices.Ping.
func (m *Server) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if m.ctrl == nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: Controller not initialized"),
		})
		return
	}

	if err := m.ctrl.Ping(); err != nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: ping failed %v", err),
		})
		return
	}

	writeResponse(w, &response{
		Data: "OK",
	})
}

// SetParam sets http console params.
func (m *Server) SetParam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	if v := r.Form.Get("servoDelta"); v != "" {
		val, _ := strconv.ParseUint(v, 10, 8)
		m.servoDelta = uint8(val)
	}

	if v := r.Form.Get("velocity"); v != "" {
		val, _ := strconv.ParseInt(v, 10, 16)
		m.velocity = int16(val)
	}

	writeResponse(w, &response{
		Err:  "",
		Data: "OK",
	})
}

// Move is the wrapper around ctrl.Move.
func (m *Server) Move(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if m.roomba == nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: roomba not initialized"),
		})
		return
	}

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	dir := strings.ToLower(r.Form.Get("dir")) // Motor button { up, down, left, right}
	var err error

	switch dir {
	case "fwd":
		err = m.roomba.Drive(m.velocity, 32767)

	case "bwd":
		err = m.roomba.Drive(-1*m.velocity, 32767)

	case "right":
		err = m.roomba.Drive(m.velocity, -1)

	case "left":
		err = m.roomba.Drive(m.velocity, 1)
	}

	if err != nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: motor failed %v", err),
		})
		glog.Errorf("Failed to run motor: %v", err)
		return
	}

	if run := m.timer.Reset(500 * time.Millisecond); !run {
		go func() {
			<-m.timer.C
			m.roomba.Drive(0, 0)
		}()
	}
}

// LEDBlink is the http wrapper for devices.LEDBlink().
func (m *Server) LEDBlink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if m.ctrl == nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: Controller not initialized"),
		})
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dur, err := strconv.ParseUint(strings.ToLower(r.Form.Get("duration")), 10, 16) // Duration of blink in ms.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	no, err := strconv.ParseUint(strings.ToLower(r.Form.Get("times")), 10, 8) // Number of times to blink.
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := m.ctrl.LEDBlink(uint16(dur), byte(no)); err != nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: failed to blink LED %v", err),
		})
		return
	}

	writeResponse(w, &response{
		Data: "OK",
	})

}

// LEDOn is the http wrapper for devices.LEDOn().
func (m *Server) LEDOn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if m.ctrl == nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: Controller not initialized"),
		})
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a := strings.ToLower(r.Form.Get("cmd")) // valid values = {on,off}.

	switch a {
	case "on":
		if err := m.ctrl.LEDOn(true); err != nil {
			writeResponse(w, &response{
				Err: fmt.Sprintf("Error: LED failed %v", err),
			})
			return
		}
	case "off":
		if err := m.ctrl.LEDOn(false); err != nil {
			writeResponse(w, &response{
				Err: fmt.Sprintf("Error: LED failed  %v", err),
			})
			return
		}
	default:
		writeResponse(w, &response{
			Err: "Error: unknown command",
		})
		return
	}

	writeResponse(w, &response{
		Data: "OK",
	})
}

// ServoRotate is the http wrapper for devices.ServoRotate().
func (m *Server) ServoRotate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if m.ctrl == nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: Controller not initialized"),
		})
		return
	}

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	dir := strings.ToLower(r.Form.Get("dir")) // Servo button { up, down, left, right}.

	delta := int(m.servoDelta)
	var servo byte

	switch dir {
	case "up":
		servo = 2
		m.servoAngle[servo] -= delta
	case "down":
		servo = 2
		m.servoAngle[servo] += delta
	case "right":
		servo = 1
		m.servoAngle[servo] -= delta
	case "left":
		servo = 1
		m.servoAngle[servo] += delta
	}

	// Set rotation boundary angles.
	if m.servoAngle[servo] < 0 {
		writeResponse(w, &response{
			Err: "Error: servo angle below 0",
		})
		m.servoAngle[servo] = 0
		return
	}
	if m.servoAngle[servo] > 180 {
		writeResponse(w, &response{
			Err: "Error servo angle above 180",
		})
		m.servoAngle[servo] = 180
		return
	}

	if err := m.ctrl.ServoRotate(servo, m.servoAngle[servo]); err != nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: servo failed %v", err),
		})
		return
	}

	writeResponse(w, &response{
		Data: map[string]int{"horiz": m.servoAngle[byte(1)], "vert": m.servoAngle[byte(2)]},
	})
}

// RoombaCmd sets the roomba mode.
func (m *Server) RoombaCmd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if m.roomba == nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: Roomba not initialized"),
		})
		return
	}

	if err := r.ParseForm(); err != nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error parsing form %v", err),
		})
		return
	}

	cmd := strings.ToLower(r.Form.Get("cmd"))
	var err error
	glog.V(2).Infof("HTTP command %v", cmd)

	switch cmd {
	case "safe_mode":
		err = m.roomba.Safe()
		err = m.roomba.MainBrush(true, true)

	case "full_mode":
		err = m.roomba.Full()
		err = m.roomba.MainBrush(true, true)

	case "passive_mode":
		err = m.roomba.Passive()
		m.data.Enabled = map[byte]bool{
			TEMP:     false,
			HUMIDITY: false,
			LDR:      false,
			PIR:      false,
			MAG:      false,
			BATT:     false,
		}

	case "power_off":
		m.roomba.Start(false)
		err = m.roomba.Power()
		err = m.roomba.MainBrush(false, true)
		m.data.Enabled = map[byte]bool{
			TEMP:     false,
			HUMIDITY: false,
			LDR:      false,
			PIR:      false,
			MAG:      false,
			BATT:     false,
		}

	case "power_on":
		err = m.roomba.Start(true)
		err = m.roomba.Safe()
		time.Sleep(100 * time.Millisecond)
		err = m.roomba.MainBrush(true, true)
		m.data.Enabled = map[byte]bool{
			TEMP:     true,
			HUMIDITY: true,
			LDR:      true,
			PIR:      true,
			MAG:      true,
			BATT:     true,
		}

	case "seek_dock":
		err = m.roomba.SeekDock()
	}

	if err != nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error changing mode %v", err),
		})
	}
}
