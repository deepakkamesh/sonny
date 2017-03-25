package httphandler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/deepakkamesh/sonny/devices"
	"github.com/deepakkamesh/sonny/rpc"
	"github.com/golang/glog"
	_ "github.com/kidoman/embd/host/chip"
	"github.com/kidoman/embd/sensor/lsm303"
	"github.com/kidoman/embd/sensor/us020"
)

type Server struct {
	ctrl       *devices.Controller
	mag        *lsm303.LSM303
	us         *us020.US020
	pir        string
	ssl        bool
	resources  string
	servoAngle map[byte]byte // Map to hold state of each servo.
}

func New(d *rpc.Devices, ssl bool, resources string) *Server {
	return &Server{
		ctrl:       d.Ctrl,
		mag:        d.Mag,
		pir:        d.Pir,
		us:         d.Us,
		ssl:        ssl,
		resources:  resources,
		servoAngle: map[byte]byte{1: 90, 2: 90},
	}
}

func (m *Server) Start() error {

	// Validate devices.
	if m.ctrl == nil {
		return errors.New("Controller not enabled")
	}

	http.HandleFunc("/", m.ServeIndex)
	http.HandleFunc("/api/ping", m.Ping)
	http.HandleFunc("/api/ledon/", m.LEDOn)
	http.HandleFunc("/api/ledblink/", m.LEDBlink)
	http.HandleFunc("/api/servorotate/", m.ServoRotate)
	http.HandleFunc("/api/distance/", m.Distance)
	return http.ListenAndServe(":8080", nil)
}

func (m *Server) ServeIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	idx, err := ioutil.ReadFile(m.resources + "/index.html")
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
	if _, err = w.Write(idx); err != nil {
		glog.Warning("Unable to write http response")
	}
}

func (m *Server) Distance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	d, err := m.us.Distance()
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	fmt.Fprintf(w, "%v", d)
}

// Ping is a http wrapper for devices.Ping.
func (m *Server) Ping(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	if err := m.ctrl.Ping(); err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	fmt.Fprintf(w, "%s", "OK")
}

// LEDBlink is the http wrapper for devices.LEDBlink().
func (m *Server) LEDBlink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	dur, err := strconv.ParseUint(strings.ToLower(r.Form.Get("duration")), 10, 16) // Duration of blink in ms.
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	no, err := strconv.ParseUint(strings.ToLower(r.Form.Get("times")), 10, 8) // Number of times to blink.
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	if err := m.ctrl.LEDBlink(uint16(dur), byte(no)); err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	fmt.Fprint(w, "OK")
}

// LEDOn is the http wrapper for devices.LEDOn().
func (m *Server) LEDOn(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	a := strings.ToLower(r.Form.Get("cmd")) // valid values = {on,off}.

	switch a {
	case "on":
		if err := m.ctrl.LEDOn(true); err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}
	case "off":
		if err := m.ctrl.LEDOn(false); err != nil {
			fmt.Fprintf(w, "Error: %v", err)
			return
		}
	default:
		fmt.Fprintf(w, "Error: unknown cmd %v", a)
		return
	}

	fmt.Fprint(w, "OK")
}

// ServoRotate is the http wrapper for devices.ServoRotate().
func (m *Server) ServoRotate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

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
		m.servoAngle[servo] += delta
	case "down":
		servo = 2
		m.servoAngle[servo] -= delta
	case "left":
		servo = 1
		m.servoAngle[servo] += delta
	case "right":
		servo = 1
		m.servoAngle[servo] -= delta
	}

	// Set rotation boundary angles.
	if m.servoAngle[servo] < 0 {
		m.servoAngle[servo] = 0
	}
	if m.servoAngle[servo] > 180 {
		m.servoAngle[servo] = 180
	}

	if err := m.ctrl.ServoRotate(servo, m.servoAngle[servo]); err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	fmt.Fprint(w, "OK")
}
