package httphandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/deepakkamesh/sonny/devices"
	"github.com/deepakkamesh/sonny/navigator"
	"github.com/golang/glog"
)

type Server struct {
	sonny      *devices.Sonny
	navigator  *navigator.Ogrid
	ssl        bool
	resources  string
	servoAngle map[byte]int // Map to hold state of each servo.
	servoDelta uint8
	velocity   int16
	timer      *time.Timer
	data       *sensorData
	connCount  int // Count of wvebsockets.
}

// Struct to return JSON.
type response struct {
	Err  string
	Data interface{}
}

func New(d *devices.Sonny, n *navigator.Ogrid, ssl bool, resources string) *Server {
	t := time.NewTimer(500 * time.Millisecond)
	t.Stop()

	return &Server{
		sonny:      d,
		navigator:  n,
		ssl:        ssl,
		resources:  resources,
		servoAngle: map[byte]int{1: 90, 2: 90},
		servoDelta: 10,
		velocity:   100,
		timer:      t,
		data: &sensorData{
			Controller: make(map[byte]float32),
			Roomba:     make(map[byte]int16),
			Pi:         make(map[byte]int),
		},
	}
}

func (m *Server) Start(hostPort string) error {

	// http routers.
	http.HandleFunc("/api/setparam/", m.SetParam)
	http.HandleFunc("/api/ping/", m.Ping)
	http.HandleFunc("/api/ledon/", m.LEDOn)
	http.HandleFunc("/api/ledblink/", m.LEDBlink)
	http.HandleFunc("/api/servorotate/", m.ServoRotate)
	http.HandleFunc("/api/move/", m.Move)
	http.HandleFunc("/api/roomba_cmd/", m.RoombaCmd)
	http.HandleFunc("/api/i2c_en/", m.I2CEn)
	http.HandleFunc("/datastream", m.dataStream)
	http.HandleFunc("/gridDisp", m.gridDisp)
	http.HandleFunc("/api/navi/", m.Navi)

	// Serve static content from resources dir.
	fs := http.FileServer(http.Dir(m.resources))
	http.Handle("/", fs)

	// Startup data collection routine.
	go m.dataCollector()
	return http.ListenAndServe(hostPort, nil)
}

// writeResponse writes the response json object to w. If unable to marshal
// it writes a http 500.
func writeResponse(w http.ResponseWriter, resp *response) {
	w.Header().Set("Content-Type", "application/json")
	js, e := json.Marshal(resp)
	if e != nil {
		http.Error(w, e.Error(), http.StatusInternalServerError)
		return
	}
	glog.V(3).Infof("Writing json response %s", js)
	w.Write(js)
}

// I2CEn enables or disables I2C bus.
func (m *Server) I2CEn(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	var err error
	action := strings.ToLower(r.Form.Get("param"))
	if action != "" {
		switch action {
		case "on":
			err = m.sonny.I2CBusEnable(true)
		case "off":
			err = m.sonny.I2CBusEnable(false)
		}
	}

	if err != nil {
		glog.Errorf("Failed to turn %v  I2C bus: %v", action, err)
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: Failed to turn %v I2C bus: %v", action, err),
		})
		return
	}

	writeResponse(w, &response{
		Data: "OK",
	})

}

// Navi is a test function for navigation.
func (m *Server) Navi(w http.ResponseWriter, r *http.Request) {
	glog.Info("Navi button pressed")
	if err := m.navigator.UpdateMap(); err != nil {
		glog.Errorf("Navi failure: %v", err)
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: update map failed %v", err),
		})
		return
	}
	writeResponse(w, &response{
		Data: "OK",
	})
}

// gridDisp streams the png image with the grid map.
func (m *Server) gridDisp(w http.ResponseWriter, r *http.Request) {

	buffer, err := m.navigator.GenerateMap()
	if err != nil {
		glog.Errorf("Failed to generate map: %v", err)
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: Failed to generate map: %v", err),
		})
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		glog.Errorf("Unable to write image: %v", err)
	}

}

// Ping is a http wrapper for devices.Ping.
func (m *Server) Ping(w http.ResponseWriter, r *http.Request) {

	if err := m.sonny.Ping(); err != nil {
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

	if m.sonny.Roomba == nil {
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
		err = m.sonny.Drive(m.velocity, 32767)

	case "bwd":
		err = m.sonny.Drive(-1*m.velocity, 32767)

	case "right":
		err = m.sonny.Drive(m.velocity, -1)

	case "left":
		err = m.sonny.Drive(m.velocity, 1)
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
			m.sonny.Drive(0, 0)
		}()
	}
}

// LEDBlink is the http wrapper for devices.LEDBlink().
func (m *Server) LEDBlink(w http.ResponseWriter, r *http.Request) {

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

	if err := m.sonny.LEDBlink(uint16(dur), byte(no)); err != nil {
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

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	a := strings.ToLower(r.Form.Get("cmd")) // valid values = {on,off}.

	switch a {
	case "on":
		if err := m.sonny.LEDOn(true); err != nil {
			writeResponse(w, &response{
				Err: fmt.Sprintf("Error: LED failed %v", err),
			})
			return
		}
	case "off":
		if err := m.sonny.LEDOn(false); err != nil {
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

	if err := m.sonny.ServoRotate(servo, m.servoAngle[servo]); err != nil {
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

	if m.sonny.Roomba == nil {
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
	param := strings.ToLower(r.Form.Get("param"))

	var err error
	glog.V(2).Infof("Roomba command %v", cmd)

	switch cmd {
	case "safe_mode":
		err = m.sonny.Safe()

	case "aux_power":
		if m.sonny.GetRoombaMode() <= 1 {
			err = fmt.Errorf("Aux can only be enabled in Safe or Full mode")
			break
		}
		switch param {
		case "on":
			err = m.sonny.AuxPower(true)
		case "off":
			err = m.sonny.AuxPower(false)
		}

	case "reset":
		err = m.sonny.Reset()

	case "full_mode":
		err = m.sonny.Full()

	case "passive_mode":
		err = m.sonny.Passive()

	case "power_off":
		err = m.sonny.Power()

	case "power_on":
		err = m.sonny.Roomba.Start(true)

	case "seek_dock":
		err = m.sonny.SeekDock()
	}

	if err != nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error execution roomba cmd: %v", err),
		})
		return
	}
	writeResponse(w, &response{
		Data: "OK",
	})
}
