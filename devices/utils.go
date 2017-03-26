package devices

import (
	"time"

	_ "github.com/kidoman/embd/host/chip"
	"github.com/kidoman/embd/sensor/us020"
)

func ForwardSweep(ctrl *Controller, us *us020.US020, angle int) ([]int32, error) {
	val := []int32{}

	// Sleep to allow servo to move to starting position.
	time.Sleep(320 * time.Millisecond)
	for i := 20; i <= 160; i += angle {
		if err := ctrl.ServoRotate(1, i); err != nil {
			return nil, err
		}

		time.Sleep(40 * time.Millisecond) // Sleep to allow servo to finish turning.
		dist, err := us.Distance()
		if err == nil {
			break
		}

		val = append(val, int32(dist))
		time.Sleep(100 * time.Millisecond)
	}

	return val, nil
}
