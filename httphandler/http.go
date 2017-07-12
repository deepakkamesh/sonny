package httphandler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	roomba "github.com/deepakkamesh/go-roomba"
	"github.com/deepakkamesh/go-roomba/constants"
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
	roomba     *roomba.Roomba
	ssl        bool
	resources  string
	servoAngle map[byte]int // Map to hold state of each servo.
	servoDelta uint8
	velocity   int16
	timer      *time.Timer
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
}

func New(d *rpc.Devices, ssl bool, resources string) *Server {
	t := time.NewTimer(500 * time.Millisecond)
	t.Stop()

	return &Server{
		ctrl:       d.Ctrl,
		mag:        d.Mag,
		pir:        d.Pir,
		roomba:     d.Roomba,
		ssl:        ssl,
		resources:  resources,
		servoAngle: map[byte]int{1: 90, 2: 90},
		servoDelta: 10,
		velocity:   100,
		timer:      t,
	}

}

func (m *Server) Start() error {

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

	for {
		envData, errStr1 := m.getEnvData()
		rbData, errStr2 := m.getRoombaData()

		m := &sensorData{
			Err:        errStr1 + errStr2,
			Roomba:     rbData,
			Controller: envData,
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
		time.Sleep(800 * time.Millisecond)
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

// getRoombaData returns the current value of the roomba sensors.
func (m *Server) getRoombaData() (data map[byte]int16, errStr string) {

	if m.roomba == nil {
		return nil, "roomba not initialized"
	}

	data = make(map[byte]int16)

	sg := []byte{constants.SENSOR_GROUP_6, constants.SENSOR_GROUP_101}
	pg := [][]byte{constants.PACKET_GROUP_6, constants.PACKET_GROUP_101}

	// Iterate through the packet groups. Sensor group 100 does not work as advertised.
	// Use sensor group, 6 and 101 instead.
	for grp := 0; grp < 2; grp++ {
		d, e := m.roomba.Sensors(sg[grp])
		if e != nil {
			errStr = errStr + e.Error()
			glog.Errorf("Failed to read sensors: %v", e)
		}

		i := byte(0)
		for _, p := range pg[grp] {
			pktL := constants.SENSOR_PACKET_LENGTH[p]

			if pktL == 1 {
				data[p] = int16(d[i])
			}
			if pktL == 2 {
				v := int16(d[i])<<8 | int16(d[i+1])
				data[p] = v
			}
			i = i + pktL
		}
	}
	return
}

// getEnvData queries the environment sensors and returns a map with with the data.
func (m *Server) getEnvData() (data map[byte]float32, errStr string) {

	data = make(map[byte]float32)

	for i := 0; i < 6; i++ {
		switch i {
		case 0:
			if m.ctrl == nil {
				errStr = errStr + " controller not initialized"
				continue
			}

			t, h, err := m.ctrl.DHT11()
			if err != nil {
				glog.Errorf("Failed to read DHT11: %v", err)
				errStr = errStr + err.Error()
				continue
			}
			data[0] = float32(t)*1.8 + 32
			data[1] = float32(h)

		case 2:
			if m.ctrl == nil {
				errStr = errStr + " controller not initialized"
				continue
			}
			l, err := m.ctrl.LDR()
			if err != nil {
				glog.Errorf("Failed to read LDR: %v", err)
				errStr = errStr + err.Error()
				continue
			}
			data[2] = float32(l)

		case 3:
			if m.pir == nil {
				errStr = errStr + " PIR  not initialized"
				continue
			}
			v, err := m.pir.Detect()
			if err != nil {
				glog.Errorf("Failed to read PIR: %v", err)
				errStr = errStr + err.Error()
				continue
			}
			if v {
				data[3] = float32(1)
				continue
			}
			data[3] = float32(0)

		case 4:
			if m.mag == nil {
				errStr = errStr + " compass not initialized"
				continue
			}
			h, err := m.mag.Heading()
			if err != nil {
				glog.Errorf("Failed to read Compass: %v", err)
				errStr = errStr + err.Error()
				continue
			}
			data[4] = float32(h)

		case 5:
			if m.ctrl == nil {
				errStr = errStr + " controller not initialized"
				continue
			}
			b, err := m.ctrl.BattState()
			if err != nil {
				glog.Errorf("Failed to read controller batt state: %v", err)
				errStr = errStr + err.Error()
				continue
			}
			data[5] = float32(b)

		}
	}
	return
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

	dir := strings.ToLower(r.Form.Get("dir")) // Motor button { up, down, left, right}
	var err error

	/*	err = m.roomba.Drive(-100, 32767)
		time.Sleep(500 * time.Millisecond)
		m.roomba.Drive(0, 0)*/

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

	//time.Sleep(500 * time.Millisecond)
	//m.roomba.Drive(0, 0)
	if run := m.timer.Reset(500 * time.Millisecond); !run {
		glog.V(2).Info("start timer")
		go func() {
			<-m.timer.C
			glog.V(2).Infof("timer expired stop")
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

	case "full_mode":
		err = m.roomba.Full()

	case "passive_mode":
		err = m.roomba.Passive()

	case "power_off":
		m.roomba.Start(false)
		err = m.roomba.Power()

	case "power_on":
		err = m.roomba.Start(true)

	case "seek_dock":
		err = m.roomba.SeekDock()
	}

	if err != nil {
		writeResponse(w, &response{
			Err: fmt.Sprintf("Error changing mode %v", err),
		})
	}
}
