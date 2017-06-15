package httphandler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/deepakkamesh/sonny/devices"
	"github.com/deepakkamesh/sonny/rpc"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"github.com/kidoman/embd/sensor/hcsr501"
	"github.com/kidoman/embd/sensor/hmc5883l"
)

type Server struct {
	ctrl       *devices.Controller
	mag        *hmc5883l.HMC5883L
	pir        *hcsr501.HCSR501
	ssl        bool
	resources  string
	servoAngle map[byte]int // Map to hold state of each servo.
}

// Struct to return JSON/
type response struct {
	Err  string
	Data interface{}
}

// sensor data struct.
type sensorData struct {
	Err    string
	Roomba map[int]int
}

func New(d *rpc.Devices, ssl bool, resources string) *Server {
	return &Server{
		ctrl:       d.Ctrl,
		mag:        d.Mag,
		pir:        d.Pir,
		ssl:        ssl,
		resources:  resources,
		servoAngle: map[byte]int{1: 90, 2: 90},
	}
}

func (m *Server) Start() error {

	http.HandleFunc("/api/ping/", m.Ping)
	http.HandleFunc("/api/ledon/", m.LEDOn)
	http.HandleFunc("/api/ledblink/", m.LEDBlink)
	http.HandleFunc("/api/servorotate/", m.ServoRotate)
	http.HandleFunc("/api/distance/", m.Distance)
	http.HandleFunc("/api/move/", m.Move)
	http.HandleFunc("/datastream", m.dataStream) // Websocket Handler.

	// State functions. TODO: Move to websock.
	http.HandleFunc("/api/batt/", m.BattState)
	http.HandleFunc("/api/accel/", m.Accelerometer)
	http.HandleFunc("/api/head/", m.Heading)
	http.HandleFunc("/api/temp/", m.DHT11)
	http.HandleFunc("/api/ldr/", m.LDR)
	http.HandleFunc("/api/pir/", m.PIRDetect)

	// Serve static content from resources dir.
	fs := http.FileServer(http.Dir(m.resources))
	http.Handle("/", fs)

	return http.ListenAndServe(":8080", nil)
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

// dataStream is the websocket server that streams rover stats.
func (m *Server) dataStream(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		glog.Warningf("failed to upgrade conn:%v", err)
		return
	}
	defer c.Close()

	ra := rand.New(rand.NewSource(99))
	for {
		// TODO: Gen ran data;

		data := map[int]int{}
		for i := 7; i < 58; i++ {
			data[i] = ra.Intn(255)
		}

		m := &sensorData{
			Err:    "none",
			Roomba: data,
		}
		jsMsg, err := json.Marshal(m)
		if err != nil {
			glog.Errorf("failed to unmarshall:%v", err)
		}

		err = c.WriteMessage(websocket.TextMessage, jsMsg)
		if err != nil {
			glog.Errorf("failed to write:%v", err)
			break
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

func (m *Server) Distance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if m.ctrl == nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: Controller not initialized"),
		})
		return
	}
	d, err := m.ctrl.Distance()
	if err != nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: distance query failed %v", err),
		})
		return
	}

	writeResponse(w, &response{
		Data: d,
	})
}

// Accelerometer is the http wrapper for ctrl.Accelerator.
func (m *Server) Accelerometer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if m.ctrl == nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: controller not initialized"),
		})
		return
	}

	x, y, z, err := m.ctrl.Accelerometer()
	if err != nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: failed to read accelerometer %v", err),
		})
		return
	}

	writeResponse(w, &response{
		Data: []float32{x, y, z},
	})
}

// Move is the wrapper around ctrl.Move.
func (m *Server) Move(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if m.ctrl == nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: controller not initialized"),
		})
		return
	}

	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}

	//dir := strings.ToLower(r.Form.Get("dir")) // Motor button { up, down, left, right}.
	/*
		switch dir {
		case "forward":
			// TODO: Remove hardcoded values for turns and duty percent.
			if _, _, err := m.ctrl.Move(20, true, 90); err != nil {
				writeResponse(w, &response{
					Err: fmt.Sprintf("Error: motor failed %v", err),
				})
				return
			}
		case "back":
			if _, _, err := m.ctrl.Move(20, false, 90); err != nil {
				writeResponse(w, &response{
					Err: fmt.Sprintf("Error: motor failed %v", err),
				})
				return
			}
		case "left":
			if _, _, err := m.ctrl.Turn(10, 1, 90); err != nil {
				writeResponse(w, &response{
					Err: fmt.Sprintf("Error: motor failed %v", err),
				})
				return
			}
		case "right":
			if _, _, err := m.ctrl.Turn(10, 0, 90); err != nil {
				writeResponse(w, &response{
					Err: fmt.Sprintf("Error: motor failed %v", err),
				})
				return
			}
		} */
}

// Heading is a http wrapper for mag.HEading.
func (m *Server) Heading(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if m.mag == nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: magnetometer not initialized"),
		})
		return
	}

	h, err := m.mag.Heading()
	if err != nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: failed to read magnetometer %v", err),
		})
		return
	}

	writeResponse(w, &response{
		Data: h,
	})
}

// Heading is a http wrapper for mag.HEading.
func (m *Server) DHT11(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if m.ctrl == nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: controller not initialized"),
		})
		return
	}

	t, h, err := m.ctrl.DHT11()
	if err != nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: failed to read DHT11 %v", err),
		})
		return
	}

	writeResponse(w, &response{
		Data: []uint16{uint16(t), uint16(h)},
	})
}

// LDR is a http wrapper for ctrl.LDR
func (m *Server) LDR(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if m.ctrl == nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: controller not initialized"),
		})
		return
	}

	v, err := m.ctrl.LDR()
	if err != nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: failed to read controller %v", err),
		})
		return
	}

	writeResponse(w, &response{
		Data: v,
	})
}

// PIRDetect is a http wrapper for pir.Detect.
func (m *Server) PIRDetect(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if m.pir == nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: PIR not initialized"),
		})
		return
	}

	v, err := m.pir.Detect()
	if err != nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: failed to read pir %v", err),
		})
		return
	}

	writeResponse(w, &response{
		Data: v,
	})
}

// BattState is the http wrapper for ctrl.BattState
func (m *Server) BattState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if m.ctrl == nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error: Controller not initialized"),
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
	w.Header().Set("Content-Type", "text/html")

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

	const delta = 10
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
