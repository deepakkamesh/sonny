package devices

import (
	"errors"
	"log"
	"time"

	rpio "github.com/stianeikeland/go-rpio"
)

type Pi struct {
	pir    rpio.Pin
	usEcho rpio.Pin
	usTrig rpio.Pin
}

func NewPi() *Pi {

	// Open GPIO
	if err := rpio.Open(); err != nil {
		log.Fatalf("Unable to open gpio port on Pi")
	}

	// Setup IO pins
	pir := rpio.Pin(18) // PIR Input.
	pir.Input()
	usTrig := rpio.Pin(23) // Ultrasonic Trigger Output.
	usTrig.Output()
	usEcho := rpio.Pin(24) // Ultrasonic Echo Input.
	usEcho.Input()

	return &Pi{
		pir:    pir,
		usEcho: usEcho,
		usTrig: usTrig,
	}
}

func (m *Pi) PIRDetect() bool {
	if m.pir.Read() == rpio.High {
		return true
	}
	return false
}

// Distance returns the distance from the ultrasonic sensor in inches.
func (m *Pi) Distance() (float64, error) {

	done := make(chan float64)
	t := time.NewTimer(time.Millisecond * 16) // Set a timeout as max distance is 4m.

	// Send Trigger pulse for 10uS.
	m.usTrig.High()
	time.Sleep(time.Microsecond * 10)
	m.usTrig.Low()

	// Measure return pulse.
	go func() {
		var start, end time.Time

		for st := m.usEcho.Read(); st == rpio.Low; {
			start = time.Now()
			st = m.usEcho.Read()
		}

		for st := m.usEcho.Read(); st == rpio.High; {
			end = time.Now()
			st = m.usEcho.Read()
		}
		done <- end.Sub(start).Seconds()
	}()

	select {
	case <-t.C:
		return 0, errors.New("timeout on ultrasonic pulse")
	case d := <-done:
		return d * 17150 * 0.39370, nil
	}

}
