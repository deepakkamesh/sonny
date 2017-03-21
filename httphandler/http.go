package httphandler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/deepakkamesh/sonny/devices"
	"github.com/deepakkamesh/sonny/rpc"
	"github.com/golang/glog"
	_ "github.com/kidoman/embd/host/chip"
	"github.com/kidoman/embd/sensor/lsm303"
	"github.com/kidoman/embd/sensor/us020"
)

type Server struct {
	ctrl *devices.Controller
	mag  *lsm303.LSM303
	us   *us020.US020
	pir  string
	ssl  bool
}

func New(d *rpc.Devices, ssl bool) *Server {
	return &Server{
		ctrl: d.Ctrl,
		mag:  d.Mag,
		pir:  d.Pir,
		us:   d.Us,
		ssl:  ssl,
	}
}

func (m *Server) Start() error {

	http.HandleFunc("/", m.ServeIndex)
	http.HandleFunc("/api/ping", m.Ping)
	http.HandleFunc("/api/led/", m.LEDOn)
	return http.ListenAndServe(":8080", nil)
	//return nil
}

func (m *Server) ServeIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	idx, err := ioutil.ReadFile("./index.html")
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
	}
	if _, err = w.Write(idx); err != nil {
		glog.Warning("Unable to write http response")
	}
}

func (m *Server) Ping(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	if m.ctrl == nil {
		fmt.Fprintf(w, "%s", "Error: Controller not enabled")
		return
	}

	if err := m.ctrl.Ping(); err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	fmt.Fprintf(w, "%s", "OK")
}

func (m *Server) LEDOn(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")

	if m.ctrl == nil {
		fmt.Fprintf(w, "%s", "Error: Controller not enabled")
		return
	}

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
