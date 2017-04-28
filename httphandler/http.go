package httphandler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/deepakkamesh/sonny/devices"
	"github.com/deepakkamesh/sonny/rpc"
	"github.com/golang/glog"
	_ "github.com/kidoman/embd/host/chip"
	"github.com/kidoman/embd/sensor/hmc5883l"
	"github.com/kidoman/embd/sensor/us020"
)

type Server struct {
	ctrl       *devices.Controller
	mag        *hmc5883l.HMC5883L
	us         *us020.US020
	pir        string
	ssl        bool
	resources  string
	servoAngle map[byte]int // Map to hold state of each servo.
}

// Struct to return JSON/
type response struct {
	Err  string
	Data interface{}
}

func New(d *rpc.Devices, ssl bool, resources string) *Server {
	return &Server{
		ctrl:       d.Ctrl,
		mag:        d.Mag,
		pir:        d.Pir,
		us:         d.Us,
		ssl:        ssl,
		resources:  resources,
		servoAngle: map[byte]int{1: 90, 2: 90},
	}
}

func (m *Server) Start() error {

	http.HandleFunc("/", m.ServeIndex)
	http.HandleFunc("/api/ping/", m.Ping)
	http.HandleFunc("/api/ledon/", m.LEDOn)
	http.HandleFunc("/api/ledblink/", m.LEDBlink)
	http.HandleFunc("/api/servorotate/", m.ServoRotate)
	http.HandleFunc("/api/distance/", m.Distance)
	http.HandleFunc("/api/batt/", m.BattState)

	return http.ListenAndServe(":8080", nil)
}

func (m *Server) ServeIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	// TODO: change to http.ServeFile.
	idx, err := ioutil.ReadFile(m.resources + "/index.html")
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
	if _, err = w.Write(idx); err != nil {
		glog.Warning("Unable to write http response")
	}
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

// Ping is a http wrapper for devices.Ping.
func (m *Server) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if m.ctrl == nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: Controller not enabled"),
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

func (m *Server) Distance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	d, err := m.us.Distance()
	if err != nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: distance query failed %v", err),
		})
		return
	}

	writeResponse(w, &response{
		Data: fmt.Sprintf("%3.3f", d),
	})
}

func (m *Server) BattState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if m.ctrl == nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: Controller not enabled"),
		})
		return
	}

	val, err := m.ctrl.BattState()
	if err != nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: failed to get battery level %v", err),
		})
		return
	}

	writeResponse(w, &response{
		Data: val,
	})
}

// LEDBlink is the http wrapper for devices.LEDBlink().
func (m *Server) LEDBlink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if m.ctrl == nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: Controller not enabled"),
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
			Err: fmt.Sprintf("Error: Controller not enabled"),
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
	w.Header().Set("Content-Type", "text/html")

	if m.ctrl == nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: Controller not enabled"),
		})
		return
	}

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	dir := strings.ToLower(r.Form.Get("dir")) // Servo button { up, down, left, right}.

	const delta = 10
	var servo byte

	switch dir {
	case "up":
		servo = 2
		m.servoAngle[servo] -= delta
	case "down":
		servo = 2
		m.servoAngle[servo] += delta
	case "left":
		servo = 1
		m.servoAngle[servo] -= delta
	case "right":
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
