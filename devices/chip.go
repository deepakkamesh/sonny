// +build ignore
package devices

import (
	"fmt"
	"log"
	"time"

	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/chip"
)

const (
	PIR     = "gpio0"
	US_TRIG = "gpio1"
	US_ECHO = "gpio2"
)

type Chip struct {
}

func NewChip() *Chip {

	// Open GPIO
	if err := embd.InitGPIO(); err != nil {
		log.Fatalf("Unable to open GPIO: %v", err)
	}

	// Setup IO pins.
	embd.SetDirection(PIR, embd.In)
	embd.SetDirection(US_TRIG, embd.Out)
	embd.SetDirection(US_ECHO, embd.In)
	embd.SetDirection("gpio7", embd.Out)

	return &Chip{}
}

func (m *Chip) PIRDetect() (bool, error) {
	v, err := embd.DigitalRead(PIR)
	if err != nil {
		return false, err
	}

	if v == embd.High {
		return true, nil
	}
	return false, nil
}

func (m *Chip) BlinkGPIO8() {
	on := 0
	for {
		embd.DigitalWrite("gpio7", on)
		on = 1 - on
		time.Sleep(550 * time.Millisecond)
		fmt.Printf("Sett state %d\n", on)
	}
}

/*
// Distance returns the distance from the ultrasonic sensor in inches.
func (m *Chip) Distance() (float64, error) {

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

}*/
