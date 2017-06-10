package devices

import "time"

func ForwardSweep(ctrl *Controller, angle int) ([]int32, error) {
	val := []int32{}

	// Sleep to allow servo to move to starting position.
	time.Sleep(320 * time.Millisecond)
	for i := 20; i <= 160; i += angle {
		if err := ctrl.ServoRotate(1, i); err != nil {
			return nil, err
		}

		time.Sleep(40 * time.Millisecond) // Sleep to allow servo to finish turning.
		dist, err := ctrl.Distance()
		if err == nil {
			break
		}

		val = append(val, int32(dist))
		time.Sleep(100 * time.Millisecond)
	}

	return val, nil
}
